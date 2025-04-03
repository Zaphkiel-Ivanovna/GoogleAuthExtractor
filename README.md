# 🔐 Google Authenticator Secret Extractor

[![CI](https://github.com/zaphkiel-ivanovna/googleauthextractor/actions/workflows/ci.yml/badge.svg)](https://github.com/zaphkiel-ivanovna/googleauthextractor/actions/workflows/ci.yml)
[![Lint](https://github.com/zaphkiel-ivanovna/googleauthextractor/actions/workflows/lint.yml/badge.svg)](https://github.com/zaphkiel-ivanovna/googleauthextractor/actions/workflows/lint.yml)
[![Docker Package](https://img.shields.io/badge/Container-ghcr.io-blue)](https://github.com/zaphkiel-ivanovna/googleauthextractor/pkgs/container/googleauthextractor)
[![GitHub Release](https://img.shields.io/github/v/release/zaphkiel-ivanovna/googleauthextractor)](https://github.com/zaphkiel-ivanovna/googleauthextractor/releases)

Extract TOTP/HOTP secrets from Google Authenticator export QR codes with ease.

> **Note:** This project is not affiliated with Google.

## 📚 Table of Contents

- [✨ Features](#-features)
- [📦 Installation](#-installation)
  - [📥 Prebuilt Binaries](#-prebuilt-binaries)
  - [🛠️ From Source](#️-from-source)
  - [🐳 Using Docker](#-using-docker)
- [🧰 Usage](#-usage)
  - [📺 View in Terminal](#-view-in-terminal)
  - [📄 Export to JSON](#-export-to-json)
  - [🔄 Generate QR Codes](#-generate-qr-codes)
  - [📋 Command Line Reference](#-command-line-reference)
  - [Legacy Mode](#legacy-mode)
- [📱 How to Export from Google Authenticator](#-how-to-export-from-google-authenticator)
- [🔑 Understanding Secret Formats](#-understanding-secret-formats)
- [🔒 Security Considerations](#-security-considerations)
- [📋 Data Format](#-data-format)
- [🔄 Migration Guide](#-migration-guide)
  - [To Authy](#to-authy)
  - [To Bitwarden](#to-bitwarden)
  - [To 1Password](#to-1password)
  - [To KeePass (with KeePassOTP plugin)](#to-keepass-with-keepassotp-plugin)
- [🧪 Development](#-development)
  - [🔄 CI/CD Workflows](#️-cicd-workflows)
  - [Protocol Buffer](#protocol-buffer)
- [📄 License](#-license)
- [👏 Acknowledgments](#-acknowledgments)


## ✨ Features

- **🔒 Secure Processing**: Handle your 2FA secrets locally without external services
- **🖼️ QR Image Processing**: Extract directly from screenshots containing QR codes
- **📤 Flexible Output**:
  - 📄 Export to JSON for backup or custom processing
  - 🔄 Generate individual QR codes for each account to scan with other apps
  - 🖥️ Pretty print account details directly in your terminal
  - 📟 Display QR codes as ASCII art in the terminal
  - 🔑 View full secrets securely when needed
- **🔄 Easy Migration**: Move your accounts to any authenticator app (Authy, Bitwarden, etc.)

## 📦 Installation

### 📥 Prebuilt Binaries

Download the latest release from [GitHub Releases](https://github.com/Zaphkiel-Ivanovna/GoogleAuthExtractor/releases).

### 🛠️ From Source

Requirements:

- Go 1.24 or higher

```bash
# Clone the repository
git clone https://github.com/Zaphkiel-Ivanovna/GoogleAuthExtractor.git
cd GoogleAuthExtractor

# Install dependencies
go mod download

# Build
go build -o gauth-extractor ./cmd/extractor

# Install system-wide (optional)
go install ./cmd/extractor
```

### 🐳 Using Docker

```bash
# Pull from GitHub Container Registry
docker pull ghcr.io/zaphkiel-ivanovna/googleauthextractor:latest

# Run in interactive mode
docker run -it --rm -v "$(pwd):/home/appuser/data" ghcr.io/zaphkiel-ivanovna/googleauthextractor:latest

# Or build locally
docker build -t gauth-extractor .
docker run -it --rm -v "$(pwd):/home/appuser/data" gauth-extractor
```

## 🧰 Usage

The CLI has been restructured with a more intuitive command system. There are three main commands:

- `view` - Display accounts in the terminal
- `json` - Export accounts to JSON format
- `qr` - Generate QR codes for each account

### Input Methods

All commands support these input methods (specify one):

```bash
# Interactive mode (will prompt for URI)
gauth-extractor <command> -i

# From URI string
gauth-extractor <command> -u "otpauth-migration://offline?data=..."

# From QR code image
gauth-extractor <command> -q "/path/to/qrcode-screenshot.png"
```

### 📺 View in Terminal

```bash
# View accounts in terminal with pretty formatting (default)
gauth-extractor view -u "otpauth-migration://offline?data=..."

# Simple table view (disable pretty print)
gauth-extractor view -u "otpauth-migration://offline?data=..." --pretty=false

# Show QR codes in terminal too
gauth-extractor view -u "otpauth-migration://offline?data=..." -r

# Display full secrets (USE WITH CAUTION)
gauth-extractor view -u "otpauth-migration://offline?data=..." -s

# Combine options
gauth-extractor view -u "otpauth-migration://offline?data=..." -r -s
```

### 📄 Export to JSON

```bash
# Save to JSON file (default: accounts.json)
gauth-extractor json -u "otpauth-migration://offline?data=..."

# Specify custom filename
gauth-extractor json -u "otpauth-migration://offline?data=..." -f "my-accounts.json"

# Print JSON to terminal instead of saving
gauth-extractor json -u "otpauth-migration://offline?data=..." -s=false
```

### 🔄 Generate QR Codes

```bash
# Save QR codes to directory (default: ./qrcodes)
gauth-extractor qr -u "otpauth-migration://offline?data=..."

# Specify custom directory
gauth-extractor qr -u "otpauth-migration://offline?data=..." -d "my-qrcodes"

# Display QR codes in terminal instead of saving files
gauth-extractor qr -u "otpauth-migration://offline?data=..." -s=false
```

### 📋 Command Line Reference

```
Usage:
  gauth-extractor [command]

Available Commands:
  json        Export accounts to JSON format
  qr          Generate QR codes for each account
  view        View the extracted accounts in the terminal
  help        Help about any command

Global Flags (for all commands):
  -i, --interactive       Interactive mode (prompt for input)
  -q, --qrimage string    Path to image containing Google Authenticator QR code
  -u, --uri string        Google Authenticator export URI

Flags for 'view' command:
  -p, --pretty            Enable pretty formatted output (default: true)
  -r, --show-qr           Display QR codes in the terminal
  -s, --show-secrets      Show full secrets (USE WITH CAUTION)

Flags for 'json' command:
  -f, --file string       Output file path for JSON (default: "accounts.json")
  -s, --save              Save to file (if false, prints to terminal) (default: true)

Flags for 'qr' command:
  -d, --dir string        Directory for saving QR code images (default: "qrcodes")
  -s, --save              Save to files (if false, displays in terminal) (default: true)
```

### Legacy Mode

For backward compatibility, you can still run the tool without a command:

```bash
gauth-extractor -u "otpauth-migration://offline?data=..."
```

This will run in interactive mode, prompting you to choose the output format.

## 📱 How to Export from Google Authenticator

1. **📲 Open** the Google Authenticator app
2. **⋮ Tap** the three dots menu and select "Transfer accounts"
3. **📤 Choose** "Export accounts"
4. **✅ Select** the accounts you want to export
5. **Choose one of these methods**:

   ### 📷 Method 1: Using a QR Scanner App

   - Scan the QR code using any QR scanner app
   - Copy the URI that looks like `otpauth-migration://offline?data=...`
   - Provide this URI to the tool using `-u` flag or interactive mode

   ### 📱 Method 2: Using a Screenshot

   - Take a screenshot of the QR code
   - Save the image file
   - Provide the image path to the tool using `-q` flag

## 🔑 Understanding Secret Formats

This tool extracts and presents secrets in two formats:

- **BASE32 (`totpSecret`)**: This is the format used by most authenticator apps and password managers. It typically appears as uppercase letters and numbers (A-Z, 2-7).
- **BASE64 (`secret`)**: This is the internal format used by Google Authenticator. It's usually shorter but less compatible with other apps.

**Which one should you use?**

- When manually adding accounts to other authenticator apps:

  - **Always use the `totpSecret` (BASE32) value**
  - This is the standard format expected by most apps

- When using QR codes generated by this tool:
  - The QR code already contains the correct format
  - Simply scan the QR code with your new authenticator app

## 🔒 Security Considerations

- **❌ Never** upload your Google Authenticator QR codes to online QR scanners
- **⚠️ Avoid** sharing the URI through insecure channels
- **🗑️ Delete** any screenshots or images containing QR codes after migration
- **🧹 Clear** your terminal history after viewing full secrets (`history -c` on most systems)
- **🔄 Consider** resetting your 2FA on critical accounts after migration
- **🔐 Secure** any JSON exports as they contain sensitive authentication secrets

## 📋 Data Format

The tool extracts the following data for each account:

```json
{
  "name": "example@gmail.com",
  "issuer": "Example Service",
  "secret": "BASE64_ENCODED_SECRET",
  "totpSecret": "BASE32_ENCODED_SECRET_FOR_OTHER_APPS",
  "type": "TOTP",
  "algorithm": "SHA1",
  "digits": "SIX",
  "counter": 0
}
```

## 🔄 Migration Guide

### To Authy

1. Extract your accounts:
   ```bash
   gauth-extractor view -u "otpauth-migration://offline?data=..." -s
   ```
2. In Authy:
   - Select "Add Account"
   - Choose "Enter code manually"
   - Enter account name and the BASE32 secret (totpSecret)
   - Select "6-digit" tokens (for most accounts)

### To Bitwarden

1. Extract your accounts:

   ```bash
   gauth-extractor json -u "otpauth-migration://offline?data=..." -s=false
   ```

2. In Bitwarden:
   - Create or edit a login entry
   - Scroll to the "Authenticator Key (TOTP)" section
   - Enter the BASE32 secret (totpSecret) value
   - Save the entry

### To 1Password

1. Generate individual QR codes:

   ```bash
   gauth-extractor qr -u "otpauth-migration://offline?data=..."
   ```

2. In 1Password:
   - Create or edit an item
   - Click "Add One-Time Password"
   - Select "Scan QR Code"
   - Capture each QR code generated by the tool

### To KeePass (with KeePassOTP plugin)

1. Extract your accounts:

   ```bash
   gauth-extractor view -u "otpauth-migration://offline?data=..." -s
   ```

2. In KeePass (with KeePassOTP plugin):
   - Edit an entry
   - Go to the "Additional" tab
   - Click "Set Up TOTP"
   - Enter the BASE32 secret (totpSecret)
   - Set other parameters as needed (6 digits, 30 seconds period)

## 🧪 Development

### 🔄 CI/CD Workflows

This project uses GitHub Actions for continuous integration and deployment:

- **🧪 CI**: Runs tests on PRs and pushes to the main branch
- **🧹 Lint**: Performs code linting with golangci-lint
- **🚀 Release Builder**: Manually triggered workflow to create releases

#### Creating a Release

To create a new release:

1. Go to the "Actions" tab in the GitHub repository
2. Select the "🚀 Release Builder" workflow
3. Click on "Run workflow"
4. Enter:
   - Version tag (e.g., `v1.0.0`)
   - Select whether it's a prerelease
5. Click "Run workflow"

This will:

- Run the test suite
- Build binaries for Linux, macOS (Intel and Apple Silicon), and Windows
- Create a Docker image and push it to GitHub Container Registry (ghcr.io)
- Create a GitHub release with the binaries attached

### Protocol Buffer

The tool uses Protocol Buffers to decode Google Authenticator's data format:

```protobuf
message MigrationPayload {
  enum Algorithm {
    ALGORITHM_UNSPECIFIED = 0;
    SHA1 = 1;
    SHA256 = 2;
    SHA512 = 3;
    MD5 = 4;
  }

  enum DigitCount {
    DIGIT_COUNT_UNSPECIFIED = 0;
    SIX = 1;
    EIGHT = 2;
    SEVEN = 3;
  }

  enum OtpType {
    OTP_TYPE_UNSPECIFIED = 0;
    HOTP = 1;
    TOTP = 2;
  }

  message OtpParameters {
    bytes secret = 1;
    string name = 2;
    string issuer = 3;
    Algorithm algorithm = 4;
    DigitCount digits = 5;
    OtpType type = 6;
    int64 counter = 7;
    string unique_id = 8;
  }

  repeated OtpParameters otp_parameters = 1;
  int32 version = 2;
  int32 batch_size = 3;
  int32 batch_index = 4;
  int32 batch_id = 5;
}
```

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 👏 Acknowledgments

- Inspired by [krissrex/google-authenticator-exporter](https://github.com/krissrex/google-authenticator-exporter)
- Protocol Buffer definition based on [beemdevelopment/Aegis](https://github.com/beemdevelopment/Aegis)
- Uses [makiuchi-d/gozxing](https://github.com/makiuchi-d/gozxing) for QR code decoding
- Uses [skip2/go-qrcode](https://github.com/skip2/go-qrcode) for QR code generation

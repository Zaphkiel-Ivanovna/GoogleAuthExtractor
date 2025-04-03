# Google Authenticator Secret Extractor

A Go tool to extract TOTP/HOTP secrets from Google Authenticator export QR codes.

> **Note:** This project is not affiliated with Google.

## Features

- **Strongly Typed**: Written in Go with proper error handling and type checking
- **Secure**: Process your 2FA secrets locally without any external services
- **Flexible Output**: Save as JSON or regenerate individual QR codes for each account
- **Easy Migration**: Move your accounts to other authenticator apps like Authy, Bitwarden, etc.

## Installation

### From Source

Requirements:

- Go 1.16 or higher
- Protocol Buffers compiler (`protoc`) for development

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

### Prebuilt Binaries

Download the latest release from [GitHub Releases](https://github.com/Zaphkiel-Ivanovna/GoogleAuthExtractor/releases).

## Usage

### Basic Command Line Usage

```bash
# Extract using interactive mode
gauth-extractor -i

# Extract to JSON file
gauth-extractor -u "otpauth-migration://offline?data=..." -o json -f accounts.json

# Generate QR codes
gauth-extractor -u "otpauth-migration://offline?data=..." -o qrcode -d ./qrcodes
```

### Command Line Options

```
Usage:
  gauth-extractor [flags]

Flags:
  -d, --dir string         Directory for QR codes (default "qrcodes")
  -f, --file string        Output file for JSON (default "accounts.json")
  -h, --help               help for gauth-extractor
  -i, --interactive        Interactive mode (prompt for URI)
  -o, --output string      Output type (json or qrcode) (default "json")
  -u, --uri string         Google Authenticator export URI
```

## How to Export from Google Authenticator

1. Open the Google Authenticator app
2. Tap the three dots menu (⋮) and select "Transfer accounts"
3. Choose "Export accounts"
4. Select the accounts you want to export
5. Scan the QR code using a QR code scanner app to get the URI
   - You can use apps like [ZXing Barcode Scanner](https://play.google.com/store/apps/details?id=com.google.zxing.client.android) for Android
   - You cannot take a screenshot of Google Authenticator, so use another camera/device
6. The scanned result should look like `otpauth-migration://offline?data=...`
7. Provide this URI to the tool

## Security Considerations

- **Never** upload your Google Authenticator QR codes to online QR scanners
- **Avoid** sharing the URI through insecure channels
- **Delete** any screenshots or images containing the QR codes after migration
- **Consider** resetting your 2FA on critical accounts after migration

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

- This project is inspired by [krissrex/google-authenticator-exporter](https://github.com/krissrex/google-authenticator-exporter)
- Protocol buffer definition based on [beemdevelopment/Aegis](https://github.com/beemdevelopment/Aegis)

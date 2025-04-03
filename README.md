# 🔐 Google Authenticator Extractor

Extract TOTP/HOTP secrets from Google Authenticator export QR codes with ease.

> **Note:** This project is not affiliated with Google.

## ✨ Features

- **🔒 Secure Processing**: Handle your 2FA secrets locally without external services
- **💪 Strongly Typed**: Written in Go with proper error handling
- **🖼️ QR Image Processing**: Extract directly from screenshots containing QR codes
- **📤 Flexible Output**:
  - 📄 Export to JSON for backup or custom processing
  - 🔄 Generate individual QR codes for each account to scan with other apps
- **🔄 Easy Migration**: Move your accounts to any authenticator app (Authy, Bitwarden, etc.)
- **🐳 Containerized**: Docker support for consistent execution

## 📦 Installation

### 📥 Prebuilt Binaries

Download the latest release from [GitHub Releases](https://github.com/Zaphkiel-Ivanovna/GoogleAuthExtractor/releases).

### 🛠️ From Source

Requirements:

- Go 1.16 or higher

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
# Build the Docker image
docker build -t gauth-extractor .

# Run in interactive mode
docker run -it --rm -v "$(pwd):/home/appuser/data" gauth-extractor
```

## 🧰 Usage

### 🖥️ Basic Command Line

```bash
# Extract using interactive mode
gauth-extractor -i

# Extract from URI to JSON
gauth-extractor -u "otpauth-migration://offline?data=..." -o json -f accounts.json

# Extract from URI to individual QR codes
gauth-extractor -u "otpauth-migration://offline?data=..." -o qrcode -d ./qrcodes

# Extract directly from a QR code image
gauth-extractor -p "/path/to/qrcode-screenshot.png" -o json
```

### 📋 Command Line Options

```
Usage:
  gauth-extractor [flags]

Flags:
  -d, --dir string         Directory for QR codes (default "qrcodes")
  -f, --file string        Output file for JSON (default "accounts.json")
  -h, --help               Help for gauth-extractor
  -i, --interactive        Interactive mode (prompt for URI)
  -p, --image string       Path to image containing Google Authenticator QR code
  -o, --output string      Output type (json or qrcode) (default "json")
  -u, --uri string         Google Authenticator export URI
```

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
   - Provide the image path to the tool using `-p` flag

## 🔒 Security Considerations

- **❌ Never** upload your Google Authenticator QR codes to online QR scanners
- **⚠️ Avoid** sharing the URI through insecure channels
- **🗑️ Delete** any screenshots or images containing QR codes after migration
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

> **Note**: When migrating to other apps, use the `totpSecret` value, not the `secret` value.

## 🔄 Migration Guide

### To Authy

1. Export to JSON or QR codes with `gauth-extractor`
2. For JSON: Manually add each account using the `totpSecret` value
3. For QR codes: Scan each generated QR code with Authy

### To Bitwarden

1. Export to JSON or QR codes
2. In Bitwarden, create a new login item or TOTP entry
3. For JSON: Enter the `totpSecret` as the Authenticator Key
4. For QR codes: Scan each QR code with Bitwarden

## 🧪 Development

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

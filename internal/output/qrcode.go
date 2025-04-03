// File: internal/output/qrcode.go
package output

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"regexp"

	"github.com/Zaphkiel-Ivanovna/GoogleAuthExtractor/internal/decoder"
	"github.com/fatih/color"
	"github.com/skip2/go-qrcode"
)

// SaveToQRCodes generates individual QR codes for each account
func SaveToQRCodes(accounts []decoder.Account, directory string) error {
	// Create directory if it doesn't exist
	if err := os.MkdirAll(directory, 0755); err != nil {
		return fmt.Errorf("failed to create directory '%s': %w", directory, err)
	}

	// Generate QR codes for each account
	for _, account := range accounts {
		// Generate otpauth URI
		uri := generateOtpAuthURI(account)

		// Create filename
		filename := generateFilename(account, directory)

		// Check if file already exists
		if _, err := os.Stat(filename); err == nil {
			color.Yellow("Warning: File '%s' already exists, skipping", filename)
			continue
		}

		// Generate QR code
		err := qrcode.WriteFile(uri, qrcode.Medium, 256, filename)
		if err != nil {
			color.Red("Error: Failed to create QR code for '%s': %v", account.Name, err)
			continue
		}

		color.Green("Created QR code: %s", filename)
	}

	return nil
}

// generateOtpAuthURI creates an otpauth URI for TOTP/HOTP
func generateOtpAuthURI(account decoder.Account) string {
	otpType := "totp"
	if account.Type == "HOTP" {
		otpType = "hotp"
	}

	label := url.PathEscape(account.Name)
	secret := url.QueryEscape(account.TOTPSecret)
	issuer := url.QueryEscape(account.Issuer)

	uri := fmt.Sprintf("otpauth://%s/%s?secret=%s", otpType, label, secret)

	if issuer != "" {
		uri += fmt.Sprintf("&issuer=%s", issuer)
	}

	if account.Type == "HOTP" {
		uri += fmt.Sprintf("&counter=%d", account.Counter)
	}

	return uri
}

// generateFilename creates a safe filename for the QR code
func generateFilename(account decoder.Account, directory string) string {
	issuer := account.Issuer
	if issuer == "" {
		issuer = "No_Issuer"
	}

	// Sanitize the filename
	name := sanitizeFilename(account.Name)
	issuer = sanitizeFilename(issuer)

	return filepath.Join(directory, fmt.Sprintf("%s (%s).png", issuer, name))
}

// sanitizeFilename removes characters that are not valid in filenames
func sanitizeFilename(name string) string {
	// Replace special characters with underscore
	re := regexp.MustCompile(`[\\/:*?"<>|%&{}$+!'=@]`)
	return re.ReplaceAllString(name, "_")
}

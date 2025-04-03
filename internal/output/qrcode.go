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

func SaveToQRCodes(accounts []decoder.Account, directory string) error {

	if err := os.MkdirAll(directory, 0755); err != nil {
		return fmt.Errorf("failed to create directory '%s': %w", directory, err)
	}

	for _, account := range accounts {

		uri := generateOtpAuthURI(account)

		filename := generateFilename(account, directory)

		if _, err := os.Stat(filename); err == nil {
			color.Yellow("Warning: File '%s' already exists, skipping", filename)
			continue
		}

		err := qrcode.WriteFile(uri, qrcode.Medium, 256, filename)
		if err != nil {
			color.Red("Error: Failed to create QR code for '%s': %v", account.Name, err)
			continue
		}

		color.Green("Created QR code: %s", filename)
	}

	return nil
}

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

func generateFilename(account decoder.Account, directory string) string {
	issuer := account.Issuer
	if issuer == "" {
		issuer = "No_Issuer"
	}

	name := sanitizeFilename(account.Name)
	issuer = sanitizeFilename(issuer)

	return filepath.Join(directory, fmt.Sprintf("%s (%s).png", issuer, name))
}

func sanitizeFilename(name string) string {

	re := regexp.MustCompile(`[\\/:*?"<>|%&{}$+!'=@]`)
	return re.ReplaceAllString(name, "_")
}

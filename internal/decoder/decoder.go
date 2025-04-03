// File: internal/decoder/decoder.go
package decoder

import (
	"encoding/base32"
	"encoding/base64"
	"fmt"
	"net/url"
	"strings"

	"github.com/Zaphkiel-Ivanovna/GoogleAuthExtractor/internal/proto"
	pb "google.golang.org/protobuf/proto"
)

// Account represents a decoded TOTP/HOTP account
type Account struct {
	Name       string `json:"name"`
	Issuer     string `json:"issuer,omitempty"`
	Secret     string `json:"secret"`            // Base64 encoded secret
	TOTPSecret string `json:"totpSecret"`        // Base32 encoded secret (for TOTP apps)
	Type       string `json:"type"`              // "TOTP" or "HOTP"
	Algorithm  string `json:"algorithm"`         // "SHA1", "SHA256", etc.
	Digits     string `json:"digits"`            // "SIX", "EIGHT", etc.
	Counter    int64  `json:"counter,omitempty"` // Only for HOTP
}

// DecodeExportURI decodes a Google Authenticator export URI
func DecodeExportURI(uri string) ([]Account, error) {
	// Parse the URI
	parsedURL, err := url.Parse(uri)
	if err != nil {
		return nil, fmt.Errorf("invalid URI format: %w", err)
	}

	// Check for correct URI format
	if !strings.HasPrefix(parsedURL.Scheme, "otpauth-migration") {
		return nil, fmt.Errorf("invalid URI scheme: expected 'otpauth-migration', got '%s'", parsedURL.Scheme)
	}

	// Extract the data parameter
	queryParams := parsedURL.Query()
	dataParam := queryParams.Get("data")
	if dataParam == "" {
		return nil, fmt.Errorf("missing 'data' parameter in URI")
	}

	// URL decode and then base64 decode
	decodedData, err := url.QueryUnescape(dataParam)
	if err != nil {
		return nil, fmt.Errorf("failed to URL-decode data parameter: %w", err)
	}

	rawData, err := base64.StdEncoding.DecodeString(decodedData)
	if err != nil {
		return nil, fmt.Errorf("failed to base64-decode data: %w", err)
	}
	// Unmarshal protobuf
	payload := &proto.MigrationPayload{}
	err = pb.Unmarshal(rawData, payload)
	if err != nil {
		return nil, fmt.Errorf("failed to decode protobuf data: %w", err)
	}

	// Check payload version
	if payload.Version != 1 {
		fmt.Printf("Warning: Expected payload version 1, but got %d. This might cause issues.\n", payload.Version)
	}

	// Convert to account structs
	accounts := make([]Account, 0, len(payload.OtpParameters))
	for _, otpParams := range payload.OtpParameters {
		account := Account{
			Name:       otpParams.Name,
			Issuer:     otpParams.Issuer,
			Secret:     base64.StdEncoding.EncodeToString(otpParams.Secret),
			TOTPSecret: toBase32(otpParams.Secret),
			Type:       otpParams.Type.String(),
			Algorithm:  otpParams.Algorithm.String(),
			Digits:     otpParams.Digits.String(),
		}

		if otpParams.Type == proto.MigrationPayload_HOTP {
			account.Counter = otpParams.Counter
		}

		accounts = append(accounts, account)
	}

	return accounts, nil
}

// toBase32 converts a byte array to Base32 encoded string (RFC 4648)
// This is what most TOTP authenticator apps use for the secret key
func toBase32(data []byte) string {
	// Standard base32 encoding with padding (=)
	encoder := base32.StdEncoding.WithPadding(base32.NoPadding)
	return encoder.EncodeToString(data)
}

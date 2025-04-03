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

type Account struct {
	Name       string `json:"name"`
	Issuer     string `json:"issuer,omitempty"`
	Secret     string `json:"secret"`
	TOTPSecret string `json:"totpSecret"`
	Type       string `json:"type"`
	Algorithm  string `json:"algorithm"`
	Digits     string `json:"digits"`
	Counter    int64  `json:"counter,omitempty"`
}

func DecodeExportURI(uri string) ([]Account, error) {

	parsedURL, err := url.Parse(uri)
	if err != nil {
		return nil, fmt.Errorf("invalid URI format: %w", err)
	}

	if !strings.HasPrefix(parsedURL.Scheme, "otpauth-migration") {
		return nil, fmt.Errorf("invalid URI scheme: expected 'otpauth-migration', got '%s'", parsedURL.Scheme)
	}

	queryParams := parsedURL.Query()
	dataParam := queryParams.Get("data")
	if dataParam == "" {
		return nil, fmt.Errorf("missing 'data' parameter in URI")
	}

	decodedData, err := url.QueryUnescape(dataParam)
	if err != nil {
		return nil, fmt.Errorf("failed to URL-decode data parameter: %w", err)
	}

	rawData, err := base64.StdEncoding.DecodeString(decodedData)
	if err != nil {
		return nil, fmt.Errorf("failed to base64-decode data: %w", err)
	}

	payload := &proto.MigrationPayload{}
	err = pb.Unmarshal(rawData, payload)
	if err != nil {
		return nil, fmt.Errorf("failed to decode protobuf data: %w", err)
	}

	if payload.Version != 1 {
		fmt.Printf("Warning: Expected payload version 1, but got %d. This might cause issues.\n", payload.Version)
	}

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

func toBase32(data []byte) string {

	encoder := base32.StdEncoding.WithPadding(base32.NoPadding)
	return encoder.EncodeToString(data)
}

package decoder

import (
	"testing"
)

func TestDecodeExportURI(t *testing.T) {
	tests := []struct {
		name          string
		uri           string
		expectedCount int
		expectError   bool
	}{
		{
			name:          "Valid single QR code",
			uri:           "otpauth-migration://offline?data=CiIKCkhlbGwPId6tvugSDlRlc3QgYWNjb3VudCAxIAEoATACCiIKCgBlbGxvId6tvu8SDlRlc3QgYWNjb3VudCAyIAEoATACCiMKCgBEjWxkLzvjHR8SDUNvdW50ZXIga2V5IDEgASgBMAE4ARABGAEgACj8nJf4Bg%3D%3D",
			expectedCount: 3,
			expectError:   false,
		},
		{
			name:          "Valid QR code with SHA512 and 8 digits",
			uri:           "otpauth-migration://offline?data=CjoKFD1jwRTgK6xTGKA0gdTWaGMebxmTEg1UT1RQZ2VuZXJhdG9yGg1UT1RQZ2VuZXJhdG9yIAMoAjACEAEYASAAKJ2G1g0%3D",
			expectedCount: 1,
			expectError:   false,
		},
		{
			name:          "Invalid URI scheme",
			uri:           "https://example.com",
			expectedCount: 0,
			expectError:   true,
		},
		{
			name:          "Missing data parameter",
			uri:           "otpauth-migration://offline",
			expectedCount: 0,
			expectError:   true,
		},
		{
			name:          "Invalid data (not base64)",
			uri:           "otpauth-migration://offline?data=not-base64",
			expectedCount: 0,
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			accounts, err := DecodeExportURI(tt.uri)

			if tt.expectError && err == nil {
				t.Errorf("Expected error but got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Did not expect error but got: %v", err)
			}

			if tt.expectError {
				return
			}

			if len(accounts) != tt.expectedCount {
				t.Errorf("Expected %d accounts, got %d", tt.expectedCount, len(accounts))
			}

			if tt.name == "Valid single QR code" && len(accounts) >= 3 {

				if accounts[0].Name != "Test account 1" {
					t.Errorf("Expected name 'Test account 1', got '%s'", accounts[0].Name)
				}
				if accounts[0].Type != "TOTP" {
					t.Errorf("Expected type 'TOTP', got '%s'", accounts[0].Type)
				}

				if accounts[2].Name != "Counter key 1" {
					t.Errorf("Expected name 'Counter key 1', got '%s'", accounts[2].Name)
				}
				if accounts[2].Type != "HOTP" {
					t.Errorf("Expected type 'HOTP', got '%s'", accounts[2].Type)
				}
				if accounts[2].Counter != 1 {
					t.Errorf("Expected counter 1, got %d", accounts[2].Counter)
				}
			}

			if tt.name == "Valid QR code with SHA512 and 8 digits" && len(accounts) >= 1 {
				if accounts[0].Algorithm != "SHA512" {
					t.Errorf("Expected algorithm 'SHA512', got '%s'", accounts[0].Algorithm)
				}
				if accounts[0].Digits != "EIGHT" {
					t.Errorf("Expected digits 'EIGHT', got '%s'", accounts[0].Digits)
				}
				if accounts[0].Issuer != "TOTPgenerator" {
					t.Errorf("Expected issuer 'TOTPgenerator', got '%s'", accounts[0].Issuer)
				}
			}
		})
	}
}

func TestInvalidURIs(t *testing.T) {
	invalidURIs := []string{
		"",
		"not-a-uri",
		"otpauth-migration://offline?data=",
		"otpauth-migration://offline?otherParam=123",
	}

	for _, uri := range invalidURIs {
		_, err := DecodeExportURI(uri)
		if err == nil {
			t.Errorf("Expected error for URI '%s' but got nil", uri)
		}
	}
}

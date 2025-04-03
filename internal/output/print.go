package output

import (
	"fmt"
	"strings"
	"time"

	"github.com/Zaphkiel-Ivanovna/GoogleAuthExtractor/internal/decoder"
	"github.com/fatih/color"
)

// PrettyPrintAccounts prints account information in a formatted way in the terminal
func PrettyPrintAccounts(accounts []decoder.Account, pretty bool, showFullSecrets bool) {
	if !pretty {
		// Simple table format if not pretty
		fmt.Println("+-----------------------+-----------------------+----------+----------+")
		fmt.Println("| Name                  | Issuer                | Type     | Digits   |")
		fmt.Println("+-----------------------+-----------------------+----------+----------+")
		
		for _, account := range accounts {
			name := truncateString(account.Name, 21)
			issuer := truncateString(account.Issuer, 21)
			fmt.Printf("| %-21s | %-21s | %-8s | %-8s |\n", 
				name, issuer, account.Type, account.Digits)
		}
		
		fmt.Println("+-----------------------+-----------------------+----------+----------+")
		return
	}

	// Pretty detailed output with colors
	cyan := color.New(color.FgCyan, color.Bold).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	red := color.New(color.FgHiRed).SprintFunc()

	fmt.Printf("%s\n\n", cyan("🔐 Google Authenticator Accounts"))
	fmt.Printf("Extracted at: %s\n", green(time.Now().Format("2006-01-02 15:04:05")))
	fmt.Printf("Total accounts: %s\n\n", green(len(accounts)))

	for i, account := range accounts {
		if i > 0 {
			fmt.Println(strings.Repeat("─", 50))
		}
		
		fmt.Printf("Account #%d:\n", i+1)
		
		fmt.Printf("  %s: %s\n", cyan("Name"), account.Name)
		
		if account.Issuer != "" {
			fmt.Printf("  %s: %s\n", cyan("Issuer"), account.Issuer)
		}
		
		fmt.Printf("  %s: %s\n", cyan("Type"), account.Type)
		fmt.Printf("  %s: %s\n", cyan("Algorithm"), account.Algorithm)
		fmt.Printf("  %s: %s\n", cyan("Digits"), account.Digits)
		
		if account.Type == "HOTP" {
			fmt.Printf("  %s: %d\n", cyan("Counter"), account.Counter)
		}
		
		// Display secrets based on user preference
		if showFullSecrets {
			fmt.Printf("  %s: %s\n", cyan("Secret (BASE32)"), account.TOTPSecret)
			fmt.Printf("  %s: %s\n", cyan("Secret (BASE64)"), account.Secret)
			
			// Warning about displayed secrets
			fmt.Printf("\n  %s\n", red("⚠️  Warning: Full secrets are displayed. Clear your terminal history when done."))
		} else {
			// Only show truncated version of secrets for security
			secretLen := len(account.TOTPSecret)
			if secretLen > 4 {
				visiblePart := account.TOTPSecret[0:4]
				hiddenPart := strings.Repeat("*", secretLen-4)
				fmt.Printf("  %s: %s%s\n", cyan("Secret"), visiblePart, hiddenPart)
			} else {
				fmt.Printf("  %s: %s\n", cyan("Secret"), "****")
			}
			
			// Option for displaying full secret
			fmt.Printf("\n  To view full secret: %s\n", yellow("Use 'gauth-extractor view -s' or 'gauth-extractor json -s=false'"))
		}
	}
	
	fmt.Println()
	color.Yellow("Note: When adding accounts to other authenticator apps, use the 'totpSecret' value as the secret key, not the 'secret' value.")
}

// truncateString trims a string if it's longer than maxLen and adds "..." at the end
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/Zaphkiel-Ivanovna/GoogleAuthExtractor/internal/decoder"
	"github.com/Zaphkiel-Ivanovna/GoogleAuthExtractor/internal/input"
	"github.com/Zaphkiel-Ivanovna/GoogleAuthExtractor/internal/output"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	qrImagePath      string
	uri              string
	interactiveInput bool

	jsonFile       string
	qrCodesDir     string
	displayPretty  bool
	displayQR      bool
	saveToFiles    bool
	showFullSecret bool
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "gauth-extractor",
		Short: "Extract TOTP/HOTP secrets from Google Authenticator export",
		Long: `🔐 Google Authenticator Secret Extractor

Extract TOTP/HOTP secrets from Google Authenticator by decoding 
the QR codes generated during account export.

This tool allows you to extract secrets from the Google Authenticator app
and migrate them to other authenticator apps or password managers.

To use:
1. Export accounts from Google Authenticator app
2. Use one of the following methods:
   a) Scan the QR code with a QR code scanner app to get the URI text
   b) Take a screenshot of the QR code and provide the image path to this tool
   c) Run in interactive mode and follow the prompts`,
	}

	viewCmd := &cobra.Command{
		Use:   "view",
		Short: "View the extracted accounts in the terminal",
		Long: `View the extracted accounts directly in the terminal

This command displays the accounts in the terminal without saving them to files.
You can customize the display using the --pretty and --qr flags.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			accounts, err := getAccounts(args)
			if err != nil {
				return err
			}

			output.PrettyPrintAccounts(accounts, displayPretty, showFullSecret)

			if displayQR {
				err = output.DisplayQRCodesInTerminal(accounts)
				if err != nil {
					return fmt.Errorf("failed to display QR codes in terminal: %w", err)
				}
			}

			return nil
		},
	}

	jsonCmd := &cobra.Command{
		Use:   "json",
		Short: "Export accounts to JSON format",
		Long: `Export accounts to JSON format

This command exports the extracted accounts to a JSON file
which can be used for backup or for importing into other applications.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			accounts, err := getAccounts(args)
			if err != nil {
				return err
			}

			if saveToFiles {
				err = output.SaveToJSON(accounts, jsonFile)
				if err != nil {
					return fmt.Errorf("failed to save JSON: %w", err)
				}
				return nil
			}

			output.PrintJSON(accounts)
			return nil
		},
	}

	qrCmd := &cobra.Command{
		Use:   "qr",
		Short: "Generate QR codes for each account",
		Long: `Generate QR codes for each account

This command generates individual QR code images for each extracted account.
These QR codes can be scanned by other authenticator apps.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			accounts, err := getAccounts(args)
			if err != nil {
				return err
			}

			if saveToFiles {
				err = output.SaveToQRCodes(accounts, qrCodesDir)
				if err != nil {
					return fmt.Errorf("failed to generate QR codes: %w", err)
				}
				return nil
			}

			err = output.DisplayQRCodesInTerminal(accounts)
			if err != nil {
				return fmt.Errorf("failed to display QR codes in terminal: %w", err)
			}
			return nil
		},
	}

	rootCmd.PersistentFlags().StringVarP(&uri, "uri", "u", "", "Google Authenticator export URI (otpauth-migration://...)")
	rootCmd.PersistentFlags().StringVarP(&qrImagePath, "qrimage", "q", "", "Path to image containing Google Authenticator QR code")
	rootCmd.PersistentFlags().BoolVarP(&interactiveInput, "interactive", "i", false, "Interactive mode (prompt for input)")

	viewCmd.Flags().BoolVarP(&displayPretty, "pretty", "p", true, "Enable pretty formatted output (colorful and detailed)")
	viewCmd.Flags().BoolVarP(&displayQR, "show-qr", "r", false, "Display QR codes in the terminal")
	viewCmd.Flags().BoolVarP(&showFullSecret, "show-secrets", "s", false, "Show full secrets (USE WITH CAUTION)")

	jsonCmd.Flags().StringVarP(&jsonFile, "file", "f", "accounts.json", "Output file path for JSON")
	jsonCmd.Flags().BoolVarP(&saveToFiles, "save", "s", true, "Save to file (if false, prints to terminal)")

	qrCmd.Flags().StringVarP(&qrCodesDir, "dir", "d", "qrcodes", "Directory for saving QR code images")
	qrCmd.Flags().BoolVarP(&saveToFiles, "save", "s", true, "Save to files (if false, displays in terminal)")

	rootCmd.AddCommand(viewCmd, jsonCmd, qrCmd)

	rootCmd.RunE = func(cmd *cobra.Command, args []string) error {
		if cmd.CalledAs() == "gauth-extractor" {
			fmt.Println("Running in legacy mode...")
			return handleLegacyCommand(args)
		}
		return cmd.Help()
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func getAccounts(args []string) ([]decoder.Account, error) {
	extractedURI := ""

	if qrImagePath != "" {
		var err error
		extractedURI, err = input.ExtractQRCodeFromImage(qrImagePath)
		if err != nil {
			return nil, fmt.Errorf("failed to extract QR code from image: %w", err)
		}
		color.Green("Successfully extracted QR code from image")
	} else if uri != "" {

		extractedURI = uri
	} else {

		if !interactiveInput && len(args) == 0 {
			interactiveInput = true
		} else if len(args) > 0 {
			extractedURI = args[0]
		}

		if interactiveInput {
			extractedURI = promptURI()
		}
	}

	if extractedURI == "" {
		return nil, fmt.Errorf("no URI provided. Use --uri, --qrimage, or --interactive flags")
	}

	accounts, err := decoder.DecodeExportURI(extractedURI)
	if err != nil {
		return nil, fmt.Errorf("failed to decode URI: %w", err)
	}

	color.Green("Successfully decoded %d accounts", len(accounts))
	return accounts, nil
}

func handleLegacyCommand(args []string) error {
	accounts, err := getAccounts(args)
	if err != nil {
		return err
	}

	fmt.Println("\nHow would you like to output the accounts?")
	fmt.Println("1. Save to JSON file")
	fmt.Println("2. Print JSON to terminal")
	fmt.Println("3. Generate QR code files")
	fmt.Println("4. Display in terminal")
	fmt.Print("\nEnter option (1-4): ")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	choice := strings.TrimSpace(scanner.Text())

	switch choice {
	case "1":
		fmt.Printf("Save to JSON file '%s'? [y/N]: ", jsonFile)
		scanner.Scan()
		if strings.HasPrefix(strings.ToLower(scanner.Text()), "y") {
			err = output.SaveToJSON(accounts, jsonFile)
			if err != nil {
				return fmt.Errorf("failed to save JSON: %w", err)
			}
		}
	case "2":
		output.PrintJSON(accounts)
	case "3":
		fmt.Printf("Save QR codes to directory '%s'? [y/N]: ", qrCodesDir)
		scanner.Scan()
		if strings.HasPrefix(strings.ToLower(scanner.Text()), "y") {
			err = output.SaveToQRCodes(accounts, qrCodesDir)
			if err != nil {
				return fmt.Errorf("failed to generate QR codes: %w", err)
			}
		}
	case "4":
		fmt.Print("Use pretty formatting? [Y/n]: ")
		scanner.Scan()
		pretty := !strings.HasPrefix(strings.ToLower(scanner.Text()), "n")

		fmt.Print("Show QR codes in terminal? [y/N]: ")
		scanner.Scan()
		showQR := strings.HasPrefix(strings.ToLower(scanner.Text()), "y")

		fmt.Print("Show full secrets? (CAUTION: Secrets will be visible) [y/N]: ")
		scanner.Scan()
		showSecrets := strings.HasPrefix(strings.ToLower(scanner.Text()), "y")

		output.PrettyPrintAccounts(accounts, pretty, showSecrets)
		if showQR {
			err = output.DisplayQRCodesInTerminal(accounts)
			if err != nil {
				return fmt.Errorf("failed to display QR codes in terminal: %w", err)
			}
		}
	default:
		return fmt.Errorf("invalid choice")
	}

	return nil
}

func promptURI() string {
	color.Red("WARNING: By using online QR decoders or untrusted ways of transferring the URI text,")
	color.Red("you risk someone storing the QR code or URI text and stealing your 2FA codes!")
	color.Red("Remember that the data contains the website, your email and the 2FA code!")
	fmt.Println()

	fmt.Println("Enter the URI from Google Authenticator QR code.")
	fmt.Println("The URI looks like otpauth-migration://offline?data=...")
	fmt.Println("")
	fmt.Println("You can get it by exporting from Google Authenticator app, then scanning the QR with")
	fmt.Println("a QR code scanner app, and copying the text to your computer.")
	fmt.Println("")

	fmt.Print("Enter URI: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
}

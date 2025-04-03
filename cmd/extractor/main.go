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
	outputType      string
	outputFile      string
	qrCodeDir       string
	uri             string
	imagePath       string
	interactiveMode bool
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "gauth-extractor",
		Short: "Extract TOTP/HOTP secrets from Google Authenticator export",
		Long: `Google Authenticator Secret Extractor

This tool allows you to extract TOTP/HOTP secrets from Google Authenticator by
decoding the QR codes generated during account export.

To use:
1. Export accounts from Google Authenticator app
2. Use one of the following methods:
   a) Scan the QR code to obtain the "otpauth-migration://offline?data=..." URI and provide it to this tool
   b) Take a screenshot of the QR code and provide the image path to this tool using the --image flag`,
		RunE: runExtractor,
	}

	rootCmd.Flags().StringVarP(&outputType, "output", "o", "json", "Output type (json or qrcode)")
	rootCmd.Flags().StringVarP(&outputFile, "file", "f", "accounts.json", "Output file for JSON")
	rootCmd.Flags().StringVarP(&qrCodeDir, "dir", "d", "qrcodes", "Directory for QR codes")
	rootCmd.Flags().StringVarP(&uri, "uri", "u", "", "Google Authenticator export URI")
	rootCmd.Flags().StringVarP(&imagePath, "image", "p", "", "Path to image containing Google Authenticator QR code")
	rootCmd.Flags().BoolVarP(&interactiveMode, "interactive", "i", false, "Interactive mode (prompt for URI)")

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func runExtractor(cmd *cobra.Command, args []string) error {

	if imagePath != "" {
		extractedURI, err := input.ExtractQRCodeFromImage(imagePath)
		if err != nil {
			return fmt.Errorf("failed to extract QR code from image: %w", err)
		}

		color.Green("Successfully extracted QR code from image")
		uri = extractedURI
	} else {

		if uri == "" {
			if !interactiveMode && len(args) == 0 {
				interactiveMode = true
			} else if len(args) > 0 {
				uri = args[0]
			}
		}

		if interactiveMode {
			uri = promptURI()
		}
	}

	if uri == "" {
		return fmt.Errorf("no URI provided")
	}

	accounts, err := decoder.DecodeExportURI(uri)
	if err != nil {
		return fmt.Errorf("failed to decode URI: %w", err)
	}

	color.Green("Successfully decoded %d accounts", len(accounts))

	if outputType == "json" {
		if promptSaveFile() {
			err = output.SaveToJSON(accounts, outputFile)
			if err != nil {
				return fmt.Errorf("failed to save JSON: %w", err)
			}
		} else {
			output.PrintJSON(accounts)
		}
	} else if outputType == "qrcode" {
		err = output.SaveToQRCodes(accounts, qrCodeDir)
		if err != nil {
			return fmt.Errorf("failed to generate QR codes: %w", err)
		}
	} else {
		return fmt.Errorf("invalid output type: %s", outputType)
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

func promptSaveFile() bool {
	fmt.Printf("Save to file '%s'? [y/N]: ", outputFile)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	response := strings.ToLower(scanner.Text())
	return strings.HasPrefix(response, "y")
}

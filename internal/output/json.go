package output

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Zaphkiel-Ivanovna/GoogleAuthExtractor/internal/decoder"
	"github.com/fatih/color"
)

func SaveToJSON(accounts []decoder.Account, filename string) error {

	if _, err := os.Stat(filename); err == nil {
		return fmt.Errorf("file '%s' already exists", filename)
	}

	dir := filepath.Dir(filename)
	if dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory '%s': %w", dir, err)
		}
	}

	data, err := json.MarshalIndent(accounts, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal account data: %w", err)
	}

	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write to file '%s': %w", filename, err)
	}

	color.Green("Successfully saved %d accounts to %s", len(accounts), filename)
	return nil
}

func PrintJSON(accounts []decoder.Account) {

	data, err := json.MarshalIndent(accounts, "", "  ")
	if err != nil {
		color.Red("Error: Failed to marshal account data: %v", err)
		return
	}

	fmt.Println(string(data))
	color.Yellow("What you want to use as secret key in other password managers is 'totpSecret', not 'secret'!")
}

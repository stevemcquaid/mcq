package ai

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/stevemcquaid/mcq/pkg/errors"
	"github.com/stevemcquaid/mcq/pkg/logger"
)

// displayAndCopyResult copies the user story to clipboard and displays it
func displayAndCopyResult(userStory string) error {
	fmt.Println("\n📋 Copying to clipboard...")

	if err := CopyToClipboard(userStory); err != nil {
		userErr := errors.WrapError(err, "Clipboard copy failed")
		userErr.Display()
		// Don't fail the entire operation if clipboard copy fails
		fmt.Println("\n" + strings.Repeat("=", 60))
		fmt.Println(userStory)
		fmt.Println(strings.Repeat("=", 60))
		return nil
	}

	fmt.Println("✅ User story generated and copied to clipboard!")
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println(userStory)
	fmt.Println(strings.Repeat("=", 60))
	return nil
}

// CopyToClipboard copies text to clipboard on macOS
func CopyToClipboard(text string) error {
	logger.LogDetailed("Copying to clipboard", "length", len(text))
	cmd := exec.Command("pbcopy")
	cmd.Stdin = strings.NewReader(text)

	err := cmd.Run()
	if err != nil {
		logger.LogError("clipboard copy", err)
		return err
	}

	logger.LogBasic("Successfully copied to clipboard")
	return nil
}

// showConnectionProgress displays progress during API connection setup
func showConnectionProgress(provider, modelName string) {
	fmt.Printf("🔌 Connecting to %s API (%s)...\n", provider, modelName)
}

// showStreamingProgress displays progress when streaming starts
func showStreamingProgress() {
	fmt.Print("💭 ")
}

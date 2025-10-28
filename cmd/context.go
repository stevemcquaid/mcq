package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/stevemcquaid/mcq/pkg/ai"
)

// contextCmd represents the context command
var contextCmd = &cobra.Command{
	Use:   "context",
	Short: "Manage and test repository context gathering",
	Long:  `Commands for managing and testing repository context gathering for AI features.`,
}

// contextTestCmd represents the context test command
var contextTestCmd = &cobra.Command{
	Use:   "test",
	Short: "Test context gathering",
	Long: `Test the context gathering functionality to verify it's working correctly.

This command will gather all available context from your repository and display
it to help you verify that context gathering is working as expected.

Context includes:
  • Go module information (project name, version, dependencies)
  • README content (from root and docs/ directory)
  • Recent commit messages
  • Directory structure
  • Configuration files (go.mod, Makefile, Dockerfile, etc.)`,
	Run: func(cmd *cobra.Command, args []string) {
		testContext()
	},
}

func init() {
	RootCmd.AddCommand(contextCmd)
	contextCmd.AddCommand(contextTestCmd)
}

// testContext tests the context gathering functionality
func testContext() {
	fmt.Println("🧪 Testing Context Gathering")
	fmt.Println("============================")
	fmt.Println()

	// Create context config with auto-detect enabled
	config := ai.ContextConfig{
		AutoDetect:       true,
		IncludeReadme:    true,
		IncludeGoMod:     true,
		IncludeCommits:   true,
		IncludeStructure: true,
		IncludeConfigs:   true,
		MaxCommits:       10,
		MaxFileSize:      50 * 1024, // 50KB
	}

	// Gather context
	fmt.Println("📥 Gathering repository context...")
	ctx := ai.GatherContextIfNeeded(config)

	if ctx == nil {
		fmt.Println()
		fmt.Println("❌ Context gathering failed or returned no data")
		fmt.Println()
		fmt.Println("🔍 Diagnostic Information:")
		fmt.Println("-------------------------")
		fmt.Println()

		// Check for common issues
		fmt.Println("Checking for common issues...")
		fmt.Println()

		// Check if we're in a git repository
		if isInGitRepo() {
			fmt.Println("✅ In a Git repository")
		} else {
			fmt.Println("⚠️  Not in a Git repository (Git history not available)")
		}
		fmt.Println()

		// Check for go.mod
		if hasGoMod() {
			fmt.Println("✅ go.mod file found")
		} else {
			fmt.Println("⚠️  go.mod file not found")
		}
		fmt.Println()

		// Check for README (including docs directory)
		if hasReadme() {
			fmt.Println("✅ README file found (checked root and docs/)")
		} else {
			fmt.Println("⚠️  README file not found (checked root and docs/)")
		}
		fmt.Println()

		// Suggest solutions
		fmt.Println("💡 Recommendations:")
		fmt.Println("  • Try running from the repository root directory")
		fmt.Println("  • Verify files are readable (check permissions)")
		fmt.Println("  • Use 'mcq ai jira --no-context' to skip context")
		fmt.Println()
		return
	}

	fmt.Println("✅ Context gathered successfully!")
	fmt.Println()

	// Display results
	fmt.Println("📊 Context Summary:")
	fmt.Println("-------------------")
	fmt.Println()

	// Project Information
	if ctx.ProjectName != "" {
		fmt.Printf("📦 Project: %s\n", ctx.ProjectName)
	}
	if ctx.ModulePath != "" {
		fmt.Printf("   Module Path: %s\n", ctx.ModulePath)
	}
	if ctx.GoVersion != "" {
		fmt.Printf("   Go Version: %s\n", ctx.GoVersion)
	}
	if ctx.ProjectType != "" {
		fmt.Printf("   Project Type: %s\n", ctx.ProjectType)
	}

	if len(ctx.Dependencies) > 0 {
		fmt.Printf("   Dependencies: %d found\n", len(ctx.Dependencies))
		if len(ctx.Dependencies) <= 5 {
			for _, dep := range ctx.Dependencies {
				fmt.Printf("      - %s\n", dep)
			}
		} else {
			for _, dep := range ctx.Dependencies[:5] {
				fmt.Printf("      - %s\n", dep)
			}
			fmt.Printf("      ... and %d more\n", len(ctx.Dependencies)-5)
		}
	}
	fmt.Println()

	// README (including docs directory if present)
	if ctx.Readme != "" {
		hasDocsSection := strings.Contains(ctx.Readme, "## Documentation")

		// Count how many doc files are included
		docsCount := strings.Count(ctx.Readme, "\n\n###")

		if hasDocsSection || docsCount > 0 {
			var statusMsg string
			if docsCount > 0 {
				statusMsg = fmt.Sprintf("📄 README (includes root README + %d docs files from docs/ folder)", docsCount)
			} else if hasDocsSection {
				statusMsg = "📄 README (includes root README + docs/README.md)"
			} else {
				statusMsg = "📄 README"
			}
			fmt.Println(statusMsg)
		} else {
			fmt.Println("📄 README:")
		}

		// Show more content (up to 500 characters to see structure)
		readmePreview := ctx.Readme
		if len(readmePreview) > 500 {
			readmePreview = readmePreview[:500] + "..."
		}
		fmt.Println("   " + readmePreview)

		// Show summary if there's more content
		if len(ctx.Readme) > 500 {
			totalChars := len(ctx.Readme)
			fmt.Printf("   [... %d more characters of documentation]\n", totalChars-500)
		}
		fmt.Println()
	}

	// Recent Commits
	if len(ctx.RecentCommits) > 0 {
		fmt.Printf("📝 Recent Commits (%d):\n", len(ctx.RecentCommits))
		for i, commit := range ctx.RecentCommits {
			if i >= 5 { // Show max 5 commits
				fmt.Printf("   ... and %d more commits\n", len(ctx.RecentCommits)-5)
				break
			}
			if len(commit) > 80 {
				commit = commit[:80] + "..."
			}
			fmt.Printf("   %d. %s\n", i+1, commit)
		}
		fmt.Println()
	}

	// Directory Structure
	if ctx.DirectoryStructure != "" {
		structurePreview := ctx.DirectoryStructure
		lines := []rune(structurePreview)
		lineCount := 0
		lastIndex := 0
		for i, r := range lines {
			if r == '\n' {
				lineCount++
				if lineCount > 20 {
					lastIndex = i
					break
				}
			}
		}

		if lastIndex > 0 {
			structurePreview = string(lines[:lastIndex]) + "\n   ... (truncated)"
		}

		fmt.Println("📁 Directory Structure:")
		fmt.Println(structurePreview)
		fmt.Println()
	}

	// Configuration Files
	if len(ctx.ConfigFiles) > 0 {
		fmt.Printf("⚙️ Configuration Files (%d):\n", len(ctx.ConfigFiles))
		for filename, content := range ctx.ConfigFiles {
			contentPreview := content
			if len(contentPreview) > 150 {
				contentPreview = contentPreview[:150] + "..."
			}
			fmt.Printf("   • %s\n", filename)
			fmt.Printf("     %s\n", contentPreview)
			fmt.Println()
		}
	}

	fmt.Println("✅ Context test completed successfully!")
	fmt.Println()
	fmt.Println("💡 This context is used to improve AI-generated user stories")
	fmt.Println("   by providing relevant information about your repository.")
}

// Helper functions for diagnostics

// isInGitRepo checks if we're in a Git repository
func isInGitRepo() bool {
	// Try to run git rev-parse to check if we're in a git repo
	// This is a simpler check that just looks for .git directory
	return hasFileOrDir(".git")
}

// hasGoMod checks if go.mod file exists
func hasGoMod() bool {
	return hasFile("go.mod")
}

// hasReadme checks if a README file exists (with various extensions, including docs directory)
func hasReadme() bool {
	readmeFiles := []string{
		"README", "README.md", "README.txt", "README.rst",
		"docs/README", "docs/README.md", "docs/README.txt", "docs/README.rst",
	}
	for _, filename := range readmeFiles {
		if hasFile(filename) {
			return true
		}
	}
	return false
}

// hasFile checks if a file exists
func hasFile(filename string) bool {
	// Use a simple file existence check
	// We'll use the same approach as the context gathering code
	_, err := os.Stat(filename)
	return err == nil
}

// hasFileOrDir checks if a file or directory exists
func hasFileOrDir(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

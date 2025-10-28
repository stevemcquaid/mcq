package cmd

import (
	"github.com/spf13/cobra"
	"github.com/stevemcquaid/mcq/pkg/ai"
)

// extractContextConfig extracts context configuration from command flags
func extractContextConfig(cmd *cobra.Command) ai.ContextConfig {
	autoDetect, _ := cmd.Flags().GetBool("auto-context")
	includeReadme, _ := cmd.Flags().GetBool("include-readme")
	includeGoMod, _ := cmd.Flags().GetBool("include-go-mod")
	includeCommits, _ := cmd.Flags().GetBool("include-commits")
	includeStructure, _ := cmd.Flags().GetBool("include-structure")
	includeConfigs, _ := cmd.Flags().GetBool("include-configs")
	maxCommits, _ := cmd.Flags().GetInt("max-commits")
	noContext, _ := cmd.Flags().GetBool("no-context")

	// Determine context configuration
	if noContext {
		return ai.ContextConfig{}
	}

	// Default to auto-detect if no specific flags are set
	if !autoDetect && !includeReadme && !includeGoMod && !includeCommits && !includeStructure && !includeConfigs {
		// Enable auto-detect by default
		autoDetect = true
	}

	if autoDetect || includeReadme || includeGoMod || includeCommits || includeStructure || includeConfigs {
		return ai.ContextConfig{
			AutoDetect:       autoDetect,
			IncludeReadme:    includeReadme,
			IncludeGoMod:     includeGoMod,
			IncludeCommits:   includeCommits,
			IncludeStructure: includeStructure,
			IncludeConfigs:   includeConfigs,
			MaxCommits:       maxCommits,
			MaxFileSize:      50 * 1024, // 50KB default
		}
	}

	// Ask user interactively
	return ai.PromptForContext()
}

// addAIFlags adds common AI flags to a command
func addAIFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("model", "m", "", "AI model to use: 'claude', 'gpt-4o', 'gpt-5', 'gpt-5-mini', or 'gpt-5-nano' (auto-detected if not specified)")
	cmd.Flags().IntP("verbosity", "v", 0, "Set verbosity level: 0=off, 1=basic, 2=detailed, 3=verbose")
	cmd.Flags().Bool("auto-context", false, "Automatically detect and include relevant repository context")
	cmd.Flags().Bool("include-readme", false, "Include README content in context")
	cmd.Flags().Bool("include-go-mod", false, "Include go.mod information in context")
	cmd.Flags().Bool("include-commits", false, "Include recent commit messages in context")
	cmd.Flags().Bool("include-structure", false, "Include directory structure in context")
	cmd.Flags().Bool("include-configs", false, "Include configuration files in context")
	cmd.Flags().Int("max-commits", 10, "Maximum number of recent commits to include")
	cmd.Flags().Bool("no-context", false, "Skip context gathering entirely")
}

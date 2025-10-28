package ai

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/stevemcquaid/mcq/pkg/logger"
)

// GatherContextIfNeeded gathers repository context if any context options are enabled
func GatherContextIfNeeded(config ContextConfig) *RepoContext {
	if !shouldGatherContext(config) {
		return nil
	}

	logger.LogBasic("Gathering repository context")
	repoContext, err := gatherRepoContext(config)
	if err != nil {
		logger.LogError("context gathering", err)

		// Show more detailed error information
		fmt.Printf("⚠️  Context Gathering Failed\n")
		fmt.Printf("💡 Continuing without context (results may be less accurate)\n")
		fmt.Printf("\n")
		fmt.Printf("⚠️  DETAILED ERROR:\n")
		fmt.Printf("   %v\n", err)
		fmt.Printf("\n")
		fmt.Printf("💡 This usually means:\n")
		fmt.Printf("   - Not running from a Git repository root\n")
		fmt.Printf("   - Missing go.mod file\n")
		fmt.Printf("   - No README files found\n")
		fmt.Printf("   - File permission issues\n")
		fmt.Printf("\n")

		return nil
	}

	logger.LogBasic("Repository context gathered successfully")

	// Debug: Log context size information
	if repoContext != nil {
		readmeSize := len(repoContext.Readme)
		structureSize := len(repoContext.DirectoryStructure)
		totalSize := readmeSize + structureSize

		logger.LogBasic("Context size info",
			"readme_chars", readmeSize,
			"structure_chars", structureSize,
			"total_chars", totalSize,
			"commits", len(repoContext.RecentCommits),
			"deps", len(repoContext.Dependencies),
		)

		// Warn if context is very large
		if totalSize > 100000 {
			fmt.Printf("⚠️  Warning: Large context (%d chars) may exceed token limits\n", totalSize)
		}
	}

	return repoContext
}

// shouldGatherContext determines if any context should be gathered
func shouldGatherContext(config ContextConfig) bool {
	return config.AutoDetect || config.IncludeReadme || config.IncludeGoMod ||
		config.IncludeCommits || config.IncludeStructure || config.IncludeConfigs
}

// gatherRepoContext gathers repository context based on configuration
func gatherRepoContext(config ContextConfig) (*RepoContext, error) {
	ctx := &RepoContext{
		ConfigFiles: make(map[string]string),
	}

	// Apply auto-detect settings if enabled
	config = applyAutoDetectSettings(config)

	// Gather all context components
	errors := gatherContextComponents(ctx, config)

	// Determine project type
	ctx.ProjectType = determineProjectType(ctx)

	// Return error with details about what failed
	if len(errors) > 0 {
		// If we got some context, don't fail completely but log what failed
		hasSomeContext := ctx.ProjectName != "" || ctx.Readme != "" || len(ctx.RecentCommits) > 0 || ctx.DirectoryStructure != "" || len(ctx.ConfigFiles) > 0

		if !hasSomeContext {
			// No context at all - return error with details
			errMsg := fmt.Sprintf("failed to gather any context. Errors: %v", errors)
			logger.LogError("Context gathering failed completely", fmt.Errorf("%s", errMsg))
			return nil, fmt.Errorf("%s", errMsg)
		} else {
			// Some context but errors occurred - log and return partial context
			logger.LogBasic("Partial context gathered with errors", "error_count", len(errors))
			for _, err := range errors {
				logger.LogBasic("Context error", "error", err)
			}
		}
	}

	return ctx, nil
}

// applyAutoDetectSettings applies auto-detect settings to the config
func applyAutoDetectSettings(config ContextConfig) ContextConfig {
	if config.AutoDetect {
		config.IncludeReadme = true
		config.IncludeGoMod = true
		config.IncludeCommits = true
		config.IncludeStructure = true
		config.IncludeConfigs = true
		config.MaxCommits = 10
		config.MaxFileSize = 50 * 1024 // 50KB
	}
	return config
}

// gatherContextComponents gathers all enabled context components
func gatherContextComponents(ctx *RepoContext, config ContextConfig) []error {
	var errors []error

	if config.IncludeGoMod {
		if err := gatherGoModuleInfo(ctx); err != nil {
			logger.LogBasic("Failed to gather Go module info", "error", err)
			errors = append(errors, fmt.Errorf("go module info: %w", err))
		}
	}

	if config.IncludeReadme {
		if err := gatherReadme(ctx); err != nil {
			logger.LogBasic("Failed to gather README", "error", err)
			errors = append(errors, fmt.Errorf("readme: %w", err))
		}
	}

	if config.IncludeCommits {
		if err := gatherRecentCommits(ctx, config.MaxCommits); err != nil {
			logger.LogBasic("Failed to gather recent commits", "error", err)
			errors = append(errors, fmt.Errorf("recent commits: %w", err))
		}
	}

	if config.IncludeStructure {
		if err := gatherDirectoryStructure(ctx); err != nil {
			logger.LogBasic("Failed to gather directory structure", "error", err)
			errors = append(errors, fmt.Errorf("directory structure: %w", err))
		}
	}

	if config.IncludeConfigs {
		if err := gatherConfigFiles(ctx, config.MaxFileSize); err != nil {
			logger.LogBasic("Failed to gather config files", "error", err)
			errors = append(errors, fmt.Errorf("config files: %w", err))
		}
	}

	return errors
}

// gatherComponent is a helper to gather a context component with error logging (deprecated)
func gatherComponent(_ *RepoContext, shouldGather bool, componentName string, gatherFunc func() error) {
	if !shouldGather {
		return
	}

	if err := gatherFunc(); err != nil {
		logger.LogBasic("Failed to gather "+componentName, "error", err)
	}
}

// gatherGoModuleInfo extracts information from go.mod
func gatherGoModuleInfo(ctx *RepoContext) error {
	content, err := os.ReadFile("go.mod")
	if err != nil {
		return fmt.Errorf("go.mod not found: %w", err)
	}

	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Extract module path
		if strings.HasPrefix(line, "module ") {
			ctx.ModulePath = strings.TrimSpace(strings.TrimPrefix(line, "module "))
			ctx.ProjectName = filepath.Base(ctx.ModulePath)
		}

		// Extract Go version
		if strings.HasPrefix(line, "go ") {
			ctx.GoVersion = strings.TrimSpace(strings.TrimPrefix(line, "go "))
		}

		// Extract dependencies
		if isDependencyLine(line) {
			parts := strings.Fields(line)
			if len(parts) > 0 {
				ctx.Dependencies = append(ctx.Dependencies, parts[0])
			}
		}
	}

	return nil
}

// isDependencyLine checks if a line contains a dependency
func isDependencyLine(line string) bool {
	return strings.HasPrefix(line, "require ") ||
		(strings.HasPrefix(line, "\t") && !strings.HasPrefix(line, "\t//") &&
			!strings.HasPrefix(line, "go ") && !strings.HasPrefix(line, "module "))
}

// gatherReadme extracts README content from root and docs directory
func gatherReadme(ctx *RepoContext) error {
	// First, try root directory
	readmeFiles := []string{
		"README.md", "README.rst", "README.txt", "README",
	}

	for _, filename := range readmeFiles {
		content, err := os.ReadFile(filename)
		if err == nil {
			ctx.Readme = string(content)
			// Continue to check docs directory for additional content
		}
	}

	// Also check docs directory if it exists
	if hasFileOrDir("docs") {
		// First, try to read docs/README.md
		docsReadmeFiles := []string{
			"docs/README.md", "docs/README.rst", "docs/README.txt", "docs/README",
		}

		for _, filename := range docsReadmeFiles {
			content, err := os.ReadFile(filename)
			if err == nil {
				// Append docs content to existing README
				if ctx.Readme == "" {
					ctx.Readme = string(content)
				} else {
					ctx.Readme += "\n\n## Documentation\n\n" + string(content)
				}
				break // Only read the first found file
			}
		}

		// Walk the docs directory to find all .md files
		err := filepath.Walk("docs", func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil // Skip errors
			}

			// Only process .md files (excluding README which we already handled)
			if !info.IsDir() && strings.HasSuffix(path, ".md") && !strings.Contains(path, "README") {
				content, err := os.ReadFile(path)
				if err == nil {
					// Extract just the filename without path
					fileName := filepath.Base(path)
					sectionName := strings.TrimSuffix(fileName, ".md")
					// Replace hyphens and underscores with spaces for better section titles
					sectionName = strings.ReplaceAll(sectionName, "-", " ")
					sectionName = strings.ReplaceAll(sectionName, "_", " ")
					// Capitalize each word
					words := strings.Fields(sectionName)
					for i, word := range words {
						if len(word) > 0 {
							words[i] = strings.ToUpper(word[:1]) + strings.ToLower(word[1:])
						}
					}
					sectionName = strings.Join(words, " ")

					ctx.Readme += fmt.Sprintf("\n\n### %s\n\n", sectionName) + string(content)
				}
			}
			return nil
		})

		if err != nil {
			logger.LogDebug("Error walking docs directory", "error", err)
		}
	}

	if ctx.Readme == "" {
		return fmt.Errorf("no README file found")
	}

	return nil
}

// gatherRecentCommits gets recent commit messages
func gatherRecentCommits(ctx *RepoContext, maxCommits int) error {
	cmd := exec.Command("git", "log", "--oneline", fmt.Sprintf("-n%d", maxCommits))
	output, err := cmd.Output()
	if err != nil {
		return err
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	for _, line := range lines {
		if line != "" {
			ctx.RecentCommits = append(ctx.RecentCommits, line)
		}
	}

	return nil
}

// gatherDirectoryStructure gets a high-level directory structure
func gatherDirectoryStructure(ctx *RepoContext) error {
	var structure strings.Builder

	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Skip errors
		}

		if shouldSkipPath(path) {
			return getSkipAction()
		}

		addToStructure(&structure, path, info)
		return nil
	})

	ctx.DirectoryStructure = structure.String()
	return err
}

// shouldSkipPath determines if a path should be skipped
func shouldSkipPath(path string) bool {
	// Skip hidden directories
	if strings.HasPrefix(path, ".") && path != "." {
		return true
	}

	// Skip common directories that don't add value
	skipDirs := []string{"vendor", "node_modules", ".git", "build", "dist", "target", "bin", "obj"}
	for _, skipDir := range skipDirs {
		if strings.Contains(path, skipDir) {
			return true
		}
	}

	return false
}

// getSkipAction returns the appropriate skip action for a directory
func getSkipAction() error {
	return filepath.SkipDir
}

// addToStructure adds a path to the directory structure
func addToStructure(structure *strings.Builder, path string, info os.FileInfo) {
	depth := strings.Count(path, string(filepath.Separator))
	indent := strings.Repeat("  ", depth)

	if info.IsDir() {
		if depth < 3 { // Limit depth to avoid too much detail
			fmt.Fprintf(structure, "%s%s/\n", indent, info.Name())
		}
	} else if isImportantFile(path) {
		fmt.Fprintf(structure, "%s%s\n", indent, info.Name())
	}
}

// isImportantFile determines if a file is important for context
func isImportantFile(path string) bool {
	importantExts := []string{".go", ".md", ".yaml", ".yml", ".json", ".toml", ".env", ".dockerfile", "Dockerfile", "Makefile"}
	importantNames := []string{"go.mod", "go.sum", "README", "LICENSE", "CHANGELOG", "Dockerfile", "Makefile", ".gitignore"}

	ext := filepath.Ext(path)
	for _, importantExt := range importantExts {
		if ext == importantExt {
			return true
		}
	}

	filename := filepath.Base(path)
	for _, importantName := range importantNames {
		if filename == importantName {
			return true
		}
	}

	return false
}

// gatherConfigFiles collects relevant configuration files
func gatherConfigFiles(ctx *RepoContext, maxSize int64) error {
	configFiles := []string{
		"go.mod", "go.sum", "Makefile", "Dockerfile", ".dockerignore",
		"docker-compose.yml", "docker-compose.yaml", ".env", ".env.example",
		"config.yaml", "config.yml", "config.json", ".gitignore",
	}

	for _, filename := range configFiles {
		if info, err := os.Stat(filename); err == nil && info.Size() <= maxSize {
			if content, err := os.ReadFile(filename); err == nil {
				ctx.ConfigFiles[filename] = string(content)
			}
		}
	}

	return nil
}

// hasFileOrDir checks if a file or directory exists
func hasFileOrDir(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// determineProjectType analyzes the repository to determine project type
func determineProjectType(ctx *RepoContext) string {
	// Check README for project type indicators
	if strings.Contains(ctx.Readme, "CLI") || strings.Contains(ctx.Readme, "command") {
		return "CLI Tool"
	}
	if strings.Contains(ctx.Readme, "API") || strings.Contains(ctx.Readme, "server") {
		return "Web API"
	}
	if strings.Contains(ctx.Readme, "library") || strings.Contains(ctx.Readme, "package") {
		return "Library"
	}

	// Check dependencies for clues
	for _, dep := range ctx.Dependencies {
		if strings.Contains(dep, "gin") || strings.Contains(dep, "echo") || strings.Contains(dep, "fiber") {
			return "Web API"
		}
		if strings.Contains(dep, "cobra") || strings.Contains(dep, "cli") {
			return "CLI Tool"
		}
	}

	// Check directory structure
	if strings.Contains(ctx.DirectoryStructure, "cmd/") {
		return "CLI Tool"
	}
	if strings.Contains(ctx.DirectoryStructure, "api/") || strings.Contains(ctx.DirectoryStructure, "server/") {
		return "Web API"
	}

	return "Go Application"
}

// formatContextForPrompt formats the repository context for inclusion in AI prompts
func formatContextForPrompt(ctx *RepoContext) string {
	if ctx == nil {
		return ""
	}

	var contextBuilder strings.Builder
	contextBuilder.WriteString("\n## Repository Context\n\n")

	// Project information
	contextBuilder.WriteString("### Project Information\n")
	contextBuilder.WriteString(fmt.Sprintf("- **Project Name**: %s\n", ctx.ProjectName))
	contextBuilder.WriteString(fmt.Sprintf("- **Module Path**: %s\n", ctx.ModulePath))
	contextBuilder.WriteString(fmt.Sprintf("- **Go Version**: %s\n", ctx.GoVersion))
	contextBuilder.WriteString(fmt.Sprintf("- **Project Type**: %s\n\n", ctx.ProjectType))

	// Dependencies
	if len(ctx.Dependencies) > 0 {
		contextBuilder.WriteString("### Key Dependencies\n")
		for _, dep := range ctx.Dependencies[:minInt(10, len(ctx.Dependencies))] {
			contextBuilder.WriteString(fmt.Sprintf("- %s\n", dep))
		}
		contextBuilder.WriteString("\n")
	}

	// README excerpt
	if ctx.Readme != "" {
		contextBuilder.WriteString("### Project Overview\n")
		readmeExcerpt := ctx.Readme
		if len(readmeExcerpt) > 1000 {
			readmeExcerpt = readmeExcerpt[:1000] + "..."
		}
		contextBuilder.WriteString(readmeExcerpt)
		contextBuilder.WriteString("\n\n")
	}

	// Recent commits
	if len(ctx.RecentCommits) > 0 {
		contextBuilder.WriteString("### Recent Development Activity\n")
		for _, commit := range ctx.RecentCommits[:minInt(5, len(ctx.RecentCommits))] {
			contextBuilder.WriteString(fmt.Sprintf("- %s\n", commit))
		}
		contextBuilder.WriteString("\n")
	}

	// Directory structure
	if ctx.DirectoryStructure != "" {
		contextBuilder.WriteString("### Project Structure\n")
		contextBuilder.WriteString("```\n")
		contextBuilder.WriteString(ctx.DirectoryStructure)
		contextBuilder.WriteString("\n```\n\n")
	}

	// Configuration files
	if len(ctx.ConfigFiles) > 0 {
		contextBuilder.WriteString("### Configuration Files\n")
		for filename, content := range ctx.ConfigFiles {
			contextBuilder.WriteString(fmt.Sprintf("**%s**:\n", filename))
			if len(content) > 500 {
				content = content[:500] + "..."
			}
			contextBuilder.WriteString("```\n")
			contextBuilder.WriteString(content)
			contextBuilder.WriteString("\n```\n\n")
		}
	}

	return contextBuilder.String()
}

// minInt returns the minimum of two integers
func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// PromptForContext interactively asks the user if they want to include context
func PromptForContext() ContextConfig {
	fmt.Println("\n🔍 Would you like to include repository context to improve the user story?")
	fmt.Println()
	fmt.Println("This will include:")
	fmt.Println("  📄 README content and project description")
	fmt.Println("  🔧 Go module information and dependencies")
	fmt.Println("  📝 Recent commit messages (last 10 commits)")
	fmt.Println("  📁 Directory structure overview")
	fmt.Println("  ⚙️ Configuration files (go.mod, Makefile, etc.)")
	fmt.Println()
	fmt.Println("This helps generate more accurate and contextually relevant user stories.")
	fmt.Println()

	// Ask if they want context
	fmt.Print("Include repository context? (Y/n): ")
	var response string
	_, err := fmt.Scanln(&response)
	if err != nil {
		fmt.Println("\n⚠️  Error reading input, skipping context gathering.")
		fmt.Println("   This is normal in non-interactive environments.")
		return ContextConfig{}
	}

	if strings.ToLower(response) == "n" || strings.ToLower(response) == "no" {
		fmt.Println("Skipping context gathering.")
		return ContextConfig{}
	}

	fmt.Println("✅ Including repository context...")
	return ContextConfig{
		AutoDetect:       true,
		IncludeReadme:    true,
		IncludeGoMod:     true,
		IncludeCommits:   true,
		IncludeStructure: true,
		IncludeConfigs:   true,
		MaxCommits:       10,
		MaxFileSize:      50 * 1024, // 50KB default
	}
}

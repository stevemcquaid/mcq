# Todo
* create markdown summary of the jira and ai features in order to implement them on another code base. any important architectural or parsing details include those too
* add `mcq jira update ticket` command
* add mcq git mr add ac
  * using the branch name to lookup the jira ticket and add the AC into the MR description
* add mcq git mr add details/analysis
  * get the diff of the current branch versus master, and summarize the changes made and a high level overview in order to help a code reviewer understand the complexities and improves of this code.
* improve the `mcq jira add` command
  * checkbox formatting for AC
  * user selectable line items
  * add `--dry-run` flag to replace `mcq ai jira` command with `mcq jira add --dry-run`
  * add `--current-sprint` flag to `mcq jira add`
* fix the streaming ai for gpt-5




## 1. **`mcq ai review` - AI-Powered Code Review**
Automatically review code changes (git diff, staged files, or PR) and provide feedback on:
- Code quality and best practices
- Potential bugs or edge cases
- Performance concerns
- Security vulnerabilities
- Suggestions for improvement
```bash
mcq ai review               # Review staged changes
mcq ai review --commit HEAD # Review last commit
mcq ai review --pr 123      # Review GitHub PR
```

## 2. **`mcq ai commit` - Generate Commit Messages**
Analyze staged changes and generate meaningful, conventional commit messages.
```bash
mcq ai commit                    # Generate and copy to clipboard
mcq ai commit --auto             # Generate and commit automatically
mcq ai commit --conventional     # Use conventional commits format
```

## 3. **`mcq jira comment` - AI-Enhanced Issue Comments**
Add intelligent comments to Jira issues based on:
- Recent code changes
- Development progress
- Blockers or questions discovered during implementation
```bash
mcq jira comment PROJ-123 "Implemented the API endpoint, added tests"
mcq jira comment PROJ-123 --from-commits  # Auto-generate from recent commits
mcq jira comment PROJ-123 --status-update # Generate progress update
```

## 4. **`mcq ai explain` - Code Explanation Generator**
Explain complex code snippets, functions, or files in plain English.
```bash
mcq ai explain path/to/file.go
mcq ai explain --function "calculateMetrics"
mcq ai explain --clipboard      # Explain code from clipboard
cat complex.go | mcq ai explain # Pipe code directly
```

## 5. **`mcq jira estimate` - AI Story Point Estimation**
Analyze user stories and suggest story point estimates based on:
- Story complexity
- Repository context
- Historical similar issues (if available)
- Technical dependencies
```bash
mcq jira estimate PROJ-123
mcq jira estimate --update PROJ-123  # Update the issue directly
```

## 6. **`mcq ai test` - Generate Test Cases**
Generate unit tests, integration tests, or test scenarios for existing code.
```bash
mcq ai test path/to/function.go
mcq ai test --type unit path/to/file.go
mcq ai test --type integration --framework testify
mcq ai test --coverage         # Suggest tests for uncovered code
```

## 7. **`mcq jira subtask` - Auto-Generate Subtasks**
Break down a Jira story into technical subtasks automatically.
```bash
mcq jira subtask PROJ-123                    # Generate and preview subtasks
mcq jira subtask PROJ-123 --create           # Create subtasks in Jira
mcq jira subtask PROJ-123 --assign-to team   # Auto-assign based on task type
```

## 8. **`mcq ai docs` - Documentation Generator**
Generate or update documentation from code, comments, and context.
```bash
mcq ai docs README.md                # Update README with current project state
mcq ai docs --api                    # Generate API documentation
mcq ai docs --godoc path/to/package  # Generate GoDoc comments
mcq ai docs --changelog              # Generate CHANGELOG from commits
```

## 9. **`mcq jira search` - Smart Issue Search**
Natural language search for Jira issues with AI understanding.
```bash
mcq jira search "issues about authentication from last sprint"
mcq jira search "my open bugs" --format table
mcq jira search "unassigned stories in current sprint" --assign-to me
```

## 10. **`mcq ai refactor` - Refactoring Suggestions**
Analyze code and suggest refactoring opportunities with implementation details.
```bash
mcq ai refactor path/to/file.go
mcq ai refactor --type "extract-function"
mcq ai refactor --detect-patterns    # Find repeated patterns
mcq ai refactor --apply              # Apply suggested refactoring (with confirmation)
```

### 11. **`mcq ai pr` - PR Description Generator**
Generate comprehensive PR descriptions from commits and changes.
```bash
mcq ai pr                        # Generate from current branch
mcq ai pr --template detailed    # Use detailed template
mcq ai pr --create               # Create PR with generated description
```

### 12. **`mcq jira transition` - Smart Status Updates**
Intelligently transition issues based on context (commits, time, etc.).
```bash
mcq jira transition PROJ-123 --auto    # Auto-detect next status
mcq jira transition PROJ-123 "in review" --add-comment
```


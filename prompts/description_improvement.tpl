{{/* 
Description Improvement Template
Available variables:
- .OriginalDescription: The current Jira description to improve
- .RepositoryContext: Repository information (if available)
- .ProjectName: Project name from go.mod
- .ModulePath: Module path from go.mod
- .GoVersion: Go version from go.mod
- .ProjectType: Detected project type
- .Readme: README content
- .RecentCommits: Recent commit messages
- .Dependencies: Go dependencies
- .DirectoryStructure: Directory structure
- .ConfigFiles: Configuration files content
- .Now: Current timestamp
*/}}
Improve the following Jira issue description. Make it:
1. More comprehensive and detailed
2. Better structured and readable
3. Include proper user story format if missing
4. Add acceptance criteria if not present
5. Add technical considerations
6. Ensure it follows best practices for user stories

Preserve the existing intent and structure, but enhance clarity, completeness, and professionalism.

Original Description:
{{.OriginalDescription}}
{{if .RepositoryContext}}
{{formatContext .RepositoryContext}}
{{end}}

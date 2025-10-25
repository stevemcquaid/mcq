{{/* 
User Story Generation Template
Available variables:
- .FeatureRequest: The user's feature request
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
Please convert the following vague feature request into a detailed user story. The user story should follow the format: "As a [user type], I want [goal] so that [benefit]". Additionally, include any relevant acceptance criteria and technical considerations. Provide ONLY the user story. 


Please provide a comprehensive user story:
1. With the main user story in the specified format
2. With acceptance criteria
3. With any relevant technical notes or considerations
4. Keep the total output under 1000 words

Do NOT add any additional questions or commentary. 
The response must ONLY be the user story. 
NOTHING ELSE.

Feature Request: {{.FeatureRequest}}
{{if .RepositoryContext}}
{{formatContext .RepositoryContext}}
{{end}}

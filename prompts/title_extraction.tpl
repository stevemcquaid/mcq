{{/* 
Title Extraction Template
Available variables:
- .FeatureRequest: The original feature request
- .UserStory: The generated user story
- .Now: Current timestamp
*/}}
Create a NEW concise, clear title (maximum 100 characters) for a Jira issue from the following user story and old title. The new title should be action-oriented and summarize the main goal or feature.
Provide ONLY the new jira title
Do NOT provide any other output.

Original Feature Request: {{.FeatureRequest}}

User Story: 
{{.UserStory}}

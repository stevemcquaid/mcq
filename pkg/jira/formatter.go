package jira

import (
	"fmt"
	"net/url"
	"strings"
)

// TextFormatter handles all Jira text formatting and conversion
type TextFormatter struct {
	userCache map[string]string
}

// NewTextFormatter creates a new text formatter
func NewTextFormatter() *TextFormatter {
	return &TextFormatter{
		userCache: make(map[string]string),
	}
}

// FormatText applies all Jira text formatting (HTML cleaning, links, code blocks, etc.)
func (tf *TextFormatter) FormatText(text string) string {
	// Convert Jira links to markdown or plain text format
	text = tf.convertJiraLinks(text)

	// Convert Jira code blocks and formatting
	text = tf.convertJiraCodeBlocks(text)

	// Clean HTML tags
	text = tf.cleanHTML(text)

	return text
}

// convertJiraLinks converts Jira link format to markdown or plain text
func (tf *TextFormatter) convertJiraLinks(text string) string {
	// Convert Jira user links [~accountid:...] to @username
	text = tf.convertJiraUserLinks(text)

	// Convert Jira smart links [text|url|smart-link] to markdown [text](url) or just url
	text = tf.convertJiraSmartLinks(text)

	return text
}

// convertJiraUserLinks converts Jira user account IDs to @username format
func (tf *TextFormatter) convertJiraUserLinks(text string) string {
	start := 0
	for {
		// Look for the start of a user link
		userStart := strings.Index(text[start:], "[~accountid:")
		if userStart == -1 {
			break
		}
		userStart += start

		// Find the end of the user link
		userEnd := strings.Index(text[userStart:], "]")
		if userEnd == -1 {
			break
		}
		userEnd += userStart + 1

		// Extract the account ID
		accountID := text[userStart+12 : userEnd-1] // Skip "[~accountid:" and "]"

		// Resolve the account ID to actual username
		username := tf.resolveAccountID(accountID)
		if username == "" {
			// Fallback to a simplified format if resolution fails
			username = "@user-" + accountID[len(accountID)-8:]
		}

		// Replace the user link with @username
		text = text[:userStart] + username + text[userEnd:]

		// Update start position to continue searching
		start = userStart + len(username)
	}

	return text
}

// convertJiraSmartLinks converts Jira smart links to markdown or plain text format
func (tf *TextFormatter) convertJiraSmartLinks(text string) string {
	start := 0
	for {
		// Look for the start of a smart link
		linkStart := strings.Index(text[start:], "[")
		if linkStart == -1 {
			break
		}
		linkStart += start

		// Find the end of the link
		linkEnd := strings.Index(text[linkStart:], "]")
		if linkEnd == -1 {
			break
		}
		linkEnd += linkStart + 1

		// Extract the link content
		linkContent := text[linkStart+1 : linkEnd-1] // Skip "[" and "]"

		// Check if this is a smart link format [text|url|smart-link] or [text|url]
		parts := strings.Split(linkContent, "|")
		if (len(parts) == 3 && parts[2] == "smart-link") || (len(parts) == 2) {
			linkText := parts[0]
			linkURL := parts[1]

			// Parse and clean the URL to make it valid
			var decodedURL string
			parsedURL, err := url.Parse(linkURL)
			if err != nil {
				// If parsing fails, use the original URL
				decodedURL = linkURL
			} else {
				// The parsed URL should be properly formatted
				decodedURL = parsedURL.String()
			}

			var convertedLink string
			if linkText == linkURL {
				// If text and URL are the same, just print the decoded URL once
				convertedLink = decodedURL
			} else {
				// If text and URL are different, use markdown format with decoded URL
				convertedLink = fmt.Sprintf("[%s](%s)", linkText, decodedURL)
			}

			text = text[:linkStart] + convertedLink + text[linkEnd:]
			start = linkStart + len(convertedLink)
		} else {
			// Not a smart link, continue searching
			start = linkStart + 1
		}
	}

	return text
}

// convertJiraCodeBlocks converts Jira code formatting to markdown
func (tf *TextFormatter) convertJiraCodeBlocks(text string) string {
	// Convert {noformat} blocks to triple backticks
	text = tf.convertJiraNoFormatBlocks(text)

	// Convert {code} blocks to triple backticks
	text = tf.convertJiraCodeBlocksWithLang(text)

	// Convert inline code formatting
	text = tf.convertJiraInlineCode(text)

	// Convert headings
	text = tf.convertJiraHeadings(text)

	return text
}

// convertJiraNoFormatBlocks converts {noformat}...{noformat} to ```...```
func (tf *TextFormatter) convertJiraNoFormatBlocks(text string) string {
	start := 0
	for {
		// Look for {noformat}
		startTag := strings.Index(text[start:], "{noformat}")
		if startTag == -1 {
			break
		}
		startTag += start

		// Look for closing {noformat}
		endTag := strings.Index(text[startTag+10:], "{noformat}")
		if endTag == -1 {
			break
		}
		endTag += startTag + 10

		// Extract the content between tags
		content := text[startTag+10 : endTag]

		// Replace with markdown code block
		markdownBlock := fmt.Sprintf("```\n%s\n```", content)
		text = text[:startTag] + markdownBlock + text[endTag+10:]

		// Update start position
		start = startTag + len(markdownBlock)
	}

	return text
}

// convertJiraCodeBlocksWithLang converts {code:lang}...{code} to ```lang...```
func (tf *TextFormatter) convertJiraCodeBlocksWithLang(text string) string {
	start := 0
	for {
		// Look for {code:lang} or {code}
		startTag := strings.Index(text[start:], "{code")
		if startTag == -1 {
			break
		}
		startTag += start

		// Find the end of the opening tag
		tagEnd := strings.Index(text[startTag:], "}")
		if tagEnd == -1 {
			break
		}
		tagEnd += startTag + 1

		// Extract language if present
		lang := ""
		if text[startTag+5:tagEnd-1] != "" {
			lang = text[startTag+6 : tagEnd-1] // Skip ":"
		}

		// Look for closing {code}
		endTag := strings.Index(text[tagEnd:], "{code}")
		if endTag == -1 {
			break
		}
		endTag += tagEnd

		// Extract the content between tags
		content := text[tagEnd:endTag]

		// Replace with markdown code block
		markdownBlock := fmt.Sprintf("```%s\n%s\n```", lang, content)
		text = text[:startTag] + markdownBlock + text[endTag+6:]

		// Update start position
		start = startTag + len(markdownBlock)
	}

	return text
}

// convertJiraInlineCode converts {{code}} to `code`
func (tf *TextFormatter) convertJiraInlineCode(text string) string {
	start := 0
	for {
		// Look for {{code}}
		startTag := strings.Index(text[start:], "{{")
		if startTag == -1 {
			break
		}
		startTag += start

		// Look for closing }}
		endTag := strings.Index(text[startTag+2:], "}}")
		if endTag == -1 {
			break
		}
		endTag += startTag + 2

		// Extract the content between tags
		content := text[startTag+2 : endTag]

		// Replace with markdown inline code
		markdownCode := fmt.Sprintf("`%s`", content)
		text = text[:startTag] + markdownCode + text[endTag+2:]

		// Update start position
		start = startTag + len(markdownCode)
	}

	return text
}

// convertJiraHeadings converts Jira headings to markdown
func (tf *TextFormatter) convertJiraHeadings(text string) string {
	// Convert h1. to #
	text = strings.ReplaceAll(text, "h1. ", "# ")

	// Convert h2. to ##
	text = strings.ReplaceAll(text, "h2. ", "## ")

	// Convert h3. to ###
	text = strings.ReplaceAll(text, "h3. ", "### ")

	// Convert h4. to ####
	text = strings.ReplaceAll(text, "h4. ", "#### ")

	// Convert h5. to #####
	text = strings.ReplaceAll(text, "h5. ", "##### ")

	// Convert h6. to ######
	text = strings.ReplaceAll(text, "h6. ", "###### ")

	return text
}

// cleanHTML removes basic HTML tags from text
func (tf *TextFormatter) cleanHTML(text string) string {
	// Simple HTML tag removal - this is basic but covers most cases
	text = strings.ReplaceAll(text, "<br>", "\n")
	text = strings.ReplaceAll(text, "<br/>", "\n")
	text = strings.ReplaceAll(text, "<br />", "\n")
	text = strings.ReplaceAll(text, "<p>", "\n")
	text = strings.ReplaceAll(text, "</p>", "\n")
	text = strings.ReplaceAll(text, "<strong>", "")
	text = strings.ReplaceAll(text, "</strong>", "")
	text = strings.ReplaceAll(text, "<em>", "")
	text = strings.ReplaceAll(text, "</em>", "")
	text = strings.ReplaceAll(text, "<b>", "")
	text = strings.ReplaceAll(text, "</b>", "")
	text = strings.ReplaceAll(text, "<i>", "")
	text = strings.ReplaceAll(text, "</i>", "")

	// Remove any remaining HTML tags (basic regex-like approach)
	for {
		start := strings.Index(text, "<")
		if start == -1 {
			break
		}
		end := strings.Index(text[start:], ">")
		if end == -1 {
			break
		}
		text = text[:start] + text[start+end+1:]
	}

	// Clean up extra whitespace
	text = strings.TrimSpace(text)
	return text
}

// resolveAccountID attempts to resolve a Jira account ID to a username
func (tf *TextFormatter) resolveAccountID(accountID string) string {
	// Check cache first
	if username, exists := tf.userCache[accountID]; exists {
		return username
	}

	// For now, return a fallback since we don't have access to config in formatter
	// This could be improved by passing config to the formatter
	return "@user-" + accountID[len(accountID)-8:]
}

package ai

import (
	"strings"

	"github.com/stevemcquaid/mcq/pkg/errors"
	"github.com/stevemcquaid/mcq/pkg/logger"
)

// AIJira converts a vague feature request to a user story and copies it to clipboard
func AIJira(args []string, modelFlag string, verbosityLevel int, contextConfig ContextConfig) error {
	logger.SetupLogger(verbosityLevel)

	featureRequest := strings.Join(args, " ")
	logger.LogBasic("Starting AIJira", "feature_request", featureRequest)

	// Gather repository context
	repoContext := GatherContextIfNeeded(contextConfig)

	// Select and configure model
	selectedModel, err := SelectModel(modelFlag)
	if err != nil {
		userErr := errors.WrapError(err, "Failed to select AI model")
		userErr.Display()
		return userErr
	}

	// Generate user story
	userStory, err := GenerateUserStory(selectedModel, featureRequest, repoContext)
	if err != nil {
		userErr := errors.WrapError(err, "Failed to generate user story")
		userErr.Display()
		return userErr
	}

	// Display and copy result
	return displayAndCopyResult(userStory)
}

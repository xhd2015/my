package openclaw

import (
	"fmt"
	"os"
)

func runResolveSlackChannel(dataDir, containerName, channelName string) int {
	if err := ensurePodmanMachine(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	resolvedDataDir, err := requireRunningContainerDataDir(dataDir, containerName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	absPath, err := resolvePath(resolvedDataDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	view, checks, err := diagnoseSlackConfig(absPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}
	for _, check := range checks {
		if check.missing {
			fmt.Fprintf(os.Stderr, "error: Slack is not correctly configured (%s)\n", check.name)
			return 1
		}
	}
	if view.botToken == "" {
		fmt.Fprintf(os.Stderr, "error: slack bot token not found in %s\n", absPath)
		return 1
	}

	fmt.Printf("Looking up Slack channel %q via bot token from %s\n", channelName, absPath)
	channel, err := lookupSlackChannelByName(view.botToken, channelName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	fmt.Printf("Channel: #%s\n", channel.Name)
	fmt.Printf("ID: %s\n", channel.ID)
	if !channel.IsMember {
		fmt.Println("Note: the bot is not a member of this channel yet; invite it before sending messages.")
	}
	return 0
}
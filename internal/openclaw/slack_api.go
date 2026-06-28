package openclaw

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

const (
	slackAPIBase            = "https://slack.com/api/"
	defaultSlackTestChannel = "general"
)

var (
	slackConversationsListURL = slackAPIBase + "conversations.list"
	slackChatPostMessageURL   = slackAPIBase + "chat.postMessage"
	slackHTTPClient           = &http.Client{Timeout: 30 * time.Second}
)

type slackAPIResponse struct {
	OK        bool   `json:"ok"`
	Error     string `json:"error"`
	Needed    string `json:"needed"`
	Provided  string `json:"provided"`
}

type slackChannel struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	IsArchived bool   `json:"is_archived"`
	IsMember   bool   `json:"is_member"`
	IsGeneral  bool   `json:"is_general"`
	IsPrivate  bool   `json:"is_private"`
}

type slackConversationsListResponse struct {
	slackAPIResponse
	Channels []slackChannel `json:"channels"`
	Metadata struct {
		NextCursor string `json:"next_cursor"`
	} `json:"response_metadata"`
}

type slackChatPostMessageResponse struct {
	slackAPIResponse
	Channel string `json:"channel"`
	TS      string `json:"ts"`
}

func normalizeSlackChannelName(name string) string {
	name = strings.TrimSpace(name)
	name = strings.TrimPrefix(name, "#")
	return strings.ToLower(name)
}

func slackStubEnabled() bool {
	return os.Getenv("MY_OPENCLAW_SLACK_STUB") == "1"
}

func lookupSlackChannelByName(botToken, channelName string) (slackChannel, error) {
	want := normalizeSlackChannelName(channelName)
	if want == "" {
		return slackChannel{}, fmt.Errorf("channel name is required")
	}
	if strings.TrimSpace(botToken) == "" {
		return slackChannel{}, fmt.Errorf("slack bot token is required")
	}
	if slackStubEnabled() {
		return slackChannel{ID: "CSTUB" + want, Name: want}, nil
	}

	channelTypes := []string{"public_channel"}
	if want != normalizeSlackChannelName(defaultSlackTestChannel) {
		channelTypes = append(channelTypes, "private_channel")
	}

	for _, channelType := range channelTypes {
		channel, found, err := findSlackChannelInType(botToken, channelType, want)
		if err != nil {
			return slackChannel{}, err
		}
		if found {
			return channel, nil
		}
	}

	if want == normalizeSlackChannelName(defaultSlackTestChannel) {
		return slackChannel{}, fmt.Errorf("slack channel %q not found in public channels (workspace may use a different default channel name)", channelName)
	}
	return slackChannel{}, fmt.Errorf("slack channel %q not found (invite the bot to private channels, or add groups:read to the Slack app)", channelName)
}

func findSlackChannelInType(botToken, channelType, want string) (slackChannel, bool, error) {
	cursor := ""
	for page := 0; page < 20; page++ {
		channels, next, err := slackConversationsPage(botToken, channelType, cursor)
		if err != nil {
			return slackChannel{}, false, err
		}
		if want == normalizeSlackChannelName(defaultSlackTestChannel) {
			for _, channel := range channels {
				if !channel.IsArchived && channel.IsGeneral {
					return channel, true, nil
				}
			}
		}
		for _, channel := range channels {
			if channel.IsArchived {
				continue
			}
			if normalizeSlackChannelName(channel.Name) == want {
				return channel, true, nil
			}
		}
		if next == "" {
			break
		}
		cursor = next
	}
	return slackChannel{}, false, nil
}

func slackConversationsPage(botToken, channelType, cursor string) ([]slackChannel, string, error) {
	values := url.Values{}
	values.Set("types", channelType)
	values.Set("exclude_archived", "true")
	values.Set("limit", "200")
	if cursor != "" {
		values.Set("cursor", cursor)
	}

	req, err := http.NewRequest(http.MethodGet, slackConversationsListURL+"?"+values.Encode(), nil)
	if err != nil {
		return nil, "", err
	}
	req.Header.Set("Authorization", "Bearer "+botToken)

	resp, err := slackHTTPClient.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}

	var parsed slackConversationsListResponse
	if err := json.Unmarshal(body, &parsed); err != nil {
		return nil, "", fmt.Errorf("parse slack conversations.list: %w", err)
	}
	if !parsed.OK {
		return nil, "", formatSlackAPIError("conversations.list", parsed.slackAPIResponse)
	}

	return parsed.Channels, parsed.Metadata.NextCursor, nil
}

func formatSlackAPIError(method string, resp slackAPIResponse) error {
	if resp.Error == "" {
		resp.Error = "unknown_error"
	}
	msg := fmt.Sprintf("slack %s: %s", method, resp.Error)
	if resp.Error == "missing_scope" {
		msg += "\n\nThe bot token is missing a required Slack scope."
		if resp.Needed != "" {
			msg += fmt.Sprintf("\n  Needed:    %s", resp.Needed)
		}
		if resp.Provided != "" {
			msg += fmt.Sprintf("\n  Installed: %s", resp.Provided)
		}
		if resp.Needed == "groups:read" {
			msg += "\n\nTo resolve private channels by name, add groups:read in api.slack.com/apps → OAuth & Permissions → Bot Token Scopes, then reinstall the app to your workspace and update the bot token in openclaw.json."
		} else if resp.Needed == "channels:read" {
			msg += "\n\nAdd channels:read in api.slack.com/apps → OAuth & Permissions → Bot Token Scopes, then reinstall the app and update the bot token in openclaw.json."
		}
	}
	return fmt.Errorf("%s", msg)
}

func sendSlackChannelMessage(botToken, channelID, message string) error {
	if strings.TrimSpace(botToken) == "" {
		return fmt.Errorf("slack bot token is required")
	}
	if strings.TrimSpace(channelID) == "" {
		return fmt.Errorf("slack channel id is required")
	}
	if strings.TrimSpace(message) == "" {
		return fmt.Errorf("message is required")
	}
	if slackStubEnabled() {
		return nil
	}

	payload, err := json.Marshal(map[string]string{
		"channel": channelID,
		"text":    message,
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, slackChatPostMessageURL, strings.NewReader(string(payload)))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+botToken)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	resp, err := slackHTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var parsed slackChatPostMessageResponse
	if err := json.Unmarshal(body, &parsed); err != nil {
		return fmt.Errorf("parse slack chat.postMessage: %w", err)
	}
	if !parsed.OK {
		return formatSlackAPIError("chat.postMessage", parsed.slackAPIResponse)
	}
	return nil
}

func resolveSlackTestTargetFromConfig(view slackConfigView) (string, bool) {
	if view.testTarget != "" {
		return view.testTarget, true
	}
	for _, id := range view.channelIDs {
		if strings.HasPrefix(id, "C") || strings.HasPrefix(id, "G") || strings.HasPrefix(id, "D") {
			return id, true
		}
	}
	return "", false
}

func resolveSlackChannelTarget(botToken, channel string) (string, string, error) {
	channel = strings.TrimSpace(channel)
	if channel == "" {
		return "", "", fmt.Errorf("channel is required")
	}
	if isSlackChannelID(channel) {
		return channel, formatSlackTargetLabel(channel), nil
	}
	found, err := lookupSlackChannelByName(botToken, channel)
	if err != nil {
		return "", "", err
	}
	return found.ID, "#" + found.Name, nil
}

func isSlackChannelID(channel string) bool {
	return strings.HasPrefix(channel, "C") ||
		strings.HasPrefix(channel, "G") ||
		strings.HasPrefix(channel, "D")
}

func resolveSlackTestTarget(view slackConfigView, botToken, slackChannel string) (targetID, displayLabel string, err error) {
	if strings.TrimSpace(slackChannel) != "" {
		return resolveSlackChannelTarget(botToken, slackChannel)
	}
	if targetID, ok := resolveSlackTestTargetFromConfig(view); ok {
		return targetID, formatSlackTargetLabel(targetID), nil
	}

	channel, err := lookupSlackChannelByName(botToken, defaultSlackTestChannel)
	if err != nil {
		return "", "", err
	}
	return channel.ID, "#" + channel.Name, nil
}

func formatSlackTargetLabel(target string) string {
	target = strings.TrimSpace(target)
	if target == "" {
		return target
	}
	if strings.HasPrefix(target, "#") || strings.HasPrefix(target, "C") ||
		strings.HasPrefix(target, "G") || strings.HasPrefix(target, "D") ||
		strings.HasPrefix(target, "channel:") || strings.HasPrefix(target, "user:") {
		return target
	}
	return target
}
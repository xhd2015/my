package openclaw

import (
	"fmt"
	"os"
	"strings"
)

const (
	topHelpText = `my — personal activity tracker

Usage: my <command> [options]

Commands:
  openclaw    OpenClaw data-dir registry and Podman launcher
  opencode    OpenCode auth import and utilities

Run 'my <command> --help' for command-specific help.
`
	openclawHelpText = `Usage: my openclaw <subcommand> [options]

Subcommands:
  add data-dir <path> [--note "..."]   Register a .openclaw data directory
  list                                  List registered data directories
  run [action] [options]               Run or manage OpenClaw gateway on the host
  run-in-podman --data-dir <path>      Start OpenClaw gateway in Podman

Options:
  -h, --help   Show this help
`
	addHelpText = `Usage: my openclaw add data-dir <path> [--note "..."]

Register an OpenClaw .openclaw data directory for bookkeeping.

Options:
  --note <text>   Optional note for this data directory
  -h, --help      Show this help
`
	listHelpText = `Usage: my openclaw list

List registered OpenClaw data directories.

Options:
  -h, --help   Show this help
`
	runInPodmanHelpText = `Usage: my openclaw run-in-podman [action] [options]

Manage the OpenClaw gateway in Podman. Most actions require a running session;
data dir is auto-detected from the container unless --data-dir is given.

Lifecycle:
  (default)    Start gateway — requires --data-dir when container is not running
               [--rebuild] [--container-name NAME] [--port PORT]
  --stop       Stop and remove container [--container-name NAME]
  --restart    Stop (if running) and start gateway
               [--data-dir <path>] [--rebuild] [--container-name NAME] [--port PORT]

Inspect:
  --status     Show gateway URLs and info [--container-name NAME] [--port PORT]
  --logs       Show container logs [--container-name NAME]
  --show-tokens
               Show gateway tokens and auth URLs [--data-dir <path>]
  --dashboard  Open authenticated dashboard [--data-dir <path>] [--port PORT]

Settings:
  --export-settings <zip>
               Export running session to a zip (contains openclaw/ tree)

Slack:
  --install-slack-plugin
               Install @openclaw/slack into the session data dir (restart after)
  --resolve-slack-channel <name>
               Look up a channel ID by name (e.g. general, #random)
  --doctor-slack
               Check Slack config and probe connectivity
  --test-slack Send a test message (default: #general)
               [--slack-channel <name|id>]

Grok:
  --import-local-grok
               Import ~/.grok OAuth, sync models, reset main session
  --install-grok
               Install Grok CLI in the running container

Advanced:
  --exec <command...>
               Run a command in the running container (must be first flag)

Shared options:
  --data-dir <path>       Path to .openclaw directory
  --container-name <name> Container name (default: openclaw-gateway)
  --port <port>           Host port (default: 18789)
  --rebuild               Force container image rebuild (start/restart only)
  -h, --help              Show this help (action-specific with action --help)
`
	runInPodmanTestSlackHelpText = `Usage: my openclaw run-in-podman --test-slack [options]

Send a test message to Slack via the bot token in the running session.

Requires a running gateway container. Destination resolution order:
  1. --slack-channel <name|id>
  2. channels.slack.testTarget in openclaw.json
  3. First configured channel ID in channels.slack.channels
  4. #general (workspace default public channel)

Examples:
  my openclaw run-in-podman --test-slack
  my openclaw run-in-podman --test-slack --slack-channel random
  my openclaw run-in-podman --test-slack --slack-channel C01234567

Options:
  --slack-channel <name|id>   Channel name (#general) or Slack ID (C…)
  --container-name <name>       Container name (default: openclaw-gateway)
  --data-dir <path>             Override auto-detected data dir
  -h, --help                    Show this help
`
	runInPodmanExportSettingsHelpText = `Usage: my openclaw run-in-podman --export-settings <zip> [options]

Export the running session's OpenClaw data directory into a zip archive.
The archive contains a single top-level openclaw/ directory with the full
session tree (openclaw.json, .env, agents/, workspace/, plugins/, etc.).

Requires a running gateway container.

Examples:
  my openclaw run-in-podman --export-settings ./openclaw-settings.zip

Restore:
  unzip ./openclaw-settings.zip -d ~/backup
  my openclaw run-in-podman --data-dir ~/backup/openclaw

Options:
  --container-name <name>   Container name (default: openclaw-gateway)
  --data-dir <path>         Override auto-detected data dir
  -h, --help                Show this help
`
	runInPodmanStatusHelpText = `Usage: my openclaw run-in-podman --status [options]

Show gateway URLs, data dir, port, and handy follow-up commands.

Requires a running gateway container.

Options:
  --container-name <name>   Container name (default: openclaw-gateway)
  --port <port>             Host port (default: 18789)
  --data-dir <path>         Override auto-detected data dir
  -h, --help                Show this help
`
)

func wantsHelp(args []string) bool {
	for _, arg := range args {
		if arg == "-h" || arg == "--help" {
			return true
		}
	}
	return false
}

func hasFlag(args []string, name string) bool {
	for _, arg := range args {
		if arg == name {
			return true
		}
	}
	return false
}

func hasFlagValue(args []string, name string) bool {
	for i, arg := range args {
		if arg == name {
			return i+1 < len(args) && args[i+1] != "" && !strings.HasPrefix(args[i+1], "-")
		}
		if strings.HasPrefix(arg, name+"=") {
			return len(strings.TrimPrefix(arg, name+"=")) > 0
		}
	}
	return false
}

func runInPodmanHelpForArgs(args []string) string {
	if !wantsHelp(args) {
		return ""
	}
	switch {
	case hasFlag(args, "--test-slack"):
		return runInPodmanTestSlackHelpText
	case hasFlag(args, "--status"):
		return runInPodmanStatusHelpText
	case hasFlag(args, "--export-settings") || hasFlagValue(args, "--export-settings"):
		return runInPodmanExportSettingsHelpText
	default:
		return runInPodmanHelpText
	}
}

func printHelp(text string) {
	fmt.Print(strings.TrimPrefix(text, "\n"))
}

func TopHelp() string {
	return topHelpText
}

func Run(args []string) int {
	if len(args) == 0 || args[0] == "-h" || args[0] == "--help" {
		printHelp(openclawHelpText)
		return 0
	}

	switch args[0] {
	case "add":
		return runAdd(args[1:])
	case "list":
		return runList(args[1:])
	case "run":
		return runLocal(args[1:])
	case "run-in-podman":
		return runInPodman(args[1:])
	default:
		fmt.Fprintf(os.Stderr, "unknown openclaw subcommand: %s\n", args[0])
		return 1
	}
}
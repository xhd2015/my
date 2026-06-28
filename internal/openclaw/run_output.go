package openclaw

import "fmt"

func printGatewayURLs(port, token string) {
	base := fmt.Sprintf("http://127.0.0.1:%s", port)
	fmt.Println("URLs:")
	fmt.Printf("  Dashboard:  %s/\n", base)
	fmt.Printf("  Chat:       %s/chat?session=main\n", base)
	if token != "" {
		fmt.Printf("  Auth chat:  %s/chat?session=main#token=%s\n", base, token)
		fmt.Printf("  Auth dash:  %s/#token=%s\n", base, token)
	}
}

func printGatewayCommands(port, containerName string) {
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Printf("  my openclaw run-in-podman --status --container-name %s\n", containerName)
	fmt.Printf("  my openclaw run-in-podman --dashboard --port %s\n", port)
	fmt.Printf("  my openclaw run-in-podman --logs --container-name %s\n", containerName)
	fmt.Printf("  my openclaw run-in-podman --show-tokens --port %s\n", port)
	fmt.Printf("  my openclaw run-in-podman --stop --container-name %s\n", containerName)
	fmt.Printf("  my openclaw run-in-podman --restart --container-name %s\n", containerName)
	fmt.Printf("  podman logs -f %s\n", containerName)
	fmt.Printf("  podman exec %s openclaw status\n", containerName)
	fmt.Printf("  podman exec %s openclaw gateway status\n", containerName)
}

func printGatewayTokenNote(token string) {
	if token == "" {
		return
	}
	fmt.Println()
	fmt.Printf("Gateway token: configured (from openclaw.json or .env)\n")
	fmt.Println("Paste the token in Control UI settings, or open an Auth URL above.")
}

func printGatewayInfo(containerName, dataDir, port, token string, justStarted bool) {
	if justStarted {
		fmt.Printf("Gateway started: %s\n", containerName)
	} else {
		fmt.Printf("Container: %s (running)\n", containerName)
	}
	if dataDir != "" {
		fmt.Printf("Data dir: %s\n", dataDir)
	}
	fmt.Printf("Port: %s\n", port)
	fmt.Println()
	printGatewayURLs(port, token)
	printGatewayCommands(port, containerName)
	printGatewayTokenNote(token)
}

func printLaunchHelp(port, containerName, dataDir, token string) {
	printGatewayInfo(containerName, dataDir, port, token, true)
}

func printGatewayStatus(containerName, dataDir, port, token string) {
	printGatewayInfo(containerName, dataDir, port, token, false)
}

func printLocalGatewayCommands(port string) {
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Printf("  my openclaw run --status --port %s\n", port)
	fmt.Println("  my openclaw run --stop")
	fmt.Printf("  my openclaw run --restart --port %s\n", port)
	fmt.Println("  openclaw gateway status")
	fmt.Println("  openclaw gateway stop")
}

func printLocalGatewayStatus(dataDir, port, token string) {
	fmt.Println("Gateway running locally")
	if dataDir != "" {
		fmt.Printf("Data dir: %s\n", dataDir)
	}
	fmt.Printf("Port: %s\n", port)
	fmt.Println()
	printGatewayURLs(port, token)
	printLocalGatewayCommands(port)
	printGatewayTokenNote(token)
}

func printLocalLaunchHelp(dataDir, port, token string, bumped bool) {
	if bumped {
		fmt.Printf("Port %d in use; using %s instead\n", defaultGatewayPort, port)
	}
	fmt.Println("Gateway started locally")
	if dataDir != "" {
		fmt.Printf("Data dir: %s\n", dataDir)
	}
	fmt.Printf("Port: %s\n", port)
	fmt.Println()
	printGatewayURLs(port, token)
	printLocalGatewayCommands(port)
	printGatewayTokenNote(token)
}
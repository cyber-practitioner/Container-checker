package checks

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

// CheckContainerNetworks retrieves and prints network configurations for all containers.
func CheckContainerNetworks(cli *client.Client) error {
	// List all containers
	containers, err := cli.ContainerList(context.Background(), container.ListOptions{All: true})
	if err != nil {
		return fmt.Errorf("error listing containers: %v", err)
	}

	var insecureContainers []string
	for _, container := range containers {
		insecureNetworkSettings, err := checkContainerNetworks(cli, container.ID)
		if err != nil {
			return fmt.Errorf("error checking network settings for container %s: %v", container.ID[:12], err)
		}

		if len(insecureNetworkSettings) > 0 {
			insecureContainers = append(insecureContainers, container.ID[:12])
		}
	}

	if len(insecureContainers) > 0 {
		fmt.Println("Containers with insecure network configurations:")
		for _, containerID := range insecureContainers {
			fmt.Println("- Container ID:", containerID)
		}
		fmt.Println()
		fmt.Println("Recommendations:")
		fmt.Println("1. Do not enable the TCP Docker daemon socket. If you are running the Docker daemon with -H tcp://0.0.0.0:XXX or similar, you are exposing unencrypted and unauthenticated direct access to the Docker daemon. This can be a security risk if the host is internet-connected.")
		fmt.Println("2. Do not expose the /var/run/docker.sock to other containers. Mounting the socket read-only is not a solution, as it only makes it harder to exploit.")
		fmt.Println("3. Ensure that the ownership and permissions of the /var/run/docker.sock file are properly configured to restrict access to authorized users only.")
	} else {
		fmt.Println("No containers with insecure network configurations found.")
	}

	return nil
}

func checkContainerNetworks(cli *client.Client, containerID string) ([]string, error) {
	containerJSON, err := cli.ContainerInspect(context.Background(), containerID)
	if err != nil {
		return nil, err
	}

	// Get network settings
	networkSettings := containerJSON.NetworkSettings.Networks
	if len(networkSettings) == 0 {
		return nil, nil
	}

	var insecureNetworkSettings []string
	// Analyze network settings
	for networkName, settings := range networkSettings {
		// Check for insecure network configurations
		if settings.NetworkID == "host" {
			insecureNetworkSettings = append(insecureNetworkSettings, fmt.Sprintf("Container %s is using host network mode for network %s", containerID[:12], networkName))
		}

		// Check for exposed ports on public IP addresses
		for port := range containerJSON.Config.ExposedPorts {
			if isPublicIP(settings.IPAddress) {
				insecureNetworkSettings = append(insecureNetworkSettings, fmt.Sprintf("Container %s has exposed port %s on a public IP address (%s) on network %s", containerID[:12], port.Port(), settings.IPAddress, networkName))
			}
		}

		// Check for public IP addresses
		if settings.IPAddress != "" && settings.IPAddress != "127.0.0.1" {
			insecureNetworkSettings = append(insecureNetworkSettings, fmt.Sprintf("Container %s has IP address %s on network %s", containerID[:12], settings.IPAddress, networkName))
		}

		// Check for Docker socket exposure
		if _, err := os.Stat(filepath.Join("/", "var", "run", "docker.sock")); err == nil {
			insecureNetworkSettings = append(insecureNetworkSettings, fmt.Sprintf("Container %s has access to the Docker socket at /var/run/docker.sock", containerID[:12]))
		}
	}

	return insecureNetworkSettings, nil
}

// isPublicIP checks if an IP address is public
func isPublicIP(ip string) bool {
	return ip != "127.0.0.1" && !strings.HasPrefix(ip, "10.") && !strings.HasPrefix(ip, "192.168.") && !strings.HasPrefix(ip, "172.16.")
}

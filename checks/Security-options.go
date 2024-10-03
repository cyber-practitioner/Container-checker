// container-checker/checks/Security-options.go
package checks

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

// CheckSecurityOptions retrieves the security options for all containers and
// provides recommendations based on the security options.
func CheckSecurityOptions(cli *client.Client) (map[string][]string, map[string]string, error) {
	// List all containers
	containers, err := cli.ContainerList(context.Background(), container.ListOptions{All: true})
	if err != nil {
		return nil, nil, fmt.Errorf("error listing containers: %v", err)
	}

	securityOptions := make(map[string][]string)
	recommendations := make(map[string]string)

	for _, container := range containers {
		opts, err := checkContainerSecurityOptions(cli, container.ID)
		if err != nil {
			return nil, nil, fmt.Errorf("error checking security options for container %s: %v", container.ID[:12], err)
		}

		securityOptions[container.ID] = opts
		recommendations[container.ID] = getSecurityRecommendations(opts)
	}

	return securityOptions, recommendations, nil
}

func checkContainerSecurityOptions(cli *client.Client, containerID string) ([]string, error) {
	containerJSON, err := cli.ContainerInspect(context.Background(), containerID)
	if err != nil {
		return nil, err
	}

	// Get security options
	securityOpts := containerJSON.HostConfig.SecurityOpt

	return securityOpts, nil
}

func getSecurityRecommendations(securityOpts []string) string {
	var recommendations string

	if containsDisableLabel(securityOpts) {
		recommendations += "The 'label=disable' security option is enabled, which disables SELinux or AppArmor. This can reduce the security of the container. Consider enabling SELinux or AppArmor for better security.\n"
	}

	// Add more security recommendations based on the security options
	// For example, check for the presence of other security options and provide recommendations accordingly

	return recommendations
}

// containsDisableLabel checks if the security options contain "label=disable".
func containsDisableLabel(securityOpts []string) bool {
	for _, opt := range securityOpts {
		if opt == "label=disable" || opt == "" {
			return true
		}
	}
	return false
}

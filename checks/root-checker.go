package checks

import (
	"container-checker/utils"
	"context"
	"fmt"

	"github.com/docker/docker/client"
)

// CheckPrivilegedContainers lists and checks containers for privileged mode and security options.
func CheckPrivilegedContainers(cli *client.Client) ([]PrivilegedContainer, string, error) {
	// Check the Docker daemon's API version
	serverVersion, err := cli.ServerVersion(context.Background())
	if err != nil {
		return nil, "", err
	}

	fmt.Printf("Docker daemon version: %s\n", serverVersion.Version)

	// List all containers
	containers, err := utils.ListContainers(cli)
	if err != nil {
		return nil, "", err
	}

	if len(containers) == 0 {
		fmt.Println("No containers found.")
		return nil, "", nil
	}

	var privilegedContainers []PrivilegedContainer
	var nonPrivilegedContainers []NonPrivilegedContainer
	var recommendations string

	for _, container := range containers {
		containerJSON, err := cli.ContainerInspect(context.Background(), container.ID)
		if err != nil {
			return nil, "", err
		}

		if containerJSON.HostConfig.Privileged {
			fmt.Printf("Privileged Container Detected: ID: %s, Image: %s, Name: %s, Status: %s\n", container.ID[:12], container.Image, container.Names[:], container.Status)
			recommendations += fmt.Sprintf("Container %s has elevated privileges. Recommendation: Configure a user with low privileges or run Docker in rootless mode.\n", container.ID[:12])

			privilegedContainers = append(privilegedContainers, PrivilegedContainer{
				ID:     container.ID[:12],
				Image:  container.Image,
				Status: container.Status,
				Names:  container.Names[:],
			})
		} else {
			nonPrivilegedContainers = append(nonPrivilegedContainers, NonPrivilegedContainer{
				ID:            container.ID[:12],
				Image:         container.Image,
				Status:        container.Status,
				ContainerName: container.Names[0][1:],
			})
		}

		// Check if the container is running as root
		if containerJSON.Config.User == "root" {
			recommendations += fmt.Sprintf("Container %s is running as root. Recommendation: Configure a user with low privileges or run Docker in rootless mode.\n", container.ID[:12])
		}
	}

	if len(privilegedContainers) == 0 {
		fmt.Println("No privileged containers found.")
	}

	return privilegedContainers, recommendations, nil
}

type PrivilegedContainer struct {
	ID     string
	Image  string
	Status string
	Names  []string
}

type NonPrivilegedContainer struct {
	ID            string
	Image         string
	Status        string
	ContainerName string
}

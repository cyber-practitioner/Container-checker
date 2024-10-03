// all-checks.go
package checks

import (
	"container-checker/utils"
	"context"
	"fmt"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

// ContainerInfo holds the unified information for each container.
type ContainerInfo struct {
	ID                        string   `json:"id"`
	ContainerName             string   `json:"containerName"`
	IsRunningAsRoot           bool     `json:"isRunningAsRoot"`
	PrivilegedContainer       bool     `json:"privilegedContainer"`
	ReadOnlyRootFilesystem    bool     `json:"readOnlyRootFilesystem"`
	PrivilegedContainerImage  string   `json:"privilegedContainerImage"`
	PrivilegedContainerStatus string   `json:"privilegedContainerStatus"`
	SecurityOptions           []string `json:"securityOptions"`
	AdvancedCapabilities      []string `json:"advancedCapabilities"`
	RestartPolicy             string   `json:"restartPolicy"`
	MaxProcesses              string   `json:"maxProcesses"`
	Recommendations           string   `json:"recommendations"`
}

func hasAdvancedCapabilities(hostConfig *container.HostConfig) bool {
	// Check for advanced capabilities
	if hostConfig.CapAdd != nil && len(hostConfig.CapAdd) > 0 {
		for _, cap := range hostConfig.CapAdd {
			if isAdvancedCapability(cap) != false {
				return true
			}
		}
	}
	return false
}
func isAdvancedCapability(capability string) bool {
	// List of advanced capabilities that can pose a security risk
	advancedCapabilities := []string{
		"CAP_SYS_ADMIN",
		"CAP_NET_ADMIN",
		"CAP_NET_RAW",
		"CAP_MKNOD",
		"CAP_SETFCAP",
		"CAP_SETPCAP",
		"CAP_SYS_PTRACE",
		"CAP_SYS_MODULE",
	}

	for _, cap := range advancedCapabilities {
		if cap == capability {
			return true
		}
	}
	return false
}

// CheckAllContainers lists all containers and checks their privileged status, security options, and read-only root filesystem.
func CheckAllContainers(cli *client.Client) ([]ContainerInfo, error) {
	// List all containers
	containers, err := utils.ListContainers(cli)
	if err != nil {
		return nil, err
	}

	if len(containers) == 0 {
		return nil, fmt.Errorf("no containers found")
	}

	var containerInfo []ContainerInfo

	for _, container := range containers {
		containerJSON, err := cli.ContainerInspect(context.Background(), container.ID)
		if err != nil {
			return nil, err
		}

		isPrivileged := containerJSON.HostConfig.Privileged
		isRunningAsRoot := containerJSON.Config.User == "root"
		isReadOnlyRootFilesystem := containerJSON.HostConfig.ReadonlyRootfs
		hasAdvancedCapabilities := hasAdvancedCapabilities(containerJSON.HostConfig)
		capabilities := containerJSON.HostConfig.CapAdd
		restartPolicy := containerJSON.HostConfig.RestartPolicy.Name
		maxProcesses := containerJSON.HostConfig.Resources.PidsLimit

		var recommendations string

		if isPrivileged {
			recommendations = fmt.Sprintf("Configure a user Container %s  with low privileges or run the container in rootless mode.\n", container.ID[:12])
		}
		if isRunningAsRoot {
			recommendations += fmt.Sprintf("Configure Container %s a user with low privileges or run container in rootless mode.\n", container.ID[:12])
		}
		if isReadOnlyRootFilesystem {
			recommendations += fmt.Sprintf("\nConfigure Container %s  to use a read-only root filesystem.\n", container.ID[:12])
		}
		if containerJSON.HostConfig.SecurityOpt == nil {
			recommendations += fmt.Sprintf("\nContainer %s needs to be relaunched either with SECCOMP or use the --security-opt=no-new-privileges flag when running Docker images to prevent priviledge escalation.\n", container.ID[:12])
		}
		if containerJSON.HostConfig.CapAdd == nil {
			recommendations += fmt.Sprintf("\nContainer %s needs to be relaunched either with lower prvileges or use the --cap-add= flag when running Docker images to prevent priviledge escalation.\n", container.ID[:12])
		}
		if hasAdvancedCapabilities {
			recommendations += fmt.Sprintf("\nContainer %s has advanced capabilities that can pose a security risk: %v. Recommendation: Minimize the use of capabilities and run the container with the least privileges required.\n", container.ID[:12], containerJSON.HostConfig.CapAdd)
		}
		if capabilities != nil {
			recommendations += fmt.Sprintf("\nContainer %s has capabilities that can pose a security risk: %v. Recommendation: Minimize the use of capabilities and run the container with the least privileges required.\n", container.ID[:12], containerJSON.HostConfig.CapAdd)
		}

		if restartPolicy == "0" {
			recommendations += fmt.Sprintf("\nConfigure a restart policy of on-failure or no-restart for Container %s \n", container.ID[:12])
		}
		if maxProcesses != nil && *maxProcesses == 0 {
			recommendations += fmt.Sprintf("\nConfigure a maximum number of processes to prevent DOS attacks for Container %s \n", container.ID[:12])
		}

		var maxProcessesStr string
		if maxProcesses != nil {
			maxProcessesStr = fmt.Sprintf("%d", *maxProcesses)
		} else {
			maxProcessesStr = "no limt set"
		}

		info := ContainerInfo{
			ID:                        container.ID[:12],
			ContainerName:             container.Names[0],
			IsRunningAsRoot:           isRunningAsRoot,
			PrivilegedContainer:       isPrivileged,
			ReadOnlyRootFilesystem:    isReadOnlyRootFilesystem,
			PrivilegedContainerImage:  container.Image,
			PrivilegedContainerStatus: container.Status,
			SecurityOptions:           containerJSON.HostConfig.SecurityOpt,
			AdvancedCapabilities:      capabilities,
			RestartPolicy:             string(restartPolicy),
			MaxProcesses:              maxProcessesStr,
			Recommendations:           recommendations,
		}
		containerInfo = append(containerInfo, info)
	}

	return containerInfo, nil
}

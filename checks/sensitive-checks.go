package checks

/*
import (
	//"container-checker/utils"
	"context"
	"container-checker/utils"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

// SensitivePaths defines a map of sensitive files for both Linux and Windows
var SensitivePaths = map[string][]string{
	"linux": {
		"/etc/passwd",
		"/etc/shadow",
		"/etc/group",
		"/root",
		"/var/run/docker.sock",
		"/proc",
		"/sys",
	},
	"windows": {
		"C:\\Windows\\System32\\config\\SAM",
		"C:\\Windows\\System32\\config\\system",
		"C:\\Users\\Administrator",
		"C:\\Windows\\System32\\drivers\\etc\\hosts",
		"C:\\Windows\\System32\\winlogon.exe",
		"C:\\Windows\\System32\\config\\RegBack",
	},
}

// CheckSensitiveMounts checks if a container has mounted sensitive files
func CheckSensitiveMounts(cli *client.Client, container types.Container) bool {
	osType := utils.DetectOSType(cli, container.ID)
	sensitiveFiles := SensitivePaths[osType]

	containerJSON, err := cli.ContainerInspect(context.Background(), container.ID)
	if err != nil {
		fmt.Printf("Error inspecting container %s: %v\n", container.ID, err)
		return false
	}

	for _, mount := range containerJSON.Mounts {
		for _, path := range sensitiveFiles {
			if mount.Source == path {
				fmt.Printf("Sensitive file %s is mounted in container %s\n", path, container.ID)
				return true
			}
		}
	}
	return false
}
*/

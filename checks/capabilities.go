package checks

import (
	"context"
	"fmt"
	"log"

	"container-checker/utils"

	"github.com/docker/docker/client"
)

// CapabilityRecommendation represents a recommendation for a container with advanced capabilities.
type CapabilityRecommendation struct {
	ContainerID     string
	ContainerName   string
	Capabilities    []string
	Recommendations []string
	RiskScore       int
}

var capabilityRecommendations = map[string]string{
	"SYS_ADMIN":       "The container has the SYS_ADMIN capability, which grants broad privileges. Consider removing this capability or using a more restrictive set of capabilities.",
	"NET_ADMIN":       "The container has the NET_ADMIN capability, which grants broad network administration privileges. Consider removing this capability or using a more restrictive set of capabilities.",
	"IPC_LOCK":        "The container has the IPC_LOCK capability, which allows locking of shared memory segments. Consider removing this capability or using a more restrictive set of capabilities.",
	"DAC_OVERRIDE":    "The container has the DAC_OVERRIDE capability, which allows overriding file system permissions. Consider removing this capability or using a more restrictive set of capabilities.",
	"SYS_MODULE":      "The container has the SYS_MODULE capability, which allows loading and unloading kernel modules. This can be a significant security risk. Consider removing this capability or using a more restrictive set of capabilities.",
	"SYS_RAWIO":       "The container has the SYS_RAWIO capability, which allows direct access to I/O ports. This can be a security risk. Consider removing this capability or using a more restrictive set of capabilities.",
	"SYS_PTRACE":      "The container has the SYS_PTRACE capability, which allows tracing and inspecting other processes. This can be a security risk. Consider removing this capability or using a more restrictive set of capabilities.",
	"DAC_READ_SEARCH": "The container has the DAC_READ_SEARCH capability, which allows reading and searching files and directories that would otherwise be inaccessible. This can be a security risk. Consider removing this capability or using a more restrictive set of capabilities.",
}

var capabilityRiskScores = map[string]int{
	"SYS_ADMIN":       5,
	"NET_ADMIN":       4,
	"IPC_LOCK":        3,
	"DAC_OVERRIDE":    4,
	"SYS_MODULE":      5,
	"SYS_RAWIO":       4,
	"SYS_PTRACE":      4,
	"DAC_READ_SEARCH": 3,
}

// CheckAdvancedCapabilities inspects containers for advanced capabilities and provides recommendations.
func CheckAdvancedCapabilities(cli *client.Client) []CapabilityRecommendation {
	var recommendations []CapabilityRecommendation

	containers, err := utils.ListContainers(cli)
	if err != nil {
		log.Printf("Error listing containers: %v", err)
		return nil
	}

	for _, container := range containers {
		containerJSON, err := cli.ContainerInspect(context.Background(), container.ID)
		if err != nil {
			log.Printf("Error inspecting container %s: %v", container.ID, err)
			continue
		}

		capabilityMap := make(map[string]struct{})
		capabilityRecommendationMap := make(map[string]string)
		riskScore := 0

		for _, capability := range containerJSON.HostConfig.CapAdd {
			if _, ok := capabilityMap[capability]; !ok {
				capabilityMap[capability] = struct{}{}
				if recommendation, ok := capabilityRecommendations[capability]; ok {
					capabilityRecommendationMap[capability] = recommendation
					riskScore += capabilityRiskScores[capability]
				}
			}
		}

		capabilities := make([]string, 0, len(capabilityMap))
		for capability := range capabilityMap {
			capabilities = append(capabilities, capability)
		}

		capabilityRecommendations := make([]string, 0, len(capabilityRecommendationMap))
		for _, recommendation := range capabilityRecommendationMap {
			capabilityRecommendations = append(capabilityRecommendations, recommendation)
		}

		if len(capabilities) > 0 {
			rec := CapabilityRecommendation{
				ContainerID:     container.ID,
				ContainerName:   container.Names[0], // Add the container name
				Capabilities:    capabilities,
				Recommendations: capabilityRecommendations,
				RiskScore:       riskScore,
			}
			recommendations = append(recommendations, rec)
		}
	}

	return recommendations
}

// PrintCapabilityRecommendations prints the capability recommendations.
func PrintCapabilityRecommendations(recommendations []CapabilityRecommendation) {
	if len(recommendations) == 0 {
		fmt.Println("No advanced capabilities detected.")
		return
	}

	fmt.Println("Advanced Capability Recommendations:")
	for _, rec := range recommendations {
		fmt.Printf("- Container Name: %s\n", rec.ContainerName)
		fmt.Printf("  Container ID: %s\n", rec.ContainerID[:12])
		fmt.Println("  Capabilities:")
		for _, capability := range rec.Capabilities {
			fmt.Printf("    - %s\n", capability)
		}
		fmt.Println("  Recommendations:")
		for _, recommendation := range rec.Recommendations {
			fmt.Printf("    - %s\n", recommendation)
		}
		fmt.Printf("  Risk Score: %d\n", rec.RiskScore)
		fmt.Println()
	}

}

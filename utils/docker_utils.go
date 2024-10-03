package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"

	//"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
)

// PrintContainerJSON prints the detailed JSON representation of the containers.
func PrintContainerJSON(cli *client.Client) {
	// List all containers
	containers, err := cli.ContainerList(context.Background(), container.ListOptions{All: true})
	if err != nil {
		log.Fatalf("Error listing containers: %v", err)
	}

	for _, container := range containers {
		containerID := container.ID

		containerJSON, err := cli.ContainerInspect(context.Background(), containerID)
		if err != nil {
			fmt.Printf("Error inspecting container %s: %v\n", containerID, err)
			continue
		}

		// Print the container JSON in a formatted way
		containerJSONBytes, err := json.MarshalIndent(containerJSON, "", "  ")
		if err != nil {
			fmt.Printf("Error marshaling container JSON: %v\n", err)
			continue
		}

		fmt.Printf("Container JSON for %s:\n%s\n", containerID, string(containerJSONBytes))
	}
}

// InspectDockerImages lists and inspects all Docker images
func InspectDockerImages(cli *client.Client) {
	fmt.Println("Listing all Docker images...")
	images, err := cli.ImageList(context.Background(), image.ListOptions{})
	if err != nil {
		log.Fatal(err)
	}

	for _, img := range images {
		imageInspect, rawJSON, err := cli.ImageInspectWithRaw(context.Background(), img.ID)
		if err != nil {
			fmt.Printf("Failed to inspect image %s: %v\n", img.ID, err)
			continue
		}

		// Print out basic image details
		fmt.Printf("Image ID: %s\n", imageInspect.ID)
		fmt.Printf("Repo Tags: %v\n", imageInspect.RepoTags)
		fmt.Printf("Size: %d bytes\n", imageInspect.Size)
		fmt.Printf("Image ports exposed: %v\n", imageInspect.Config.ExposedPorts)
		fmt.Printf("OS: %s\n", imageInspect.Os)
		//fmt.Printf("Architecture: %s\n", imageInspect.Architecture)
		//fmt.Printf("Created: %s\n", imageInspect.Created)
		fmt.Printf("Docker Version: %s\n", imageInspect.DockerVersion)

		// Call the function to check for security risks
		//checks.CheckForSecurityRisks(imageInspect)

		// Print raw JSON for additional manual inspection
		fmt.Printf("Raw JSON: %s\n", string(rawJSON))
		fmt.Println("-------------------------------------------------")
	}
}

// getAllNetworks retrieves all networks available on the Docker host.
func getAllNetworks(cli *client.Client) (map[string]types.NetworkResource, error) {
	networks, err := cli.NetworkList(context.Background(), types.NetworkListOptions{})
	if err != nil {
		return nil, err
	}

	networkMap := make(map[string]types.NetworkResource)
	for _, nw := range networks {
		networkMap[nw.ID] = nw
	}

	return networkMap, nil
}

// PrintAllNetworks prints the details of all Docker networks.
func PrintAllNetworks(cli *client.Client) {
	networks, err := getAllNetworks(cli)
	if err != nil {
		log.Fatalf("Error retrieving networks: %v", err)
	}

	for id, network := range networks {
		// Print network details
		networkJSONBytes, err := json.MarshalIndent(network, "", "  ")
		if err != nil {
			fmt.Printf("Error marshaling network JSON for %s: %v\n", id, err)
			continue
		}

		fmt.Printf("Network JSON for %s:\n%s\n", id, string(networkJSONBytes))
		fmt.Println("-------------------------------------------------")
	}
}

// ListContainers returns a list of all containers.
func ListContainers(cli *client.Client) ([]types.Container, error) {
	// List all containers
	containers, err := cli.ContainerList(context.Background(), container.ListOptions{All: true})
	if err != nil {
		return nil, fmt.Errorf("error listing containers: %v", err)
	}

	return containers, nil
}

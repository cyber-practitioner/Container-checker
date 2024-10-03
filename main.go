package main

import (
	//"container-checker/checks"

	"container-checker/web"
	"fmt"
	// Import the checks package
	//	"github.com/docker/docker/api/types/container"
)

func main() {

	fmt.Println("Starting Web server...")
	if err := web.StartWebServer(); err != nil {
		panic(err)
	}

	fmt.Println("Starting container checks...")
	/*
		// Initialize Docker client
		cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
		if err != nil {
			log.Fatal(err)
		}

		//utils.InspectDockerImages(cli)
		//utils.PrintContainerJSON(cli)
		fmt.Println("Checking for containers running with root priviledges...")
		checks.CheckPrivilegedContainers(cli)
		//utils.PrintContainerJSON(cli)
		fmt.Println("Printing Security options...")
		//checks.CheckSecurityOptions(cli)
		//utils.PrintAllNetworks(cli)
		//checks.CheckContainerNetworks(cli)
		//checks.CheckAdvancedCapabilities(cli)
		//checks.PrintCapabilityRecommendations(checks.CheckAdvancedCapabilities(cli))
	*/
}

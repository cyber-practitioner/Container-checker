package main

import (
	"container-checker/checks"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/docker/docker/client"
)

// StartWebServer starts the web server and handles the container checking.
func StartWebServer(cli *client.Client) {
	// Start periodic container checks
	containerInfoChan := make(chan []checks.ContainerInfo)
	go periodicContainerCheck(cli, containerInfoChan)

	// Define template functions
	funcMap := template.FuncMap{
		"split": func(s, sep string) []string {
			return strings.Split(s, sep)
		},
	}

	// Parse template with function map
	tmpl := template.Must(template.New("index.html").Funcs(funcMap).ParseFiles("web/index.html"))

	// Start web server and define handler
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		containerInfo := <-containerInfoChan

		if err := tmpl.Execute(w, containerInfo); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	fmt.Println("Starting web server on http://localhost:8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}

// periodicContainerCheck periodically checks the containers and sends data to the web server.
func periodicContainerCheck(cli *client.Client, containerInfoChan chan<- []checks.ContainerInfo) {
	for {
		containerInfo, err := checks.CheckAllContainers(cli)
		if err != nil {
			log.Printf("Error checking containers: %v", err)
		} else {
			containerInfoChan <- containerInfo
		}
		time.Sleep(10 * time.Second) // Check containers every 10 seconds
	}
}

func main() {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatalf("Error creating Docker client: %v", err)
	}

	StartWebServer(cli)
}

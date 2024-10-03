package checks

import (
	"fmt"

	"github.com/docker/docker/api/types"
)

// SecurityChecker defines the interface for checking security risks
type SecurityChecker interface {
	CheckForSecurityRisks(imageInspect types.ImageInspect)
}

// CheckForSecurityRisks analyzes the image configuration for potential security issues
func CheckForSecurityRisks(imageInspect types.ImageInspect) {
	// Check if the image is running as root
	if imageInspect.Config.User == "" || imageInspect.Config.User == "root" {
		fmt.Println("Warning: Image is running as root. This can be a security risk.")
	}

	// Check the base image's Python version for known vulnerabilities
	for _, envVar := range imageInspect.Config.Env {
		if envVarContains(envVar, "PYTHON_VERSION") {
			fmt.Printf("Python Version: %s\n", envVar)
			checkPythonVersionSecurity(envVar)
		}
	}

	// List all environment variables
	fmt.Println("Environment Variables:")
	for _, envVar := range imageInspect.Config.Env {
		fmt.Printf("- %s\n", envVar)
	}

	// Check for exposed ports
	if len(imageInspect.Config.ExposedPorts) > 0 {
		fmt.Println("Warning: Exposed ports detected. Ensure they are necessary and secure.")
	}

	// Suggest running the container as a non-root user if not already
	if imageInspect.Config.User == "" || imageInspect.Config.User == "root" {
		fmt.Println("Recommendation: Consider configuring the container to run as a non-root user.")
	}

	// Suggest using a minimal base image if the size is large
	if imageInspect.Size > 500*1024*1024 { // 500MB threshold for example
		fmt.Println("Recommendation: Consider using a smaller base image to reduce the attack surface.")
	}
}

// envVarContains checks if an environment variable contains a specific key
func envVarContains(envVar, key string) bool {
	return len(envVar) > len(key) && envVar[:len(key)] == key
}

// checkPythonVersionSecurity checks for security issues in the Python version
func checkPythonVersionSecurity(pythonVersion string) {
	version := pythonVersion[len("PYTHON_VERSION="):]
	if version < "3.8" { // Arbitrary version check, update based on actual vulnerabilities
		fmt.Printf("Warning: Python version %s may have known vulnerabilities. Consider updating.\n", version)
	}
}

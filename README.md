# Container Checker

**Work in Progress**

Container Checker is a tool that helps you identify privileged and non-privileged containers running on your Docker host. It provides a web-based interface to display information about the containers, including their privileged status, security options, and recommendations for improving security.

## Features

- Displays a list of all containers running on the Docker host
- Identifies privileged containers and containers running as root
- Checks the security options for each container
- Provides recommendations for improving container security
- Periodically updates the container information in the web interface

## Usage

1. Ensure you have Docker installed on your system.
2. Clone the repository and navigate to the project directory:
4. Open your web browser and navigate to `http://localhost:8081` to access the Container Checker web interface.

## Architecture

The Container Checker application consists of the following components:

1. **all-checks.go**: This file contains the `CheckAllContainers` function, which is responsible for fetching information about all the containers running on the Docker host, including their privileged status, security options, and recommendations.

2. **server.go**: This file sets up the web server and handles the periodic container checks. It uses the `CheckAllContainers` function from `all-checks.go` to retrieve the container information and passes it to the `index.html` template for rendering.

3. **index.html**: This is the HTML template that displays the container information in the web interface.

## Future Improvements

- Implement pagination or batching to handle a large number of containers efficiently
- Add support for filtering and sorting the container list
- Provide more detailed security recommendations based on the container's configuration
- Integrate with a security monitoring or alerting system
- Add support for exporting the container information to a file or external system

## Contributing

This project is currently a work in progress. If you'd like to contribute, please feel free to submit a pull request or open an issue on the [GitHub repository](https://github.com/your-username/container-checker).

## License

This project is licensed under the [MIT License](LICENSE).

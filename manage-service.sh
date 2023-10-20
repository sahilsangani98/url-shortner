#!/bin/bash

# Define the Docker Compose file
DOCKER_COMPOSE_FILE="docker-compose.yml"

# Display usage information
usage() {
    echo "Usage: $0 [start|stop]"
    echo "  start - Start the Docker Compose services"
    echo "  stop  - Stop the Docker Compose services"
    exit 1
}

# Check for the correct number of arguments
if [ $# -ne 1 ]; then
    usage
fi

# Start or stop the Docker Compose services based on the argument
case "$1" in
    "start")
        echo "Starting Docker Compose services..."
        docker-compose -f $DOCKER_COMPOSE_FILE up -d
        ;;
    "stop")
        echo "Stopping Docker Compose services..."
        docker-compose -f $DOCKER_COMPOSE_FILE down --volumes --rmi all
        ;;
    *)
        usage
        ;;
esac

# Exit with success status
exit 0

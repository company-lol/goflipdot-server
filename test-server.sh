#!/bin/bash

# Default values
HOST="localhost"
PORT="8080"
PATTERN='[[1,0,1,0],[0,1,0,1],[1,0,1,0],[0,1,0,1]]'

# Function to display usage
usage() {
    echo "Usage: $0 [-h host] [-p port] [-f pattern_file]"
    echo "  -h host   : Specify the host (default: localhost)"
    echo "  -p port   : Specify the port (default: 8080)"
    echo "  -f file   : Specify a JSON file containing the pattern"
    echo "Without -f, a default 4x4 checkerboard pattern will be used."
    exit 1
}

# Parse command line options
while getopts "h:p:f:" opt; do
    case ${opt} in
        h )
            HOST=$OPTARG
            ;;
        p )
            PORT=$OPTARG
            ;;
        f )
            if [[ -f $OPTARG ]]; then
                PATTERN=$(cat $OPTARG)
            else
                echo "Error: File $OPTARG does not exist."
                exit 1
            fi
            ;;
        \? )
            usage
            ;;
    esac
done

# Construct the URL
URL="http://$HOST:$PORT/api/dots"

# Send the request
echo "Sending pattern to $URL"
curl -X POST -H "Content-Type: application/json" -d "$PATTERN" "$URL"
echo  # New line for better readability of output

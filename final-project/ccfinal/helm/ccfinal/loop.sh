#!/bin/bash

url="http://monitor:8080/api/server/all"  # Change this to the desired URL
interval=1  # Interval between requests in seconds

trap "echo Stopped sending requests. && exit" SIGINT  # Trap Ctrl+C to stop the script gracefully

while true; do
    response=$(curl -s -o /dev/null -w "%{http_code}" $url)
    echo "GET request sent to $url. Response status code: $response"
done

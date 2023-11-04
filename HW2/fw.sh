#!/bin/bash

# Try to get data from Redis
cached_data=$(redis-cli -h red GET "$CITY")

# If data is found in Redis, use it
if [ "$cached_data" != "null" ] && [ ! -z "$cached_data" ]; then
    echo "Data found in cache: $CITY:$cached_data c"
else
    # If data is not in cache, fetch it using curl
    data=$(curl --request GET \
        --url "https://weatherapi-com.p.rapidapi.com/current.json?q=$CITY" \
        --header 'X-RapidAPI-Host: weatherapi-com.p.rapidapi.com' \
        --header "X-RapidAPI-Key: $API_KEY")
    echo $data | jq -j --arg ttl "$TTL" '"\(.location.name) \($ttl) \(.current.temp_c)"' | xargs redis-cli -h red SETEX
    echo "Data fetched and cached: $CITY:$data c"
fi

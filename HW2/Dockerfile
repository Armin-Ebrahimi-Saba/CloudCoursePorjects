# install ubuntu alpine
FROM ubuntu
WORKDIR /usr/src
# install curl package
# Set an alias for your script
RUN apt-get update && \
    apt-get install -y curl redis jq python3; \
    echo "alias fw='/usr/src/fw.sh'" >> ~/.bashrc; 
# Copy scripts into the container
COPY . .
ENTRYPOINT ["python3", "listener.py"]

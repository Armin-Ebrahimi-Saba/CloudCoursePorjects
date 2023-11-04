# install ubuntu alpine
FROM ubuntu
WORKDIR /usr/src
# install curl package
RUN apt-get update && \
    apt-get install -y curl redis jq; \
    echo "alias fw='/usr/src/fw.sh'" >> ~/.bashrc; \
    apt-get install redis;
# Set an alias for your script
# RUN echo "alias fw='/usr/src/fw.sh'" >> ~/.bashrc
# Copy your script into the container
COPY . .
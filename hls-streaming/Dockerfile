FROM ubuntu:focal

RUN /usr/bin/apt-get update && \
    /usr/bin/apt-get install -y curl && \
    curl -fsSL https://deb.nodesource.com/setup_24.x -o nodesource_setup.sh && \
    /usr/bin/apt-get update && \
    /usr/bin/apt-get upgrade -y && \
    /usr/bin/apt-get install -y nodejs ffmpeg

WORKDIR /home/app

ENTRYPOINT ["bash"]


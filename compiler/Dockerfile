# Use an official base image with C++ and CMake installed
FROM ubuntu:20.04
ARG DEBIAN_FRONTEND=noninteractive

# Install necessary packages
RUN apt-get update && \
    apt-get install -y git build-essential software-properties-common lsb-release wget

# Install latest CMake from Kitware (CMake's maintainer)
RUN wget -O - https://apt.kitware.com/keys/kitware-archive-latest.asc | gpg --dearmor - > /usr/share/keyrings/kitware-archive-keyring.gpg
RUN echo "deb [signed-by=/usr/share/keyrings/kitware-archive-keyring.gpg] https://apt.kitware.com/ubuntu/ $(lsb_release -cs) main" | tee /etc/apt/sources.list.d/kitware.list >/dev/null
RUN apt-get update && apt-get install -y cmake

# Clone the repository and build the project
RUN git clone https://github.com/JyotinderSingh/ctok.git /ctok && \
    mkdir /ctok/build && \
    cd /ctok/build && \
    cmake .. && \
    cmake --build .

# Install Go (necessary to run the Go server)
RUN wget https://golang.org/dl/go1.17.3.linux-amd64.tar.gz && \
    tar -xzf go1.17.3.linux-amd64.tar.gz && \
    mv go /usr/local

ENV PATH="${PATH}:/usr/local/go/bin"

# Add the Go server code
ADD server.go /server.go

# Expose the port the server will run on
EXPOSE 8080

# Command to run the Go server
CMD ["go", "run", "/server.go"]

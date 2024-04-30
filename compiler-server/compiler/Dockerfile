# Use an official base image with C++ and CMake installed
FROM ubuntu:20.04 as builder
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
    cmake -DCMAKE_BUILD_TYPE=Release -DBUILD_SHARED_LIBS=OFF -DCMAKE_EXE_LINKER_FLAGS="-static" .. && \
    cmake --build .

RUN chmod +x /ctok/build/ctok

# Create an alpine image with the built object
FROM alpine:latest

# Install libc6-compat if there are minimal dependencies on glibc
RUN apk --no-cache add libc6-compat

# Set the Current Working Directory inside the container
WORKDIR /

# Copy the statically linked binary from the builder stage
COPY --from=builder /ctok/build/ctok /ctok

# Command to run the binary
CMD ["/ctok"]

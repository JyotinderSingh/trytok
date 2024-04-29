# Compiler Server

A Docker container with a lightweight server that compiles Tok code and returns the output.

## Build the docker image

```bash
docker build -t compiler-server .
```

If you're building on darwin (Mac M1/M2/M3), use the following command

```bash
docker build --platform linux/x86_64 -t compiler-server .
```

## Run the docker container

```bash
docker run -d -p 8080:8080 compiler-server
```

If you're running on darwin (Mac M1/M2/M3), use the following command

```bash
docker run --platform linux/x86_64 -d -p 8080:8080 compiler-server
```

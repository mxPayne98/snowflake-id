# Snowflake ID Generator

A flexible Snowflake ID generator implemented in Go. It provides a REST API to generate unique IDs and can be run both natively and within Docker. The Snowflake algorithm guarantees that IDs are ordered sequentially and can be generated at a high rate across multiple nodes without collision.

## Running the Program

### 1. Via CLI

First, ensure you have Go (version 1.20 or later) installed on your machine.

Navigate to the project directory and compile the Go program:

```bash
go build -o snowflake-generator
```

Run the compiled binary:

```bash
./snowflake-generator
```

This will start the API server, and you can generate an ID by accessing the `/generate-id` endpoint:

```bash
curl http://localhost:8080/generate-id
```

### 2. Via Docker

Make sure Docker is installed and the daemon is running on your machine.

Build the Docker image:

```bash
docker build -t snowflake-generator .
```

Run the Docker container:

```bash
docker run -p 8080:8080 snowflake-generator
```

Again, access the API via:

```bash
curl http://localhost:8080/generate-id
```

## Using the Snowflake Package

The Snowflake package can be imported and used independently in your Go programs. Here's a brief guide:

1. **Initialization**: Instantiate a new Snowflake generator by providing a worker ID. Ensure that the worker ID is unique across nodes to avoid ID collisions.

```go
import "path/to/snowflake-package"

sf, err := snowflake.NewSnowflake(workerId)
if err != nil {
    log.Fatal(err)
}
```

2. **Generate ID**: Once you have a Snowflake instance, you can generate IDs:

```go
id, err := sf.GenerateId()
if err != nil {
    log.Fatal(err)
}
fmt.Println(id)
```

Note: Make sure to handle errors appropriately, as shown above.

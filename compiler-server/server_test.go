package main

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// Sends a POST request to the compiler server and checks the response.
// The request contains a code snippet that plays a simple text-based game.
// The server should compile and run the code, and return the game output.
func TestE2ECodeRun(t *testing.T) {
	// Launch the compiler server container.
	ctx := context.Background()

	// Define the container request using a custom image
	containerRequest := testcontainers.ContainerRequest{
		Image:        "jyotindersingh/compiler-server", // Use your custom image
		ExposedPorts: []string{"6543:8080/tcp"},
		WaitingFor:   wait.ForListeningPort("8080/tcp"),
		Privileged:   true,
		Mounts: []testcontainers.ContainerMount{
			{
				Source: testcontainers.GenericVolumeMountSource{Name: "/var/run/docker.sock"},
				Target: "/var/run/docker.sock",
			},
		},
	}

	// Start the container
	compilerServer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: containerRequest,
		Started:          true,
	})
	assert.NoError(t, err, "Failed to start container")
	defer compilerServer.Terminate(ctx)

	// time.Sleep(1 * time.Minute) // Wait for the server to start

	// Send a POST request to the server
	url := "http://localhost:6543"
	code := `
	// Game initialization
	var currentRoom = "entry";
	var treasureRoom = "library";
	var gameRunning = true;
	var hasFoundTreasure = false;
	
	// Room descriptions
	fun describeRoom(room) {
			if (room == "entry") {
					print "You are in the entry hall. There are doors to the north and east.";
			} else if (room == "living room") {
					print "You are in a cozy living room. There is a door to the south and a staircase leading up.";
			} else if (room == "kitchen") {
					print "You are in the large kitchen. There are doors to the west and north.";
			} else if (room == "library") {
					print "You are in the dusty library. You see many books and a desk.";
					hasFoundTreasure = true;
			}
	}
	
	// Moving between rooms
	fun move(direction) {
			if (currentRoom == "entry" and direction == "north") {
					currentRoom = "kitchen";
					describeRoom(currentRoom);
			} else if (currentRoom == "entry" and direction == "east") {
					currentRoom = "living room";
					describeRoom(currentRoom);
			} else if (currentRoom == "kitchen" and direction == "north") {
					currentRoom = "library";
					describeRoom(currentRoom);
			}
	}
	
	// Main game sequence
	print "Starting Treasure Hunt Game...";
	describeRoom(currentRoom);
	move("north"); // Move to kitchen
	if (!hasFoundTreasure) {
			move("north"); // Move to library
	}
	
	if (hasFoundTreasure) {
			print "Congratulations! You have found the hidden treasure!";
	} else {
			print "No treasure found. End of the game.";
	}
	print "Thanks for playing Treasure Hunt!";
	`

	req, err := http.NewRequest("POST", url, bytes.NewBufferString(code))
	assert.NoError(t, err, "Failed to create HTTP request")

	req.Header.Set("Content-Type", "text/plain")

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err, "Failed to send HTTP request")
	defer resp.Body.Close()

	// Read and check the response body
	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err, "Failed to read response body")

	expected := "Starting Treasure Hunt Game...\nYou are in the entry hall. There are doors to the north and east.\nYou are in the large kitchen. There are doors to the west and north.\nYou are in the dusty library. You see many books and a desk.\nCongratulations! You have found the hidden treasure!\nThanks for playing Treasure Hunt!\n"
	assert.Equal(t, expected, string(body), "Unexpected response from server")
}

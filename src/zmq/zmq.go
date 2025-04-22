package zmqcomm

import (
	"fmt"
	zmq "github.com/pebbe/zmq4"
)

// SendMessage sends a ZMQ message to the specified target and port
func SendMessage(target string, port string, data []byte) error {
	// Create a new ZMQ context and request socket
	socket, err := zmq.NewSocket(zmq.REQ)
	if err != nil {
		return fmt.Errorf("Error creating ZMQ socket: %v", err)
	}
	defer socket.Close() // Ensure cleanup

	// Connect to the ZMQ endpoint
	address := fmt.Sprintf("tcp://%s:%s", target, port)
	err = socket.Connect(address)
	if err != nil {
		return fmt.Errorf("Error connecting to ZMQ: %v", err)
	}

	// Send message
	_, err = socket.SendBytes(data, 0)
	if err != nil {
		return fmt.Errorf("Error sending ZMQ message: %v", err)
	}

	return nil // Return nil if no errors occurred
}

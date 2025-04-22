package tcp

import (
	"fmt"
	"net"
)

// SendPacket sends a TCP packet to the specified target and port
func SendPacket(target string, port string, data []byte) error {
	// Establish TCP Connection
	conn, err := net.Dial("tcp", target+":"+port)
	if err != nil {
		return fmt.Errorf("Error Connecting: %v", err)
	}

	// Ensure FIN Flag on function finish
	defer conn.Close()

	// Send Message
	_, err = conn.Write(data)
	if err != nil {
		return fmt.Errorf("Error sending data: %v", err)
	}

	return nil // Return nil if no errors occurred
}

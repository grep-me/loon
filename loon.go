package main

import (
	"loon/src/fuzz"
	"loon/src/head"
	"loon/src/load"
	"fmt"
	"sync"
	"flag"
)

func main() {
	// Print Banner:
	head.PrintHeader()

	// Load CMD Line Args:
	ipAddress := flag.String("target", "127.0.0.1", "Set Target IP")
	port := flag.String("port", "6969", "Set Target Port")
	filePath := flag.String("path", "./bin", "Set location of binary generators")
	requestType := flag.String("type", "TCP", "Choose between TCP and ZMQ")

	// Parse the flags
	flag.Parse()

	// Initialize progress and mutex for synchronization
	progress := &fuzz.Progress{}
	mu := &sync.Mutex{}

	// Set up the output configuration (whether to show progress or not)
	stdOutput := &fuzz.Stdoutput{
		Config: fuzz.Config{
			Quiet: false, // Set to true for quiet mode
		},
	}

	// Get executable paths
	executablePaths, err := load.GetExecutables(*filePath)
	if err != nil {
		fmt.Println("Error getting executable paths:", err)
		return // Early exit on error
	}

	fmt.Println("[!] Number of Packet Generators:", len(executablePaths))

	if *requestType == "TCP" {
		// Call GeneratePackets with the list of paths
		err = fuzz.GeneratePacketsTCP(executablePaths, progress, mu, stdOutput, *ipAddress, *port)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	} 

	if *requestType == "ZMQ" {
		// Call GeneratePackets with the list of paths
		err = fuzz.GeneratePacketsZMQ(executablePaths, progress, mu, stdOutput, *ipAddress, *port)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	} 



	
}
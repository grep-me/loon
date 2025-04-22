package fuzz

import (
	"loon/src/tcp"
	"loon/src/zmq"
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"sync"
	"time"
)

// Progress holds the stats related to the fuzzing process
type Progress struct {
	ReqCount        int64
	ReqTotal        int64
	QueuePos        int64
	QueueTotal      int64
	StartedAt       time.Time
	ReqSec          int64
	ErrorCount      int64
	SuccessfulSends int64
	FailedSends     int64
}

// Config holds the configuration for the output (e.g., quiet mode)
type Config struct {
	Quiet bool
}

// Stdoutput handles printing progress updates
type Stdoutput struct {
	Config Config
}

// TERMINAL_CLEAR_LINE is used to clear the terminal line for progress updates
const TERMINAL_CLEAR_LINE = "\r"

// GeneratePackets captures the output of an executable and sends each line as a TCP packet
func GeneratePacketsZMQ(executablePaths []string, progress *Progress, mu *sync.Mutex, stdOutput *Stdoutput, targetHost string, targetPort string) error {	// Create a WaitGroup to wait for all goroutines to finish
	var wg sync.WaitGroup

	// Loop over each executable path
	for _, executablePath := range executablePaths {
		wg.Add(1)

		go func(executablePath string) {
			defer wg.Done()

			// Initialize StartedAt if it's not already set
			if progress.StartedAt.IsZero() {
				progress.StartedAt = time.Now()
			}

			// Set and Execute Command
			command := exec.Command(executablePath)

			// Create a pipe to capture the output
			stdoutPipe, err := command.StdoutPipe()
			if err != nil {
				stdOutput.Progress(*progress) // Print the progress before returning
				fmt.Printf("Error creating stdout pipe: %v\n", err)
				return
			}

			// Start the command
			err = command.Start()
			if err != nil {
				stdOutput.Progress(*progress) // Print the progress before returning
				fmt.Printf("Error starting command: %v\n", err)
				return
			}

			// Read the output line by line
			scanner := bufio.NewScanner(stdoutPipe)

			// Loop through each line of the output from the executable
			for scanner.Scan() {
				line := scanner.Text()

				// Convert the string to []byte
				data := []byte(line)

				// Send the TCP packet

				// if err := tcp.SendPacket(targetHost, targetPort, data); err != nil {

				if err := zmqcomm.SendMessage(targetHost, targetPort, data); err != nil {
					// Lock and update error count in progress
					mu.Lock()
					progress.ErrorCount++
					progress.FailedSends++
					mu.Unlock()

				} else {
					// Lock and update successful packet count in progress
					mu.Lock()
					progress.ReqCount++
					progress.SuccessfulSends++
					mu.Unlock()
				}

				// Update ReqSec (requests per second)
				mu.Lock()
				elapsed := time.Since(progress.StartedAt).Seconds()
				if elapsed > 0 {
					progress.ReqSec = int64(float64(progress.ReqCount) / elapsed)
				}
				mu.Unlock()

				// Print progress
				stdOutput.Progress(*progress)
			}

			// Check for errors during reading
			if err := scanner.Err(); err != nil {
				stdOutput.Progress(*progress)
				fmt.Printf("Error reading command output: %v\n", err)
			}

			// Wait for the command to finish
			err = command.Wait()
			if err != nil {
				stdOutput.Progress(*progress)
				fmt.Printf("Error waiting for command to finish: %v\n", err)
			}
		}(executablePath) // Pass executablePath to the goroutine
	}

	// Wait for all goroutines to finish
	wg.Wait()

	return nil
}

// GeneratePackets captures the output of an executable and sends each line as a TCP packet
func GeneratePacketsTCP(executablePaths []string, progress *Progress, mu *sync.Mutex, stdOutput *Stdoutput, targetHost string, targetPort string) error {	// Create a WaitGroup to wait for all goroutines to finish
	var wg sync.WaitGroup

	// Loop over each executable path
	for _, executablePath := range executablePaths {
		wg.Add(1)

		go func(executablePath string) {
			defer wg.Done()

			// Initialize StartedAt if it's not already set
			if progress.StartedAt.IsZero() {
				progress.StartedAt = time.Now()
			}

			// Set and Execute Command
			command := exec.Command(executablePath)

			// Create a pipe to capture the output
			stdoutPipe, err := command.StdoutPipe()
			if err != nil {
				stdOutput.Progress(*progress) // Print the progress before returning
				fmt.Printf("Error creating stdout pipe: %v\n", err)
				return
			}

			// Start the command
			err = command.Start()
			if err != nil {
				stdOutput.Progress(*progress) // Print the progress before returning
				fmt.Printf("Error starting command: %v\n", err)
				return
			}

			// Read the output line by line
			scanner := bufio.NewScanner(stdoutPipe)

			// Loop through each line of the output from the executable
			for scanner.Scan() {
				line := scanner.Text()

				// Convert the string to []byte
				data := []byte(line)

				// Send the TCP packet

				if err := tcp.SendPacket(targetHost, targetPort, data); err != nil {
					// Lock and update error count in progress
					mu.Lock()
					progress.ErrorCount++
					progress.FailedSends++
					mu.Unlock()

				} else {
					// Lock and update successful packet count in progress
					mu.Lock()
					progress.ReqCount++
					progress.SuccessfulSends++
					mu.Unlock()
				}

				// Update ReqSec (requests per second)
				mu.Lock()
				elapsed := time.Since(progress.StartedAt).Seconds()
				if elapsed > 0 {
					progress.ReqSec = int64(float64(progress.ReqCount) / elapsed)
				}
				mu.Unlock()

				// Print progress
				stdOutput.Progress(*progress)
			}

			// Check for errors during reading
			if err := scanner.Err(); err != nil {
				stdOutput.Progress(*progress)
				fmt.Printf("Error reading command output: %v\n", err)
			}

			// Wait for the command to finish
			err = command.Wait()
			if err != nil {
				stdOutput.Progress(*progress)
				fmt.Printf("Error waiting for command to finish: %v\n", err)
			}
		}(executablePath) // Pass executablePath to the goroutine
	}

	// Wait for all goroutines to finish
	wg.Wait()

	return nil
}

// Progress method updates the progress of fuzzing
func (s *Stdoutput) Progress(status Progress) {
	if s.Config.Quiet {
		// No progress for quiet mode
		return
	}

	// Calculate elapsed time
	dur := time.Since(status.StartedAt)
	runningSecs := int(dur / time.Second)
	var reqRate int64
	if runningSecs > 0 {
		reqRate = status.ReqSec
	} else {
		reqRate = 0
	}

	hours := dur / time.Hour
	dur -= hours * time.Hour
	mins := dur / time.Minute
	dur -= mins * time.Minute
	secs := dur / time.Second

	// Update the output on the same line
	// Clear the line before printing the new stats using \r (carriage return)
	fmt.Fprintf(os.Stderr, "%sDuration: [%02d:%02d:%02d] :: TCP Sent: %d :: Speed: %d/sec :: Err: %d", 
		TERMINAL_CLEAR_LINE, hours, mins, secs, status.ReqCount, reqRate, status.ErrorCount)

	// Ensure flushing to output stream
	os.Stderr.Sync()
}

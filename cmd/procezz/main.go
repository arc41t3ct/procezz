package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/fatih/color"
)

// ServiceConfig defines the structure for each service in the config file.
type ServiceConfig struct {
	Name    string `json:"name"`
	OnStart string `json:"on_start"`
	OnStop  string `json:"on_stop"`
}

// Config defines the overall structure for the configuration.
type Config struct {
	Services []ServiceConfig `json:"services"`
}

var (
	serviceConfigs Config
	lastStates     = make(map[string]bool)
	mu             sync.Mutex
)

// LoadConfig reads the configuration from a JSON file.
func LoadConfig(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &serviceConfigs)
}

// checkProcessExists checks if a process with the given name is running.
func checkProcessExists(name string) (bool, error) {
	cmd := exec.Command("pgrep", name) // Using pgrep to check if the process exists
	err := cmd.Run()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			// The process is not running
			if exitErr.ExitCode() == 1 {
				return false, nil
			}
			return false, err
		}
		return false, err
	}
	return true, nil
}

// executeCommand executes a shell command or script from the config.
func executeCommand(script string) error {
	color.Green("executing: %s", script)
	if script == "" {
		return nil // Do nothing if the command is empty
	}
	cmd := exec.Command("sh", "-c", script)
	return cmd.Run()
}

// monitorService checks the status of each service listed in the configuration.
func monitorService(service ServiceConfig) {
	mu.Lock()
	lastExists := lastStates[service.Name]
	mu.Unlock()
	exists, err := checkProcessExists(service.Name)
	if err != nil {
		fmt.Printf("Error checking process %s: %v\n", service.Name, err)
		return
	}
	mu.Lock()
	if exists && !lastExists {
		color.Yellow("Process %s has started", service.Name)
		if err = executeCommand(service.OnStart); err != nil {
			color.Red("Error executing on_start command for %s: %v", service.Name, err)
		}
	} else if !exists && lastExists {
		color.Yellow("Process %s has stopped", service.Name)
		if err = executeCommand(service.OnStop); err != nil {
			color.Red("Error executing on_stop command for %s: %v", service.Name, err)
		}
	}
	lastStates[service.Name] = exists // Update the last state
	mu.Unlock()
}

// validateInput validates the arguments being passed to Imperator
func validateInput() (string, error) {
	var arg1 string
	if len(os.Args) > 1 {
		arg1 = os.Args[1]
		if arg1 == "" {
			return "", errors.New("must specify config file")
		}
	}
	return arg1, nil
}

func usage() {
	color.Yellow(`Usage:
    procezz <full path to config>
  `)
}

func main() {
	arg1, err := validateInput()
	if err != nil {
		color.Red("Error parsing locating config: %s", err)
		usage()
		os.Exit(1)
	}
	// Load the configuration
	if err := LoadConfig(arg1); err != nil {
		color.Red("Error loading config: %v", err)
		usage()
		os.Exit(1)
	}
	// Initialize last states
	for _, service := range serviceConfigs.Services {
		lastStates[service.Name] = false
	}
	ticker := time.NewTicker(1 * time.Second) // Check every second
	defer ticker.Stop()
	done := make(chan bool)
	// Handle interrupt signals to cleanly exit the application
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c
		done <- true
	}()

	for {
		select {
		case <-done:
			color.Cyan("Shutting down...")
			return
		case <-ticker.C:
			// Monitor all services listed in the configuration
			for _, service := range serviceConfigs.Services {
				monitorService(service)
			}
		}
	}
}

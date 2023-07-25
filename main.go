package main

import (
	"flag"
	"github.com/go-ini/ini"
	"github.com/go-vgo/robotgo"
	"github.com/tarm/serial"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

// Define a structure to store the content of the configuration file
type Config struct {
	Baud int    `ini:"baud"`
	Port string `ini:"port"`
}

// Read key value pairs from configuration file
func readConfig(filename string) (map[string][]string, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(data), "\n")
	config := make(map[string][]string)
	for _, line := range lines {
		parts := strings.Split(line, "=")
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			if key == "" || value == "" {
				log.Println("Warning: empty key or value in config file")
				continue
			}
			keys := strings.Split(value, "+")
			config[key] = keys
		} else {
			log.Println("Warning: invalid format in config file")
			continue
		}
	}
	return config, nil
}

func main() {
	// Define Console Parameters
	baud := flag.Int("b", 0, "baud rate")
	port := flag.String("p", "", "port name")
	file := flag.String("f", "config.ini", "config file path")
	flag.Parse()

	// Check if the configuration file exists
	if _, err := os.Stat(*file); os.IsNotExist(err) {
		log.Panic("Error: config file does not exist")
		return
	}

	// Read Configuration File
	cfg := new(Config)
	err := ini.MapTo(cfg, *file)
	if err != nil {
		log.Panic("Error: config file is not valid ini format")
		return
	}

	// Determine the final Baud and serial port
	if *baud == 0 {
		*baud = cfg.Baud
	}
	if *port == "" {
		*port = cfg.Port
	}

	// Check whether the Baud and serial port are valid
	if *baud <= 0 {
		log.Panic("Error: baud rate is not valid")
		return
	}
	if *port == "" {
		log.Panic("Error: port name is not valid")
		return
	}

	// Open serial port connection
	c := &serial.Config{Name: *port, Baud: *baud}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Panic(err)
		return
	}
	defer s.Close()

	// Read key value pairs in the configuration file
	config, err := readConfig(*file)
	if err != nil {
		log.Panic(err)
		return
	}

	// Circular reading of serial port data
	for {
		buf := make([]byte, 128)
		n, err := s.Read(buf)
		if err != nil {
			log.Panic(err)
			break
		}
		if n > 0 {
			data := string(buf[:n])
			log.Println(data)
			// Send corresponding single or combination keys according to the configuration file
			if keys, ok := config[data]; ok {
				log.Println("Sending keys:", keys)
				for _, key := range keys {
					// Check if the key value exists in the key list supported by the robotgo library
					if robotgo.Keycode[key] == 0 {
						log.Println("Warning: key is not supported by robotgo library")
						continue
					}
				}
				robotgo.KeyTap(keys...)
			}
		}
	}
}


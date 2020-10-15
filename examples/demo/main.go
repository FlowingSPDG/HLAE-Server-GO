package main

import (
	"fmt"
	"log"
	"path"

	mirvpgl "github.com/FlowingSPDG/HLAE-Server-GO"
)

const (
	demDir = "./sample.dem"
)

// ExampleHandler for HLAE Server
func ExampleHandler(cmd string) {
	log.Printf("Received %s\n", cmd)
}

func main() {
	log.Println("Starting...")
	hlaeserver, err := mirvpgl.New(":65535", "/mirv")
	if err != nil {
		panic(err)
	}
	hlaeserver.RegisterHandler(ExampleHandler)
	go func() {
		err := hlaeserver.Start()
		if err != nil {
			panic(err)
		}
	}()

	// NOTE : enclose ws URL with double quotes...
	// mirv_pgl url "ws://localhost:65535/mirv"
	// mirv_pgl start
	// mirv_pgl datastart

	dem := path.Dir(demDir)
	cmd := fmt.Sprintf("playdemo \"%s\"", dem)
	hlaeserver.BroadcastRCON(cmd)
}

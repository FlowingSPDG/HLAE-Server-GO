package main

import (
	"fmt"

	mirvpgl "github.com/FlowingSPDG/HLAE-Server-GO"
	"github.com/c-bata/go-prompt"
)

// ExampleHandler for HLAE Server
func ExampleHandler(cmd string) {
	fmt.Printf("Received %s\n", cmd)
}

// ExampleCamHandler for cam datas
func ExampleCamHandler(cam *mirvpgl.CamData) {
	fmt.Printf("Received cam data %v\n", cam)
}

func completer(in prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{}
	return prompt.FilterHasPrefix(s, in.GetWordBeforeCursor(), true)
}

func main() {
	hlaeserver, err := mirvpgl.New(":65535", "/mirv")
	if err != nil {
		panic(err)
	}
	hlaeserver.RegisterHandler(ExampleHandler)
	hlaeserver.RegisterCamHandler(ExampleCamHandler)
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
	for {
		cmd := prompt.Input("HLAE >>> ", completer)
		if cmd == "exit" {
			break
		}
		hlaeserver.BroadcastRCON(cmd)
	}
}

package main

import (
	"log"

	mirvpgl "github.com/FlowingSPDG/HLAE-Server-GO"
	"github.com/c-bata/go-prompt"
)

var (
	hlaeserver *mirvpgl.HLAEServer
)

func init() {
	var err error
	hlaeserver, err = mirvpgl.New(":65535", "/mirv")
	if err != nil {
		panic(err)
	}
}

// ExampleHandler for HLAE Server
func ExampleHandler(cmd string) {
	log.Printf("Received %s\n", cmd)
}

// ExampleCamHandler for cam datas
func ExampleCamHandler(cam *mirvpgl.CamData) {
	log.Printf("Received cam data %v\n", cam)
}

func completer(in prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{}
	return prompt.FilterHasPrefix(s, in.GetWordBeforeCursor(), true)
}

func main() {
	hlaeserver.RegisterHandler(ExampleHandler)
	hlaeserver.RegisterCamHandler(ExampleCamHandler)
	go hlaeserver.Start()
	// NOTE : enclose ws URL with double quotes...
	// mirv_pgl url "ws://localhost:65535/mirv"
	// mirv_pgl start
	// mirv_pgl datastart
	for {
		cmd := prompt.Input("HLAE >>> ", completer)
		hlaeserver.SendRCON(cmd)
	}
}

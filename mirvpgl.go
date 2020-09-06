package mirvpgl

import (

	//"bufio"

	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"

	"github.com/gin-gonic/gin"
	"gopkg.in/olahol/melody.v1"
	// "strconv"
)

var (
	nullstr = byte('\x00')
)

// CamData Camera datas
type CamData struct {
	Time float32
	XPos float32
	YPos float32
	ZPos float32
	XRot float32
	YRot float32
	Zrot float32
	Fov  float32
}

// New Get new instance of HLAEServer
func New(host, path string) (*HLAEServer, error) {
	if host == "" || path == "" {
		return nil, fmt.Errorf("Empty path or host")
	}
	srv := &HLAEServer{
		host:     host,
		path:     path,
		melody:   nil,
		sessions: make([]*melody.Session, 0),
		engine:   nil,
	}
	m := melody.New()
	m.HandleMessageBinary(func(s *melody.Session, data []byte) {
		// TODO
		buf := bytes.NewBuffer(data)
		cmd, err := buf.ReadString(nullstr)
		if err != nil {
			if err == io.EOF {
				log.Println("EOF.")
			} else {
				log.Println("Failed to read string from buffer : ", err)
			}
			return
		}
		// TODO with l
		switch cmd {
		case "hello":
			log.Println("HLAE Client connection established...")
		case "version":
			var version uint32
			if err := binary.Read(buf, binary.LittleEndian, &version); err != nil {
				log.Println("Failed to read version message buffer : ", err)
				return
			}
			log.Println("Current Version :", version)
			// h.handleRequest(cmd)
		case "dataStop":
			// h.handleRequest(cmd)
		case "dataStart":
			// h.handleRequest(cmd)
		case "levelInit":
			mapname, err := buf.ReadString(nullstr)
			if err != nil {
				log.Println("Failed to read levelInit message buffer : ", err)
				return
			}
			log.Printf("map : %s", mapname) //
			// h.handleRequest(cmd)
		case "levelShutdown":
			// h.handleRequest(cmd)
		case "cam":
			camData := CamData{}

			if err := binary.Read(buf, binary.LittleEndian, &camData); err != nil {
				log.Println("Failed to read cam message buffer : ", err)
				return
			}
			log.Printf("CamData : %v\n", camData)
			// h.handleCamRequest(camdata)
			//
		case "gameEvent":
			//TODO. JSON
			// h.handleRequest(cmd)
		default:
			fmt.Println("Unknown message")
			// h.handleRequest(cmd)
		}
	})
	m.HandleConnect(func(s *melody.Session) {
		srv.sessions = append(srv.sessions, s)
		s.WriteBinary(commandToByteSlice("echo Hello World from Go!"))
	})
	m.HandleDisconnect(func(s *melody.Session) {
		// Remove session from session slice
		for k, v := range srv.sessions {
			if v == s {
				newsession := make([]*melody.Session, len(srv.sessions)-1)
				newsession = append(srv.sessions[:k], srv.sessions[k+1:]...)
				copy(srv.sessions, newsession)
				return
			}
		}
	})
	srv.melody = m

	r := gin.Default()
	r.GET(path, func(c *gin.Context) {
		m.HandleRequest(c.Writer, c.Request)
	})

	srv.engine = r
	return &HLAEServer{
		host:   host,
		path:   path,
		melody: m,
		engine: r,
	}, nil
}

// HLAEServer Main struct
type HLAEServer struct {
	host     string
	path     string
	melody   *melody.Melody
	sessions []*melody.Session
	engine   *gin.Engine

	handlers    []func(string)
	camhandlers []func(*CamData)
}

func commandToByteSlice(cmd string) []byte {
	length := len("exec") + 1 + len(cmd) + 1 // "exec" + (nullstr) + command + (nullstr)
	command := make([]byte, 0, length)
	command = append(command, []byte("exec")...)
	command = append(command, nullstr)
	command = append(command, []byte(cmd)...)
	command = append(command, nullstr)

	return command
}

// SendRCON command
func (h *HLAEServer) SendRCON(cmd string) error {
	command := commandToByteSlice(cmd)
	if err := h.melody.BroadcastBinary(command); err != nil {
		return err
	}
	return nil
}

// RegisterHandler to handle each requests
func (h *HLAEServer) RegisterHandler(handler func(string)) {
	if h.handlers == nil {
		h.handlers = make([]func(string), 0)
	}
	h.handlers = append(h.handlers, handler)
	log.Printf("Registered handler. Currently %d handlers are active\n", len(h.handlers))
}

// RegisterCamHandler to handle each requests
func (h *HLAEServer) RegisterCamHandler(handler func(*CamData)) {
	if h.camhandlers == nil {
		h.camhandlers = make([]func(*CamData), 0)
	}
	h.camhandlers = append(h.camhandlers, handler)
	log.Printf("Registered Camera handler. Currently %d handlers are active\n", len(h.handlers))
}

func (h *HLAEServer) handleRequest(cmd string) {
	log.Printf("Sending handler request for %d clients...\n", len(h.handlers))
	for i := 0; i < len(h.handlers); i++ {
		go h.handlers[i](cmd)
	}
}

func (h *HLAEServer) handleCamRequest(cam *CamData) {
	log.Printf("Sending camera handler request for %d clients...\n", len(h.handlers))
	for i := 0; i < len(h.handlers); i++ {
		go h.camhandlers[i](cam)
	}
}

// Start listening
func (h *HLAEServer) Start() {
	log.Printf("Listening on %s%s", h.host, h.path)
	if err := h.engine.Run(h.host); err != nil {
		panic(err)
	}
}

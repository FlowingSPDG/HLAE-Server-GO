package mirvpgl

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"gopkg.in/olahol/melody.v1"
)

const (
	mirvPglVersion = 2
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
	srv.melody = melody.New()
	srv.eventSerializer = newGameEventUnserializer(enrichments)

	srv.melody.HandleConnect(func(s *melody.Session) {

	})
	srv.melody.HandleMessageBinary(func(s *melody.Session, data []byte) {
		buf := bytes.NewBuffer(data)
		cmd, err := buf.ReadString(nullstr)
		if err != nil {
			if err == io.EOF {
				fmt.Println("EOF.")
			} else {
				fmt.Println("Failed to read string from buffer : ", err)
			}
			return
		}
		cmd = strings.ReplaceAll(cmd, string(nullstr), "")
		fmt.Println("Received cmd:", cmd)
		switch cmd {
		case "hello":
			fmt.Println("HLAE Client connection established...")
			var version uint32
			if err := binary.Read(buf, binary.LittleEndian, &version); err != nil {
				fmt.Println("Failed to read version message buffer : ", err)
				return
			}
			fmt.Println("Current Version :", version)
			if version != mirvPglVersion {
				return
			}
			srv.handleRequest(cmd)
		case "dataStop":
			fmt.Println("HLAE Client stopped sending data...")
			srv.handleRequest(cmd)
		case "dataStart":
			fmt.Println("HLAE Client started sending data...")
			srv.handleRequest(cmd)
		case "levelInit":
			mapname, err := buf.ReadString(nullstr)
			if err != nil {
				fmt.Println("Failed to read levelInit message buffer : ", err)
				return
			}
			fmt.Printf("map : %s", mapname)
			// srv.handleRequest(cmd) // Add mapname info...
		case "levelShutdown":
			fmt.Println("Received levelShutdown...")
			srv.handleRequest(cmd)
		case "cam":
			camData := &CamData{}
			if err := binary.Read(buf, binary.LittleEndian, camData); err != nil {
				fmt.Println("Failed to parse cam message buffer : ", err)
				return
			}
			srv.handleCamRequest(camData)
		case "gameEvent":
			fmt.Println("Received gameEvent data...")
			ev, err := srv.eventSerializer.Unserialize(buf)
			if err != nil {
				fmt.Println("Failed to parse event descriptions... ERR:", err)
				return
			}
			log.Printf("EVENT : %v\n", ev)
			srv.handleEventRequest(ev)
		default:
			fmt.Printf("Unknown message:[%s]\n", cmd)
			srv.handleRequest(cmd)
		}
	})
	srv.melody.HandleConnect(func(s *melody.Session) {
		srv.sessions = append(srv.sessions, s)
		fmt.Println("HLAE WebSocket client connected. Current clients:", len(srv.sessions))
		// s.WriteBinary(commandToByteSlice("echo Hello World from Go!"))
	})
	srv.melody.HandleDisconnect(func(s *melody.Session) {
		// Remove session from session slice
		for k, v := range srv.sessions {
			if v == s {
				newsession := make([]*melody.Session, len(srv.sessions)-1)
				newsession = append(srv.sessions[:k], srv.sessions[k+1:]...)
				srv.sessions = newsession
				fmt.Println("HLAE WebSocket client disconnected. Current clients : ", len(srv.sessions))
				return
			}
		}
	})

	gin.SetMode(gin.ReleaseMode)
	srv.engine = gin.Default()
	srv.engine.GET(path, func(c *gin.Context) {
		srv.melody.HandleRequest(c.Writer, c.Request)
	})

	return srv, nil
}

// HLAEServer Main struct
type HLAEServer struct {
	host            string
	path            string
	melody          *melody.Melody
	sessions        []*melody.Session
	engine          *gin.Engine
	eventSerializer *gameEventUnserializer

	handlers      []func(string)
	camhandlers   []func(*CamData)
	eventhandlers []func(*GameEventData)
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

// BroadcastRCON broadcast command
func (h *HLAEServer) BroadcastRCON(cmd string) error {
	command := commandToByteSlice(cmd)
	if err := h.melody.BroadcastBinary(command); err != nil {
		return err
	}
	return nil
}

// SendRCON Send RCON to specific client
func (h *HLAEServer) SendRCON(k int, cmd string) error {
	if len(h.sessions)-1 < k {
		command := commandToByteSlice(cmd)
		if err := h.melody.BroadcastBinary(command); err != nil {
			return err
		}
		return nil
	}

	return fmt.Errorf("Out of index slice")
}

// RegisterHandler to handle each requests
func (h *HLAEServer) RegisterHandler(handler func(string)) {
	if h.handlers == nil {
		h.handlers = make([]func(string), 0)
	}
	h.handlers = append(h.handlers, handler)
	fmt.Printf("Registered handler. Currently %d handlers are active\n", len(h.handlers))
}

// RegisterCamHandler to handle each requests
func (h *HLAEServer) RegisterCamHandler(handler func(*CamData)) {
	if h.camhandlers == nil {
		h.camhandlers = make([]func(*CamData), 0)
	}
	h.camhandlers = append(h.camhandlers, handler)
	fmt.Printf("Registered Camera handler. Currently %d handlers are active\n", len(h.handlers))
}

// RegisterEventHandler to handle each requests
func (h *HLAEServer) RegisterEventHandler(handler func(*GameEventData)) {
	if h.eventhandlers == nil {
		h.eventhandlers = make([]func(*GameEventData), 0)
	}
	h.eventhandlers = append(h.eventhandlers, handler)
	fmt.Printf("Registered event handler. Currently %d handlers are active\n", len(h.eventhandlers))
}

func (h *HLAEServer) handleRequest(cmd string) {
	fmt.Printf("Sending handler request for %d handlers...\n", len(h.handlers))
	for i := 0; i < len(h.handlers); i++ {
		go h.handlers[i](cmd)
	}
}

func (h *HLAEServer) handleCamRequest(cam *CamData) {
	fmt.Printf("Sending camera handler request for %d handlers...\n", len(h.handlers))
	for i := 0; i < len(h.handlers); i++ {
		go h.camhandlers[i](cam)
	}
}

func (h *HLAEServer) handleEventRequest(ev *GameEventData) {
	fmt.Printf("Sending event handler request for %d handlers...\n", len(h.eventhandlers))
	for i := 0; i < len(h.eventhandlers); i++ {
		go h.eventhandlers[i](ev)
	}
}

// Start listening
func (h *HLAEServer) Start() error {
	// fmt.Printf("Listening on %s%s\n", h.host, h.path)
	return h.engine.Run(h.host)
}

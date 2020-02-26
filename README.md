# HLAE-Server-GO
HLAE Server with Go implemention

# About
HLAE `mirv_pgl` command server implemention with Go.  
This package helps you to handle `mirv_pgl` command and their datas.

## mirv_pgl
`mirv_pgl` supports to remote-control CS:GO client by using WebSocket.  
you can handle client's camera information(position,rotation,fov,time) and game-events, and you can send commands like RCON.

Official NodeJS code : https://github.com/advancedfx/advancedfx/blob/master/misc/mirv_pgl_test/server.js


## Usage
see [examples](https://github.com/FlowingSPDG/HLAE-Server-GO/blob/master/examples/main.go).  
- 1,Launch HLAE Server by ``go run main.go``.  
- 2,Launch CSGO with HLAE.
- 3,type following commands below:  
```
mirv_pgl url "ws://localhost:65535/mirv";
mirv_pgl start;
mirv_pgl datastart;
```

once CS:GO client succeed to connect HLAE Server, you can send commands by typing window.

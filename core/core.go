// core implements shared functionality that both
// client and server can use.
//
// By exporting all the messages (Request and Response)
// it becomes very easy for the client to communicate
// back and forth with the server.
package core

import (
	"errors"
	"log"
	"os/exec"
	"strings"
	"sync"
	"time"
)

type Response struct {
	Message string
	Ok      bool
}

type Request struct {
	Name string
}

// HandlerName provider the name of the only
// method that `core` exposes via the RPC
// interface.
//
// This could be replaced by the use of the reflect
// package (e.g, `reflect.ValueOf(func).Pointer()).Name()`).
const HandlerName = "Handler.Execute"

// Handler holds the methods to be exposed by the RPC
// server as well as properties that modify the methods'
// behavior.
type Handler struct {

	// Sleep adds a little sleep between to the
	// method execution to simulate a time-consuming
	// operation.
	Sleep    time.Duration
	TakenSet map[string]bool
	Mu       sync.Mutex
}

// Execute is the exported method that a RPC client can
// make use of by calling the RPC server using `HandlerName`
// as the endpoint.
//
// It takes a Request and produces a Response if no error
// happens, possibly sleeping in between if a sleep is
// specified in Handler.
func (h *Handler) Execute(req Request, res *Response) (err error) {
	if req.Name == "" {
		err = errors.New("A name must be specified")
		return
	}
	h.Mu.Lock()
	if req.Name == "request" {
		res.Ok = false
		out, _ := exec.Command("ls", "/var/www/html/").Output()
		s := strings.Split(string(out), "\n")
		for _, fileName := range s {
			if fileName != "" && fileName[len(fileName)-4:] == "plot" {
				if _, exist := h.TakenSet[fileName]; !exist {
					res.Ok = true
					res.Message = fileName
					h.TakenSet[fileName] = true
					log.Printf("Serving %s\n", res.Message)
					break
				}
			}
		}
	} else if strings.Split(req.Name, " ")[0] == "delete" {
		exec.Command("rm", "/var/www/html/"+strings.Split(req.Name, " ")[1]).Output()
		res.Ok = true
	}

	h.Mu.Unlock()

	return
}

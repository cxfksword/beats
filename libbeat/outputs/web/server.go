/*
   package WebServer get the captured http data from ngnet,
   and send these data to frontend by websocket.

           chan                    +-----WebClient
   ngnet----------WebServer---------+-----WebClient
                                   +-----WebClient
*/
package web

import (
	"fmt"
	"net"
	"net/http"

	"github.com/cxfksword/beats/libbeat/logp"
	"golang.org/x/net/websocket"
)

type WebServer struct {
	eventChan       chan interface{}
	addr            string
	staticFileDir   string
	connectedClient map[*websocket.Conn]*WebClient
	eventBuffer     []interface{}
}

func (s *WebServer) websocketHandler(ws *websocket.Conn) {
	c := NewWebClient(ws, s)
	s.connectedClient[ws] = c
	go c.TransmitEvents()
	c.RecvAndProcessCommand()
	c.Close()
	delete(s.connectedClient, ws)
}

/*
   Dispatch the event received from ngnet to all clients connected with websocket.
*/
func (s *WebServer) DispatchEvent() {
	for ev := range s.eventChan {
		for _, c := range s.connectedClient {
			c.eventChan <- ev
		}
	}
}

func (s *WebServer) SendEvent(ev interface{}) {
	s.eventChan <- ev
}

/*
   If the flag '-s' is set and the browser sent a 'sync' command,
   the WebServer will push all the http message buffered in eventBuffer to
   the client.
*/
func (s *WebServer) Sync(c *WebClient) {
	for _, ev := range s.eventBuffer {
		c.eventChan <- ev
	}
}

/*
   Handle static files (.html, .js, .css).
*/
func (s *WebServer) handleStaticFile(w http.ResponseWriter, r *http.Request) {
	uri := r.RequestURI
	if uri == "/" {
		uri = "/index.html"
	}

	data, err := Asset("static" + uri)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	w.Write(data)
}

func (s *WebServer) Start() {
	go s.DispatchEvent()
	http.Handle("/data", websocket.Handler(s.websocketHandler))
	http.HandleFunc("/", s.handleStaticFile)

	logp.Info("web output listen on%s", s.addr)
	fmt.Printf("Please goto http://%s%s for details.\n", GetHostIp(), s.addr)
	err := http.ListenAndServe(s.addr, nil)
	if err != nil {
		logp.Err("can't start web output server: %v", err)
	}
}

func NewWebServer(addr string) *WebServer {
	s := new(WebServer)
	s.eventChan = make(chan interface{}, 1024)
	s.addr = addr
	s.connectedClient = make(map[*websocket.Conn]*WebClient)
	return s
}

func GetHostIp() string {
	ip := "127.0.0.1"

	addrs, _ := net.InterfaceAddrs()
	for _, a := range addrs {
		ipnet := net.ParseIP(a.String())
		if ipnet != nil && !ipnet.IsLoopback() && !ipnet.IsUnspecified() {
			if ipnet.To4() != nil {
				ip = ipnet.String()
				break
			}
		}
	}

	return ip
}

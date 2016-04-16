/*
   package WebServer get the captured http data from ngnet,
   and send these data to frontend by websocket.

           chan                    +-----WebClient
   ngnet----------WebServer---------+-----WebClient
                                   +-----WebClient
*/
package web

import (
	"encoding/json"

	"golang.org/x/net/websocket"
)

type WebClient struct {
	eventChan chan interface{}
	server    *WebServer
	ws        *websocket.Conn
}

func (c *WebClient) RecvAndProcessCommand() {
	for {
		var msg string
		err := websocket.Message.Receive(c.ws, &msg)
		if err != nil {
			return
		}
		if len(msg) > 0 {
			if msg == "sync" {
				c.server.Sync(c)
			}
		} else {
			panic("empty command")
		}
	}
}

func (c *WebClient) TransmitEvents() {
	for ev := range c.eventChan {

		//http: {"@timestamp":"2016-04-16T17:50:39.940Z","beat":{"hostname":"vagrant-ubuntu-trusty","name":"vagrant-ubuntu-trusty"},"bytes_in":73,"bytes_out":72785,"client_ip":"10.0.2.15","client_port":37553,"client_proc":"","client_server":"","console":"[HTTP] 01:50:39 114.67.54.12    GET  http://google.com/\n                                     ←\u001b[37m200\u001b[0m text/html; charset=utf-8 72KB","direction":"out","http":{"code":200,"content_length":71396,"phrase":"OK","request_headers":{"accept":"*/*","host":"google.com","user-agent":"curl/7.35.0"},"response_headers":{"cache-control":"no-cache","content-type":"text/html; charset=utf-8","date":"Sat, 16 Apr 2016 17:50:35 GMT","expires":"Mon, 26 Jul 1997 05:00:00 GMT","server":"nginx/1.8.0","set-cookie":"WAm_vipUINFO=newCustomerGuide%3AA; path=/; domain=.vip.com; Max-Age=86313601; path=/; domain=.vip.com-Age=63072001; path=/; domain=.vip.com","transfer-encoding":"chunked","vary":"Accept-Encoding"}},"ip":"114.67.54.12","method":"GET","params":"","path":"/","port":80,"proc":"","query":"GET /","request_uri":"http://google.com/","responsetime":112,"server":"","status":"OK","type":"http"}"
		//thrift: {"@timestamp":"2016-04-16T17:40:08.412Z","beat":{"hostname":"vagrant-ubuntu-trusty","name":"vagrant-ubuntu-trusty"},"bytes_in":31,"bytes_out":37,"client_ip":"127.0.0.1","client_port":35742,"client_proc":"","client_server":"vagrant-ubuntu-trusty","console":"[Thrift] 01:40:08 127.0.0.1:9090      0ms   sayHello(1: \"TOM\")","ip":"127.0.0.1","method":"sayHello","path":"","port":9090,"proc":"","query":"sayHello(1: \"TOM\")","responsetime":0,"server":"vagrant-ubuntu-trusty","status":"OK","thrift":{"params":"(1: \"TOM\")","return_value":"\"Hello TOM\""},"type":"thrift"}"

		data, err := json.Marshal(ev)
		if err == nil {
			websocket.Message.Send(c.ws, string(data))
		}
	}
}

func (c *WebClient) Close() {
	close(c.eventChan)
}

func NewWebClient(ws *websocket.Conn, server *WebServer) *WebClient {
	c := new(WebClient)
	c.server = server
	c.ws = ws
	c.eventChan = make(chan interface{}, 16)
	return c
}

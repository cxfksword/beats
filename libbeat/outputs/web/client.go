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
	"time"

	"github.com/cxfksword/beats/libbeat/common"
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

		//{"@timestamp":"2016-04-10T15:56:39.406Z","beat":{"hostname":"vagrant-ubuntu-trusty","name":"vagrant-ubuntu-trusty"},"bytes_in":73,"bytes_out":386,"client_ip":"10.0.2.15","client_port":52716,"client_proc":"","client_server":"","console":"[HTTP]2016-04-10 15:56:39 200 GET http://baidu.com/","direction":"out","http":{"code":200,"content_length":81,"phrase":"OK"},"ip":"180.149.132.47","method":"GET","params":"","path":"/","port":80,"proc":"","query":"GET /","request_uri":"http://baidu.com/","responsetime":68,"server":"","status":"OK","type":"http"}
		//
		request := HttpRequestEvent{}
		request.SetType("HttpRequest")
		request.AddHeader(HttpHeaderItem{"Host", string(ev.(common.MapStr)["request_uri"].(common.NetString))})
		request.AddHeader(HttpHeaderItem{"Method", string(ev.(common.MapStr)["method"].(common.NetString))})
		request.SetTimestamp(float64(time.Now().Unix()))
		request.SetEndTimestamp(float64(time.Now().Unix()) + float64(20))
		data, err := json.Marshal(request)
		if err == nil {
			websocket.Message.Send(c.ws, string(data))
		}

		response := HttpResponseEvent{}
		response.SetType("HttpResponse")
		response.SetBody("ok")
		data, err = json.Marshal(response)
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

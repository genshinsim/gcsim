package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 100000000
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     originCheck,
}

//originCheck checks origin header re websocket upgrade
func originCheck(r *http.Request) bool {
	//origin should be the localhost at port 80
	origin := r.Header.Get("Origin")
	log.Println(origin)
	return origin == fmt.Sprintf("http://localhost:%v", port) || origin == "http://localhost:3000"
}

//handleUpgrade returns a handlerfunc that upgrades a websocket connection
func (s *Server) handleUpgrade() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.Log.Debugw("got request for socket")

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			//handle error?
			s.Log.Warnw("handleUpgrade: error upgrading", "err", err)
			return
		}

		//pullout token here?
		client := wsClient{
			send: make(chan []byte, 5),
			conn: conn,
			s:    s,
		}
		s.wsRegister <- client
		go client.readPump()
		go client.writePump()
	}
}

type wsRequest struct {
	ID      string    `json:"id"`      //id used by the client - we don't care; client's problem to keep track
	Type    string    `json:"type"`    //type of request
	Target  string    `json:"target"`  //target of request
	Payload string    `json:"payload"` //request payload
	client  *wsClient //response channel
}

type wsResponse struct {
	ID      string `json:"id"`     //id referencing whatever came from client
	Status  int    `json:"status"` //status code - use http statuses
	Payload string `json:"payload"`
}

type wsClient struct {
	conn *websocket.Conn
	send chan []byte
	s    *Server
}

func (s *Server) wsHub() {
	clients := make(map[*websocket.Conn]chan []byte)
	for {
		select {
		case c := <-s.wsRegister:
			clients[c.conn] = c.send
		case c := <-s.wsUnregister:
			if send, ok := clients[c]; ok {
				delete(clients, c)
				close(send)
				//use this to gracefully close down sockets esp to clean up open db
			}
		}
	}
}

func (w *wsClient) readPump() {
	defer func() {
		w.s.wsUnregister <- w.conn
		w.conn.Close()
	}()
	w.conn.SetReadLimit(maxMessageSize)
	w.conn.SetReadDeadline(time.Now().Add(pongWait))
	w.conn.SetPongHandler(func(string) error { w.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		mt, in, err := w.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				w.s.Log.Warnw("readPump abnormal err", "err", err)
			} else {
				w.s.Log.Warnw("readPump err", "err", err)
			}
			break
		}
		if mt == websocket.TextMessage || mt == websocket.BinaryMessage {
			var x wsRequest
			x.client = w
			if err := json.Unmarshal(in, &x); err == nil {
				// w.s.Log.Infow("readPump: got request", "req", x)
				go w.s.reqHandler(x)
			} else {
				w.s.Log.Warnw("error parsing msg", "err", err, "msg", string(in))
			}
		}
	}
}

func (w *wsClient) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		w.conn.Close()
	}()
	for {
		select {
		case msg, ok := <-w.send:
			w.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				w.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			writer, err := w.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			writer.Write(msg)
			//flush the queue if any
			n := len(w.send)
			for i := 0; i < n; i++ {
				writer.Write(<-w.send)
			}

			//check if closed
			if err := writer.Close(); err != nil {
				return
			}

		case <-ticker.C:
			w.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := w.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

//reqHandler redirects incoming requests
func (s *Server) reqHandler(r wsRequest) {
	//set up context
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	switch r.Type {
	//do stuff here based on request type
	case "run":
		s.handleRun(ctx, r)
	default:
		//do nothing
	}
}

//socketWorker handles one websocket connection
// func (s *Server) socketWorker(conn net.Conn) {
// 	s.Log.Debugw("socket opened")

// 	msg := make(chan []byte, 5)
// 	for {
// 		//read msg first?
// 		in, op, err := wsutil.ReadClientData(conn)
// 		if err != nil {
// 			s.Log.Debugw("error reading connection: ", "err", err, "code", op)
// 			// handle error
// 		}
// 		s.Log.Debugw("incoming msg received: ", "in", in, "code", op)
// 		//handle request here - add it to queue and toss it to our worker
// 		//request should conform with {id: "blah", payload: "blah blah"}
// 		//if not just ignore it

// 		if op != ws.OpContinuation {
// 			var x wsRequest
// 			// x.resp = msg
// 			if err := json.Unmarshal(in, &x); err == nil {
// 				go s.reqHandler(x)
// 			} else {
// 				s.Log.Debugw("error parsing in", "err", err)
// 			}
// 		}

// 		select {
// 		//gracefully shut down
// 		case <-s.wsClosed:
// 			//i think need to close socket here somehow
// 			return
// 		case m := <-s.wsBroadcast:
// 			s.Log.Debugw("sending out broadcast msg: ", "msg", string(m))
// 			err = wsutil.WriteServerMessage(conn, ws.OpBinary, m)
// 			if err != nil {
// 				// handle error
// 			}
// 		case m := <-msg:
// 			//send msg
// 			s.Log.Debugw("sending out single msg: ", "msg", string(m))
// 			err = wsutil.WriteServerMessage(conn, ws.OpBinary, m)
// 			if err != nil {
// 				// handle error
// 			}
// 		}
// 	}
// }

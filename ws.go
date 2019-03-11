package main

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
	"github.com/labstack/gommon/random"
	"net/http"
)

type WebSocketMessage struct {
	Event   string      `json:"event"`
	Payload interface{} `json:"payload"`
}

type SocketID string

type SocketEventListener func(interface{})

type Socket struct {
	conn      *websocket.Conn
	listeners map[string]SocketEventListener
	logger    echo.Logger
}

// wait reads messages from socket connection
func (s *Socket) wait() error {
	for {
		_, msg, err := s.conn.ReadMessage()
		if err != nil {
			return err
		}
		var message WebSocketMessage
		if err := json.Unmarshal(msg, &message); err != nil {
			return err
		}
		listener, ok := s.listeners[message.Event]
		if !ok {
		} else {
			listener(message.Payload)
		}
	}
}

// On todo
func (s *Socket) On(event string, callback SocketEventListener) error {
	s.listeners[event] = callback
	return nil
}

// Emit writes message to current connection
func (s *Socket) Emit(event string, payload interface{}) error {
	err := s.conn.WriteJSON(WebSocketMessage{
		Event:   event,
		Payload: payload,
	})
	if err != nil {
		if err == websocket.ErrCloseSent {
			return nil
		}
	}
	return nil
}

type WebSocket struct {
	upgrader     websocket.Upgrader
	connected    map[SocketID]*Socket
	OnConnection func(socket *Socket)
}

// Broadcast writes message to all connections
func (ws *WebSocket) Broadcast(event string, payload interface{}) error {
	for _, socket := range ws.connected {
		if err := socket.Emit(event, payload); err != nil {

		}
	}
	//	if err == websocket.ErrCloseSent {
	//		return nil
	//	}
	return nil
}

// Handler handles incoming socket connections
func (ws *WebSocket) Handler(c echo.Context) error {
	conn, err := ws.upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer conn.Close()

	socket := ws.addConnection(conn)

	socket.On("ping", func(payload interface{}) {
		socket.Emit("ping_result", "pong")
	})

	return socket.wait()
}

func (ws *WebSocket) addConnection(conn *websocket.Conn) *Socket {
	socket := Socket{
		conn:      conn,
		listeners: make(map[string]SocketEventListener),
	}
	ws.connected[SocketID(random.String(16))] = &socket
	return &socket
}

func (ws *WebSocket) removeConnection(id SocketID) error {
	delete(ws.connected, id)
	return nil
}

func NewWebSocket() *WebSocket {
	return &WebSocket{
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		connected: make(map[SocketID]*Socket),
	}
}

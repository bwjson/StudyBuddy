package delivery

import (
	"context"
	"github.com/bwjson/StudyBuddy/internal/dto"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

const (
	templateDir = "./web/templates/html"
	staticDir   = "./web/static/"
)

type WSServer interface {
	Start() error
	Stop() error
}

type wsSrv struct {
	mux       *http.ServeMux
	srv       *http.Server
	wsUpg     *websocket.Upgrader
	broadcast chan *dto.WsMessage
	clients   clients
}

type clients struct {
	mutex     *sync.RWMutex
	wsClients map[*websocket.Conn]struct{}
}

func NewWsServer(addr string) WSServer {
	m := http.NewServeMux()
	return &wsSrv{
		mux: m,
		srv: &http.Server{
			Addr:    addr,
			Handler: m,
		},
		wsUpg: &websocket.Upgrader{},
		clients: clients{
			mutex:     &sync.RWMutex{},
			wsClients: map[*websocket.Conn]struct{}{},
		},
		broadcast: make(chan *dto.WsMessage),
	}
}

func (ws *wsSrv) Start() error {
	ws.mux.Handle("/", http.FileServer(http.Dir(templateDir)))
	ws.mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(staticDir))))
	ws.mux.HandleFunc("/ws", ws.wsHandler)
	go ws.writeToClientsBroadcast()
	return ws.srv.ListenAndServe()
}

func (ws *wsSrv) Stop() error {
	log.Info("Before", ws.clients.wsClients)
	close(ws.broadcast)
	ws.clients.mutex.Lock()
	for conn := range ws.clients.wsClients {
		if err := conn.Close(); err != nil {
			log.Errorf("Error with closing: %v", err)
		}
		delete(ws.clients.wsClients, conn)
	}
	ws.clients.mutex.Unlock()
	log.Info("After close", ws.clients.wsClients)
	return ws.srv.Shutdown(context.Background())
}

func (ws *wsSrv) wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := ws.wsUpg.Upgrade(w, r, nil)
	if err != nil {
		log.Errorf("Error with websocket connection: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Infof("Client with address %s connected", conn.RemoteAddr().String())
	ws.clients.mutex.Lock()
	ws.clients.wsClients[conn] = struct{}{}
	ws.clients.mutex.Unlock()
	go ws.readFromClient(conn)
}

func (ws *wsSrv) readFromClient(conn *websocket.Conn) {
	for {
		msg := new(dto.WsMessage)
		if err := conn.ReadJSON(msg); err != nil {
			wsErr, ok := err.(*websocket.CloseError)
			if !ok || wsErr.Code != websocket.CloseGoingAway {
				log.Errorf("Error with reading from WebSocket: %v", err)
			}
			break
		}
		host, _, err := net.SplitHostPort(conn.RemoteAddr().String())
		if err != nil {
			log.Errorf("Error with address split: %v", err)
		}
		msg.IPAddress = host
		msg.Time = time.Now().Format("15:04")
		ws.broadcast <- msg
	}
	ws.clients.mutex.Lock()
	delete(ws.clients.wsClients, conn)
	ws.clients.mutex.Unlock()
}

func (ws *wsSrv) writeToClientsBroadcast() {
	for msg := range ws.broadcast {
		ws.clients.mutex.RLock()
		for client := range ws.clients.wsClients {
			func() {
				if err := client.WriteJSON(msg); err != nil {
					log.Errorf("Error with writing message: %v", err)
				}
			}()
		}
		ws.clients.mutex.RUnlock()
	}
}

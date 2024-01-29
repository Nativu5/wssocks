package wss

import (
	"context"
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"
	"nhooyr.io/websocket"
)

type WebsocksServerConfig struct {
	EnableHttp       bool
	EnableConnKey    bool   // bale connection key
	ConnKey          string // connection key
	EnableStatusPage bool   // enable/disable status page

	EnableTWFID bool // enable/disable TWFID display
}

type ServerWS struct {
	config WebsocksServerConfig
	hc     *HubCollection
}

// return a a function handling websocket requests from the peer.
func NewServeWS(hc *HubCollection, config WebsocksServerConfig) *ServerWS {
	return &ServerWS{config: config, hc: hc}
}

func (s *ServerWS) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// check if upgrade
	if r.Header.Get("Upgrade") != "websocket" {
		// check if enable TWFID display
		if s.config.EnableTWFID {
			cookie, err := r.Cookie("TWFID")
			if err == nil {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("Find TWFID: " + cookie.Value + "\n"))
			} else {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte("TWFID not found!\n"))
			}
			return
		}

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad Request, expecting upgrade\n"))
		return
	}

	// check connection key
	if s.config.EnableConnKey && r.Header.Get("Key") != s.config.ConnKey {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Access denied!\n"))
		return
	}

	wc, err := websocket.Accept(w, r, nil)
	if err != nil {
		log.Error(err)
		return
	}
	defer wc.Close(websocket.StatusNormalClosure, "the sky is falling")

	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	// negotiate version with client.
	if err := NegVersionServer(ctx, wc, s.config.EnableStatusPage); err != nil {
		return
	}

	hub := s.hc.NewHub(wc)
	defer s.hc.RemoveProxy(hub.id)
	defer hub.Close()
	// read messages from webSocket
	wc.SetReadLimit(1 << 23) // 8 MiB
	for {
		msgType, p, err := wc.Read(ctx) // fixme context
		// if WebSocket is closed by some reason, then this func will return,
		// and 'done' channel will be set, the outer func will reach to the end.
		if err != nil && err != io.EOF {
			log.Error("error reading webSocket message:", err)
			break
		}
		if err = dispatchMessage(hub, msgType, p, s.config); err != nil {
			log.Error("error proxy:", err)
			// break skip error
		}
	}
}

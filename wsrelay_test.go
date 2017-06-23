package wshttp

import (
	"testing"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"github.com/gorilla/websocket"
)

func fakeHandleWs(rw http.ResponseWriter, r *http.Request){
	wsConnection, err := upgrader.Upgrade(rw, r, nil)
	if err != nil{
		log.Fatalf("cannot upgrade websocket: %v", err)
	}
	defer wsConnection.Close()
	msgType, msg, err := wsConnection.ReadMessage()
	if err != nil{
		log.Fatalf("cannot read websocket: %v", err)
	}
	err = wsConnection.WriteMessage(msgType, msg)
	if err != nil{
		log.Fatalf("cannot write websocket: %v", err)
	}
}

func TestWsRelay_HandleWS(t *testing.T) {
	backEndServer := httptest.NewServer(http.HandlerFunc(fakeHandleWs))
	defer backEndServer.Close()
	backEndUrl, _ := url.Parse(backEndServer.URL)
	wsRelay := WsRelay{TargetAddr: backEndUrl.Host}
	relayServer := httptest.NewServer(http.HandlerFunc(wsRelay.HandleWS))
	defer relayServer.Close()
	relayUrl, _ := url.Parse(relayServer.URL)
	relayUrl.Scheme = "ws"
	clientConnection, _, err := websocket.DefaultDialer.Dial(relayUrl.String(), nil)
	if err != nil{
		t.Errorf("ws cannot make handshake: %v", err)
	}
	defer clientConnection.Close()
	err = clientConnection.WriteMessage(websocket.TextMessage, []byte("test"))
	if err != nil{
		t.Errorf("ws cannot write message: %v", err)
	}
	_, msg, err := clientConnection.ReadMessage()
	if err != nil{
		t.Errorf("ws cannot read message: %v", err)
	}
	if string(msg) != "test"{
		t.Errorf("ws msg mismatch: %v", string(msg))
	}
}

package wshttp

import (
	"testing"
	"net/http"
	"github.com/gorilla/websocket"
	"net/http/httptest"
	"net/url"
	"log"
)

type fakeWsHandler struct {}
var upgrader = websocket.Upgrader{}
func (f fakeWsHandler)HandleWs(rw http.ResponseWriter, r *http.Request){
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

type fakeHttpHandler struct {}
func (f fakeHttpHandler)HandleHttp(rw http.ResponseWriter, r *http.Request){
	rw.WriteHeader(http.StatusOK)
}

func TestDmux(t *testing.T) {
	fh := fakeHttpHandler{}
	fw := fakeWsHandler{}
	testServer := httptest.NewServer(http.HandlerFunc(Dmux(fh, fw)))
	defer testServer.Close()
	u, _ := url.Parse(testServer.URL)
	u.Scheme = "ws"
	clientConnection, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
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
	u.Scheme = "http"
	req, _ := http.NewRequest(http.MethodGet, u.String(), nil)
	resp, _ := http.DefaultClient.Do(req)
	if resp.StatusCode != http.StatusOK{
		t.Errorf("got http status code: %v", resp.StatusCode)
	}
}
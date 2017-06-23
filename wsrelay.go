package wshttp

import (
	"net/http"
	"net"
	"io"
)

type WsRelay struct{
	TargetAddr string
}

func (ws WsRelay) HandleWS(rw http.ResponseWriter, r *http.Request){
	hijacker, ok := rw.(http.Hijacker)
	if !ok{
		http.Error(rw, "Internal server error", http.StatusInternalServerError)
		// todo log it
		return
	}
	clientConnection, _, err := hijacker.Hijack()
	if err != nil{
		http.Error(rw, "Internal server error", http.StatusInternalServerError)
		//todo log err
		return
	}
	// todo check clientConnection deadlines
	serverConnection, err := net.Dial("tcp", ws.TargetAddr)
	if err != nil{
		http.Error(rw, "Internal server error", http.StatusInternalServerError)
		//todo log err
		return
	}
	err = r.Write(serverConnection)
	if err != nil{
		http.Error(rw, "Internal server error", http.StatusInternalServerError)
		//todo log err
		return
	}
	go connectionRelay(clientConnection, serverConnection)
}

func connectionRelay(conn1 net.Conn, conn2 net.Conn){
	defer conn1.Close()
	defer conn2.Close()
	errChan := make(chan error, 2)
	go func(){
		_, err := io.Copy(conn1, conn2)
		if err != nil {
			errChan <- err
		}
	}()
	go func(){
		_, err := io.Copy(conn2, conn1)
		if err != nil {
			errChan <- err
		}
	}()
	<-errChan // todo log error
}

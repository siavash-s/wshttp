package wshttp

import (
	"net/http"
	"strings"
)

type WsHandler interface {
	HandleWs(http.ResponseWriter, *http.Request)
}

type HttpHandler interface {
	HandleHttp(http.ResponseWriter, *http.Request)
}

func Dmux(httpHandler HttpHandler, wsHandler WsHandler) func(http.ResponseWriter, *http.Request){
	return func(rw http.ResponseWriter, r *http.Request){
		if isWs(r){
			wsHandler.HandleWs(rw, r)
		}else {
			httpHandler.HandleHttp(rw, r)
		}
	}
}

func matchHeaders(headers http.Header, keyValueMatch map[string]string) bool{
	for matchKey, matchValue := range keyValueMatch{
		if headerValues, ok := headers[matchKey]; ok{
			if len(headerValues) < 1 || strings.ToLower(headerValues[0]) != strings.ToLower(matchValue){
				return false
			}
		}else {
			return false
		}
	}
	return true
}

var matches = map[string]string{"Connection": "upgrade", "Upgrade": "websocket"}
func isWs(r *http.Request) bool{
	if matchHeaders(r.Header, matches){return true}
	return false
}

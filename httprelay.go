package wshttp

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

type HttpRelay struct {
	TargetUrl string
}

func (httpRelay HttpRelay)HandleHttp(rw http.ResponseWriter, r *http.Request){
	u, err := url.Parse(httpRelay.TargetUrl)
	if err != nil{
		http.Error(rw, "Internal server error", http.StatusInternalServerError)
		// todo log it
		return
	}
	httputil.NewSingleHostReverseProxy(u).ServeHTTP(rw, r)
}

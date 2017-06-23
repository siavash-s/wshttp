package wshttp

import (
	"testing"
	"net/http"
	"net/http/httptest"
	"net/url"
)

func fakeHandleHttp(rw http.ResponseWriter, r *http.Request){
	rw.WriteHeader(http.StatusOK)
}

func TestHttpRelay_HandleHttp(t *testing.T) {
	backEndServer := httptest.NewServer(http.HandlerFunc(fakeHandleHttp))
	defer backEndServer.Close()
	backEndUrl, _ := url.Parse(backEndServer.URL)
	httpRelay := HttpRelay{TargetUrl: backEndUrl.String()}
	relayServer := httptest.NewServer(http.HandlerFunc(httpRelay.HandleHttp))
	defer relayServer.Close()
	relayUrl, _ := url.Parse(relayServer.URL)
	req, _ := http.NewRequest(http.MethodGet, relayUrl.String(), nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil{
		t.Fatalf("HTTP protocol error: %v", err)
	}
	t.Logf("HTTP status code: %v", resp.StatusCode)
	if resp.StatusCode != http.StatusOK{
		t.Errorf("not 200 status code for url: %v", relayUrl.String())
	}
}

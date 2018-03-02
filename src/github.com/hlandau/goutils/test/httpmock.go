package test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
)

type httpMockResponse struct {
	http.Response
	Data []byte
}

func (r *httpMockResponse) ServeHTTP(http.ResponseWriter, *http.Request) {
}

type HTTPMockTransport struct {
	mux *http.ServeMux
}

func (t *HTTPMockTransport) Clear() {
	t.mux = http.NewServeMux()
}

func (t *HTTPMockTransport) Add(path string, res *http.Response, data []byte) {
	if t.mux == nil {
		t.mux = http.NewServeMux()
	}

	t.mux.Handle(path, &httpMockResponse{
		Response: *res,
		Data:     data,
	})
}

func (t *HTTPMockTransport) AddHandlerFunc(path string, hf http.HandlerFunc) {
	if t.mux == nil {
		t.mux = http.NewServeMux()
	}

	t.mux.HandleFunc(path, hf)
}

func (t *HTTPMockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	h, ptn := t.mux.Handler(req)
	var res *httpMockResponse
	if ptn == "" {
		res = &httpMockResponse{
			Response: http.Response{
				StatusCode: 404,
			},
		}
	} else {
		var ok bool
		res, ok = h.(*httpMockResponse)
		if !ok {
			rw := httptest.NewRecorder()
			h.ServeHTTP(rw, req)
			res := &http.Response{
				StatusCode: rw.Code,
				Header:     rw.HeaderMap,
				Body:       ioutil.NopCloser(bytes.NewReader(rw.Body.Bytes())),
			}
			return res, nil
		}
	}
	res.Response.Body = ioutil.NopCloser(bytes.NewReader(res.Data))
	return &res.Response, nil
}

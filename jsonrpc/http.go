package jsonrpc

import (
	"context"
	"fmt"
	"log"
	"net/http"
)

type Http struct {
	rpc  *Server
	http *http.Server
}

func NewHttp(port uint16, methods []Method) *Http {
	h := &Http{
		rpc:  NewServer(),
		http: &http.Server{Addr: fmt.Sprintf(":%d", port)},
	}
	h.http.Handler = h
	for _, method := range methods {
		err := h.rpc.RegisterMethod(method.Name, method.ParamNames, method.Handler)
		if err != nil {
			panic(err)
		}
	}
	return h
}

// Run starts to listen for HTTP requests
func (h *Http) Run() error {
	go func() {
		log.Println("Listening for RPC requests on", h.http.Addr)
		h.http.ListenAndServe()
		h.Shutdown()
	}()
	return nil
}

// Shutdown stops the HTTP server operation
func (h *Http) Shutdown() error {
	return h.http.Shutdown(context.Background())
}

// ServeHTTP processes an incoming HTTP request
func (h *Http) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		req.Close = true
		return
	}

	resp, err := h.rpc.HandleReader(req.Body)
	writer.Header().Set("Content-Type", "application/json")
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
	} else {
		writer.WriteHeader(http.StatusOK)
	}
	if resp != nil {
		writer.Write(resp)
	}
}

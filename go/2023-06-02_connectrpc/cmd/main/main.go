package main

import (
	"net/http"

	playground "github.com/ras0q/connect-web-playground"
	"github.com/ras0q/connect-web-playground/internal/handler"
	"github.com/ras0q/connect-web-playground/pkg/bufgen/api/proto/protoconnect"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func main() {
	mux := http.NewServeMux()
	rh := handler.NewReadyHandler()
	path, handler := protoconnect.NewReadyServiceHandler(rh)
	mux.Handle(path, handler)
	mux.Handle("/", http.FileServer(http.FS(playground.PublicFS)))

	if err := http.ListenAndServe(
		"localhost:8080",
		h2c.NewHandler(mux, &http2.Server{}),
	); err != nil {
		panic(err)
	}
}

package server

import (
	"net/http"

	"github.com/VyacheZhuk/FINAL13-14/pkg/api"
)

const (
	port = "7540"
)

func Run() error {

	api.Init()

	http.Handle("/", http.FileServer(http.Dir("web")))

	return http.ListenAndServe(":"+port, nil)
}

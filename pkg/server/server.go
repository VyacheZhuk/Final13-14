package server

import (
	"fmt"
	"net/http"
	"os"

	"github.com/VyacheZhuk/FINAL13-14/pkg/api"
)

func Run() error {

	api.Init()

	port, ok := os.LookupEnv("TODO_PORT")
	if !ok {
		port = "7540"
	}

	http.Handle("/", http.FileServer(http.Dir("web")))

	return http.ListenAndServe(fmt.Sprintf(":%v", port), nil)
}

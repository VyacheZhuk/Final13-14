package main

import (
	"log"

	"github.com/VyacheZhuk/FINAL13-14/pkg/db"
	"github.com/VyacheZhuk/FINAL13-14/pkg/server"
)

func main() {
	if err := db.Init("scheduler.db"); err != nil {
		log.Fatalf("Не удалось выполнить инициализацию базы данных: %v", err)
	}
	defer db.GetDB().Close()

	if err := server.Run(); err != nil {
		log.Fatalf("Сбой в работе сервера: %v", err)
	}
}

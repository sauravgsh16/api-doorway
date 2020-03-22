package main

import (
	"os"

	"github.com/sauravgsh16/api-doorway/server"
)

func main() {
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_USER", "postgres")
	os.Setenv("DB_PWD", "postgres")
	os.Setenv("DB_NAME", "gateway")
	os.Setenv("DB_TYPE", "postgres")

	server.Run()
}

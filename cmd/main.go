package main

import (
	"github.com/ernur-eskermes/crud-app/internal/repository"
	"github.com/ernur-eskermes/crud-app/internal/service"
	"github.com/ernur-eskermes/crud-app/internal/transport/rest"
	"github.com/ernur-eskermes/crud-app/pkg/database"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"time"
)

const schema string = `
	create table if not exists books (
		id serial not null unique,
		title varchar(255) not null unique,
		author varchar(255) not null,
		publish_date timestamp not null default now(),
		rating int not null
	);
`

func main() {
	// init db
	db, err := database.NewPostgresConnection(database.ConnectionInfo{
		Host:     "ninja-db",
		Port:     5432,
		Username: "postgres",
		DBName:   "postgres",
		SSLMode:  "disable",
		Password: "qwerty123",
	})
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.MustExec(schema)

	// init deps
	booksRepo := postgres.NewBooks(db)
	booksService := service.NewBooks(booksRepo)
	handler := rest.NewHandler(booksService)

	// init & run server
	srv := &http.Server{
		Addr:    ":8000",
		Handler: handler.InitRouter(),
	}

	log.Println("SERVER STARTED AT", time.Now().Format(time.RFC3339))

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

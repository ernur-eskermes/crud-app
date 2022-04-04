package main

import (
	"encoding/json"
	"errors"
	"github.com/ernur-eskermes/crud-app/domain"
	"github.com/ernur-eskermes/crud-app/repository"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
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
	db, err := sqlx.Connect("postgres", "host=ninja-db user=postgres dbname=postgres password=qwerty123 sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.MustExec(schema)

	http.HandleFunc("/books", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			books, err := repository.GetBooks(*db)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			b, err := json.Marshal(books)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.Write(b)
		case http.MethodPost:
			var book domain.Book
			reqBytes, err := ioutil.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			if err = json.Unmarshal(reqBytes, &book); err != nil {
				res, err := json.Marshal(map[string]string{"detail": err.Error()})
				if err != nil {
					log.Fatalf("Error happened in JSON marshal. Err: %s", err.Error())
				}
				w.WriteHeader(http.StatusBadRequest)
				w.Header().Set("Content-Type", "application/json")
				w.Write(res)
				return
			}

			if err := repository.CreateBook(*db, book); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusCreated)
		}
	})
	http.HandleFunc("/book", func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()
		id, err := strconv.Atoi(params.Get("id"))
		if err != nil {
			log.Fatalf("Unable to convert the string into int.  %v", err)
		}
		switch r.Method {
		case http.MethodGet:
			book, err := repository.GetBookById(*db, id)
			if err != nil {
				if errors.Is(err, repository.ErrBookNotFound) {
					w.WriteHeader(http.StatusNotFound)
				} else {
					w.WriteHeader(http.StatusInternalServerError)
				}
				return
			}
			b, err := json.Marshal(book)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.Write(b)
		case http.MethodPut:
			if _, err := repository.GetBookById(*db, id); err != nil {
				if errors.Is(err, repository.ErrBookNotFound) {
					w.WriteHeader(http.StatusNotFound)
				} else {
					w.WriteHeader(http.StatusInternalServerError)
				}
				return
			}
			var book domain.Book
			reqBytes, err := ioutil.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			if err = json.Unmarshal(reqBytes, &book); err != nil {
				res, err := json.Marshal(map[string]string{"detail": err.Error()})
				if err != nil {
					log.Fatalf("Error happened in JSON marshal. Err: %s", err.Error())
				}
				w.WriteHeader(http.StatusBadRequest)
				w.Header().Set("Content-Type", "application/json")
				w.Write(res)
			}
			book.Id = id
			if err = repository.UpdateBook(*db, book); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			return
		case http.MethodDelete:
			if _, err = repository.GetBookById(*db, id); err != nil {
				if errors.Is(err, repository.ErrBookNotFound) {
					w.WriteHeader(http.StatusNotFound)
				} else {
					w.WriteHeader(http.StatusInternalServerError)
				}
				return
			}
			if err = repository.DeleteBook(*db, id); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusNoContent)
		}
	})

	if err = http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal(err)
	}
}

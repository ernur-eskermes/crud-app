package rest

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/ernur-eskermes/crud-app/internal/domain"

	"github.com/gorilla/mux"
)

type Books interface {
	Create(ctx context.Context, book domain.Book) error
	GetByID(ctx context.Context, id int) (domain.Book, error)
	GetAll(ctx context.Context) ([]domain.Book, error)
	Delete(ctx context.Context, id int) error
	Update(ctx context.Context, inp domain.Book) error
}

type Handler struct {
	booksService Books
}

func NewHandler(books Books) *Handler {
	return &Handler{
		booksService: books,
	}
}

func (h *Handler) InitRouter() *mux.Router {
	r := mux.NewRouter()
	r.Use(loggingMiddleware)

	books := r.PathPrefix("/books").Subrouter()
	{
		books.HandleFunc("", h.createBook).Methods(http.MethodPost)
		books.HandleFunc("", h.getAllBooks).Methods(http.MethodGet)
		books.HandleFunc("/{id:[0-9]+}", h.getBookByID).Methods(http.MethodGet)
		books.HandleFunc("/{id:[0-9]+}", h.deleteBook).Methods(http.MethodDelete)
		books.HandleFunc("/{id:[0-9]+}", h.updateBook).Methods(http.MethodPut)
	}

	return r
}

func (h *Handler) getBookByID(w http.ResponseWriter, r *http.Request) {
	id, err := getIdFromRequest(r)
	if err != nil {
		newResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	book, err := h.booksService.GetByID(context.TODO(), id)
	if err != nil {
		if errors.Is(err, domain.ErrBookNotFound) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		log.Println("getBookByID() error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(book)
	if err != nil {
		log.Println("getBookByID() error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(response)
}

func (h *Handler) createBook(w http.ResponseWriter, r *http.Request) {
	reqBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		newResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	var book domain.Book
	if err = json.Unmarshal(reqBytes, &book); err != nil {
		newResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	err = h.booksService.Create(context.TODO(), book)
	if err != nil {
		log.Println("createBook() error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) deleteBook(w http.ResponseWriter, r *http.Request) {
	id, err := getIdFromRequest(r)
	if err != nil {
		log.Println("deleteBook() error:", err)
		newResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	err = h.booksService.Delete(context.TODO(), id)
	if err != nil {
		log.Println("deleteBook() error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) getAllBooks(w http.ResponseWriter, r *http.Request) {
	books, err := h.booksService.GetAll(context.TODO())
	if err != nil {
		log.Println("getAllBooks() error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(books)
	if err != nil {
		log.Println("getAllBooks() error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(response)
}

func (h *Handler) updateBook(w http.ResponseWriter, r *http.Request) {
	id, err := getIdFromRequest(r)
	if err != nil {
		newResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	reqBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		newResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	var inp domain.Book
	if err = json.Unmarshal(reqBytes, &inp); err != nil {
		newResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	inp.Id = id
	err = h.booksService.Update(context.TODO(), inp)
	if err != nil {
		log.Println("error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func getIdFromRequest(r *http.Request) (int, error) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		return 0, err
	}

	if id == 0 {
		return 0, errors.New("id can't be 0")
	}

	return int(id), nil
}

type response struct {
	Message string `json:"message"`
}

func newResponse(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	if message != "" {
		res, err := json.Marshal(response{message})
		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err.Error())
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(res)
	}
}

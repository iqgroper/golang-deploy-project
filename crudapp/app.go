package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type Note struct {
	ID        int
	Text      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type NoteStorage struct {
	NoteList []*Note
	LastID   int
}

func (store *NoteStorage) List(w http.ResponseWriter, r *http.Request) {
	body, err := json.Marshal(store.NoteList)
	if err != nil {
		log.Println("Marshall error")
		http.Error(w, "storage err", http.StatusInternalServerError)
	}
	_, errWrite := w.Write(body)
	if errWrite != nil {
		http.Error(w, "write err", http.StatusInternalServerError)
	}
}

func (store *NoteStorage) Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	requiredIDStr, ok := vars["id"]
	if !ok {
		http.Error(w, "no id in query", http.StatusBadRequest)
		return
	}

	requiredID, errID := strconv.Atoi(requiredIDStr)
	if errID != nil {
		log.Println("cast error")
		http.Error(w, "id err", http.StatusBadRequest)
	}

	for _, note := range store.NoteList {
		if note.ID == requiredID {
			body, err := json.Marshal(note)
			if err != nil {
				log.Println("Marshall error")
				http.Error(w, "storage err", http.StatusInternalServerError)
			}
			_, errWrite := w.Write(body)
			if errWrite != nil {
				http.Error(w, "write err", http.StatusInternalServerError)
			}
			return
		}
	}
	http.Error(w, "No such note", http.StatusBadRequest)
}

func (store *NoteStorage) Create(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()

	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("error reading body")
		http.Error(w, "request err", http.StatusBadRequest)
	}

	req := make(map[string]string)
	errJSON := json.Unmarshal(reqBody, &req)
	if errJSON != nil {
		http.Error(w, "err unmarshalling", http.StatusInternalServerError)
	}

	if req["text"] == "" {
		http.Error(w, "request body err", http.StatusBadRequest)
	}

	newNote := &Note{
		ID:        store.LastID,
		Text:      req["text"],
		CreatedAt: time.Now().Round(time.Minute),
		UpdatedAt: time.Now().Round(time.Minute),
	}
	store.LastID++

	respBody, err := json.Marshal(newNote)
	if err != nil {
		log.Println("Marshall error")
		http.Error(w, "response err", http.StatusInternalServerError)
	}
	_, errWrite := w.Write(respBody)
	if errWrite != nil {
		http.Error(w, "write err", http.StatusInternalServerError)
	}
}

func (store *NoteStorage) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	requiredID, errID := strconv.Atoi(vars["id"])
	if errID != nil {
		log.Println("cast error")
		http.Error(w, "id err", http.StatusBadRequest)
	}

	defer r.Body.Close()

	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("error reading body")
		http.Error(w, "request err", http.StatusBadRequest)
	}

	req := make(map[string]string)
	errJSON := json.Unmarshal(reqBody, &req)
	if errJSON != nil {
		http.Error(w, "err unmarshalling", http.StatusInternalServerError)
	}

	if req["text"] == "" {
		http.Error(w, "request body err", http.StatusBadRequest)
	}

	for _, note := range store.NoteList {
		if note.ID == requiredID {
			note.Text = req["text"]
			note.UpdatedAt = time.Now().Round(time.Minute)

			_, errWrite := w.Write([]byte("success"))
			if errWrite != nil {
				http.Error(w, "write err", http.StatusInternalServerError)
			}
			return
		}
	}
	http.Error(w, "No such note", http.StatusBadRequest)
}

func (store *NoteStorage) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	requiredID, errID := strconv.Atoi(vars["id"])
	if errID != nil {
		log.Println("cast error")
		http.Error(w, "id err", http.StatusBadRequest)
	}

	delIndex := -1
	for idx, note := range store.NoteList {
		if note.ID == requiredID {
			delIndex = idx
		}
	}

	if delIndex == -1 {
		http.Error(w, "No such note", http.StatusBadRequest)
	}

	if delIndex < len(store.NoteList)-1 {
		copy(store.NoteList[delIndex:], store.NoteList[delIndex+1:])
	}
	store.NoteList[len(store.NoteList)-1] = nil
	store.NoteList = store.NoteList[:len(store.NoteList)-1]

	_, errWrite := w.Write([]byte("success"))
	if errWrite != nil {
		http.Error(w, "write err", http.StatusInternalServerError)
	}
}

func main() {
	r := mux.NewRouter()

	store := &NoteStorage{
		NoteList: make([]*Note, 5),
	}

	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		_, errWrite := w.Write([]byte("Im alive!"))
		if errWrite != nil {
			http.Error(w, "write err", http.StatusInternalServerError)
		}
	})
	r.HandleFunc("/note/{id:[0-9]+}", store.Get).Methods("GET")
	r.HandleFunc("/note", store.Create).Methods("POST")
	r.HandleFunc("/note/{id:[0-9]+}", store.Update).Methods("PUT")
	r.HandleFunc("/note/{id:[0-9]+}", store.Delete).Methods("DELETE")
	r.HandleFunc("/note", store.List).Methods("GET").Queries("order_by", "")

	port := ":8080"
	fmt.Println("starting server at", port)
	log.Fatal(http.ListenAndServe(port, r))
}

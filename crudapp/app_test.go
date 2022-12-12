package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/gorilla/mux"
)

// go test -coverprofile=cover.out && go tool cover -html=cover.out -o cover.html

func TestCreate(t *testing.T) {
	store := NoteStorage{
		NoteList: []*Note{},
	}

	newNote := map[string]string{
		"text": "first task",
	}

	expectedNote := &Note{
		ID:        0,
		Text:      "first task",
		CreatedAt: time.Now().Round(time.Second),
		UpdatedAt: time.Now().Round(time.Second),
	}

	reqBody, err := json.Marshal(newNote)
	if err != nil {
		t.Errorf("Err marshalling: %v", err)
	}

	port := "8080"
	req := httptest.NewRequest("POST", fmt.Sprintf("localhost:%s/note", port), bytes.NewReader(reqBody))

	w := httptest.NewRecorder()

	store.Create(w, req)

	if w.Code != 200 {
		t.Errorf("Wrong code\nExpected: 200\nGot: %d", w.Code)
	}

	resp := w.Result()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Failed to read body: %v", err)
	}

	recievedNote := &Note{}
	errMarshal := json.Unmarshal(responseBody, recievedNote)
	if errMarshal != nil {
		t.Errorf("Err unmarshalling: %v", errMarshal)
	}

	if !reflect.DeepEqual(recievedNote, expectedNote) {
		t.Errorf("response isnt correct\nWanted: %v\nGot: %v", expectedNote, recievedNote)
	}
}

func TestGet(t *testing.T) {

	store := NoteStorage{
		NoteList: []*Note{
			{
				ID:        0,
				Text:      "first note",
				CreatedAt: time.Now().Add(-48 * time.Hour).Round(time.Hour),
				UpdatedAt: time.Now().Round(time.Hour),
			},
			{
				ID:        1,
				Text:      "second note",
				CreatedAt: time.Now().Add(-24 * time.Hour).Round(time.Hour),
				UpdatedAt: time.Now().Add(-4 * time.Hour).Round(time.Hour),
			},
		},
	}

	port := "8080"
	req := httptest.NewRequest("GET", fmt.Sprintf("localhost:%s/note/1", port), nil)
	vars := map[string]string{
		"id": "1",
	}
	req = mux.SetURLVars(req, vars)

	w := httptest.NewRecorder()

	store.Get(w, req)

	if w.Code != 200 {
		t.Errorf("Wrong code\nExpected: 200\nGot: %d", w.Code)
	}

	resp := w.Result()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Failed to read body: %v", err)
	}

	recievedNote := &Note{}
	errMarshal := json.Unmarshal(responseBody, recievedNote)
	if errMarshal != nil {
		t.Errorf("Err unmarshalling: %v", errMarshal)
	}

	expectedNote := store.NoteList[1]

	if !reflect.DeepEqual(recievedNote, expectedNote) {
		t.Errorf("response isnt correct\nWanted: %v\nGot: %v", expectedNote, recievedNote)
	}

}

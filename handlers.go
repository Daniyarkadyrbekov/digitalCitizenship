package main

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"time"

	"github.com/digitalCitizenship/lib/storage/mock"
)

func init() {
	rand.Seed(time.Now().Unix())
}

func stubHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("stubHandler"))
}

func newInteraction(m *mock.Mock) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m.AddInteraction(rand.Int63(), rand.Int63(), time.Now().Unix())
	}
}

func interactedWithInfected(m *mock.Mock) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m.InteractedWithInfected(rand.Int63())
	}
}

func getInfectedList(m *mock.Mock) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		list, err := m.GetInfectedList()
		if err != nil {
			w.WriteHeader(http.StatusInsufficientStorage)
			return
		}
		listByte, err := json.Marshal(list)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(listByte)
	}
}

func newInfetcted(m *mock.Mock) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m.AddInfected(rand.Int63())
	}
}

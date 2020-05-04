package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/digitalCitizenship/lib/storage/mock"
	"github.com/volatiletech/authboss"
)

func stubHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("stubHandler"))
}

type InteractionReq struct {
	UserIIN string `json:"IIN"`
}

func newInteraction(m *mock.Mock) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u := r.Context().Value(authboss.CTXKeyUser)
		user, ok := u.(authboss.User)
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("error getting user from context"))
		}
		b, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		var intReq InteractionReq
		err = json.Unmarshal(b, &intReq)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		if err := m.AddInteraction(user.GetPID(), intReq.UserIIN, time.Now().Unix()); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.Header().Set("content-type", "application/json")
	}
}

func interactedWithInfected(m *mock.Mock) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u := r.Context().Value(authboss.CTXKeyUser)
		user, ok := u.(authboss.User)
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("error getting user from context"))
		}
		interacted, err := m.InteractedWithInfected(user.GetPID())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		data, err := json.Marshal(interacted)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("error marshaling user"))
		}
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	}
}

type ResponseList struct {
	List []string `json:"list"`
}

func getInfectedList(m *mock.Mock) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		list, err := m.GetInfectedList()
		if err != nil {
			w.WriteHeader(http.StatusInsufficientStorage)
			return
		}
		rl := ResponseList{List: list}
		listByte, err := json.Marshal(rl)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(listByte)
	}
}

func newInfetcted(m *mock.Mock) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		var intReq InteractionReq
		err = json.Unmarshal(b, &intReq)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		if err := m.AddInfected(intReq.UserIIN); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	}
}

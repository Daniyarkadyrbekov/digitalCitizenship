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
	Mac string `json:"mac"`
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
		secondUser, err := m.GetUserByMac(intReq.Mac)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		if err := m.AddInteraction(user.GetPID(), secondUser.IIN, time.Now().Unix()); err != nil {
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

type NewInfectedReq struct {
	IIN string `json:"IIN"`
}

func newInfetcted(m *mock.Mock) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		var intReq NewInfectedReq
		err = json.Unmarshal(b, &intReq)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		if err := m.AddInfected(intReq.IIN); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	}
}

func getUsersList(m *mock.Mock) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		list, err := m.GetUsersList()
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

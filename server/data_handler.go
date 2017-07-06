package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	ds "bitbucket.org/enticusa/kingdb/docstore"
	"bitbucket.org/enticusa/kingdb/serialization"
)

func (s *httpServer) dataHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		s.handleGet(w, r)
	case "POST":
		s.handlePost(w, r)
	case "PUT":
		s.handlePut(w, r)
	case "DELETE":
		s.handlePut(w, r)
	default:
		s.handleUnkown(w, r)
	}
}

func (s *httpServer) handlePost(w http.ResponseWriter, r *http.Request) {

	pathParts := strings.Split(r.URL.Path, "/")
	label := ds.Label(pathParts[2])

	parentID := r.URL.Query().Get("parent")
	documentID := ds.ID(pathParts[3])

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotAcceptable)
		return
	}

	fields := make(ds.KeyValueMap)
	json.Unmarshal(body, &fields)

	document := ds.Document{
		ID:             documentID,
		ParentID:       ds.ID(parentID),
		KeyValueFields: fields,
	}

	if document.ID == "" {
		document.ID = s.store.GenerateID()
	}

	id, err := s.store.CreateDocument(label, document)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotAcceptable)
		return
	}

	w.Write([]byte(id))
}

func (s *httpServer) handleGet(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(r.URL.Path, "/")

	label := ds.Label(pathParts[2])
	documentID := ds.ID(pathParts[3])

	depth, err := strconv.Atoi(r.URL.Query().Get("depth"))
	if err != nil {
		errMsg := fmt.Sprintf("Invalid depth parameter: %s", err)
		http.Error(w, errMsg, http.StatusNotAcceptable)
		return
	}

	document, err := s.store.ReadDocument(label, documentID)
	if err != nil {
		errMsg := fmt.Sprintf("Error while reading document: %s", err)
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}

	json, err := serialization.SerializeDocument(&document, int(depth))
	if err != nil {
		errMsg := fmt.Sprintf("Error while searializing document: %s", err)
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(json)
}

func (s *httpServer) handlePut(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Updates not yet implemented", http.StatusNotImplemented)
}

func (s *httpServer) handleDelete(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Deletes not yet implemented", http.StatusNotImplemented)
}

func (s *httpServer) handleUnkown(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Unkwown Action", http.StatusBadRequest)
}

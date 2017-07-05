package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	ds "github.com/arturom/docdb/docstore"
	"github.com/arturom/docdb/serialization"
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

	parentID, err := strconv.Atoi(r.URL.Query().Get("parent"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotAcceptable)
		return
	}

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotAcceptable)
		return
	}

	fields := make(ds.KeyValueMap)
	json.Unmarshal(body, &fields)

	document := ds.Document{
		ParentID:       ds.ID(parentID),
		KeyValueFields: fields,
	}

	id, err := s.store.CreateDocument(label, document)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotAcceptable)
		return
	}

	w.Write([]byte(strconv.Itoa(int(id))))
}

func (s *httpServer) handleGet(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(r.URL.Path, "/")

	label := ds.Label(pathParts[2])

	id, err := strconv.Atoi(string(pathParts[3]))
	if err != nil {
		errMsg := fmt.Sprintf("Invalid document ID: %s", err)
		http.Error(w, errMsg, http.StatusNotAcceptable)
		return
	}

	depth, err := strconv.Atoi(r.URL.Query().Get("depth"))
	if err != nil {
		errMsg := fmt.Sprintf("Invalid depth parameter: %s", err)
		http.Error(w, errMsg, http.StatusNotAcceptable)
		return
	}

	documentID := ds.ID(id)

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

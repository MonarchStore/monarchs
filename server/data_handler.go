package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	ds "bitbucket.org/enticusa/kingdb/docstore"
	"bitbucket.org/enticusa/kingdb/serialization"
)

func parsePath(p string) parsedPath {
	parts := strings.Split(p, "/")
	count := len(parts)

	parsed := parsedPath{
		count: count,
	}
	if count > 0 {
		parsed.storeName = parts[0]
	}
	if count > 1 {
		parsed.documentType = ds.Label(parts[1])
	}
	if count > 2 {
		parsed.documentID = ds.ID(parts[2])
	}

	return parsed
}

type parsedPath struct {
	count        int
	storeName    string
	documentID   ds.ID
	documentType ds.Label
}

func (s *httpServer) dataHandler(w http.ResponseWriter, r *http.Request) {
	path := parsePath(strings.Trim(r.URL.Path, "/"))
	signature := fmt.Sprintf("%s%d", r.Method, path.count)

	switch signature {
	case "GET3":
		s.readDocument(w, r, path)
	case "POST3":
		s.createDocument(w, r, path)
	case "PUT3":
		s.updateDocument(w, r, path)
	case "DELETE3":
		s.deleteDocument(w, r, path)
	case "GET1":
		s.readSchema(w, r, path)
	case "POST1":
		s.createSchema(w, r, path)
	case "DELETE1":
		s.deleteSchema(w, r, path)
	default:
		s.handleUnkown(w, r)
	}
}

func (s *httpServer) createDocument(w http.ResponseWriter, r *http.Request, path parsedPath) {
	store, ok := s.storeMap[path.storeName]
	if !ok {
		http.NotFound(w, r)
		return
	}

	parentID := r.URL.Query().Get("parent")

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fields := make(ds.KeyValueMap)
	err = json.Unmarshal(body, &fields)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotAcceptable)
		return
	}

	document := ds.Document{
		ID:             path.documentID,
		ParentID:       ds.ID(parentID),
		KeyValueFields: fields,
	}

	if document.ID == "" {
		document.ID = store.GenerateID()
	}

	id, err := store.CreateDocument(path.documentType, document)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotAcceptable)
		return
	}

	w.Write([]byte(id))
}

func (s *httpServer) readDocument(w http.ResponseWriter, r *http.Request, path parsedPath) {
	store, ok := s.storeMap[path.storeName]
	if !ok {
		http.NotFound(w, r)
		return
	}

	depth, err := strconv.Atoi(r.URL.Query().Get("depth"))
	if err != nil {
		errMsg := fmt.Sprintf("Invalid depth parameter: %s", err)
		http.Error(w, errMsg, http.StatusNotAcceptable)
		return
	}

	parentCountStr := r.URL.Query().Get("parents")
	parentCount := 0
	if parentCountStr != "" {
		parentCount, err = strconv.Atoi(parentCountStr)
		if err != nil {
			errMsg := fmt.Sprintf("Invalid parent parameter: %s", err)
			http.Error(w, errMsg, http.StatusNotAcceptable)
			return
		}
	}

	document, err := store.ReadDocument(path.documentType, path.documentID)
	if err != nil {
		errMsg := fmt.Sprintf("Error while reading document: %s", err)
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}

	parents, err := store.ReadParentDocuments(path.documentType, path.documentID, parentCount)

	json, err := serialization.SerializeDocument(&document, int(depth), parents)
	if err != nil {
		errMsg := fmt.Sprintf("Error while searializing document: %s", err)
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(json)
}

func (s *httpServer) updateDocument(w http.ResponseWriter, r *http.Request, path parsedPath) {
	store, ok := s.storeMap[path.storeName]
	if !ok {
		http.NotFound(w, r)
		return
	}
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotAcceptable)
		return
	}

	fields := make(ds.KeyValueMap)
	err = json.Unmarshal(body, &fields)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotAcceptable)
		return
	}

	document := ds.Document{
		ID:             path.documentID,
		KeyValueFields: fields,
	}

	err = store.UpdateDocument(path.documentType, document)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotAcceptable)
		return
	}

	w.Write([]byte("OK"))
}

func (s *httpServer) deleteDocument(w http.ResponseWriter, r *http.Request, path parsedPath) {
	store, ok := s.storeMap[path.storeName]
	if !ok {
		http.NotFound(w, r)
		return
	}
	err := store.DeleteDocument(path.documentType, path.documentID)
	if err != nil {
		errMsg := fmt.Sprintf("Error while deleting document: %s", err)
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}

	w.Write([]byte("OK"))
}

func (s *httpServer) readSchema(w http.ResponseWriter, r *http.Request, path parsedPath) {
	store, ok := s.storeMap[path.storeName]
	if !ok {
		http.NotFound(w, r)
		return
	}
	labels := store.GetHierarchyLabels()

	json, err := json.Marshal(labels)
	if err != nil {
		errMsg := fmt.Sprintf("Error while searializing schema: %s", err)
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func (s *httpServer) createSchema(w http.ResponseWriter, r *http.Request, path parsedPath) {
	_, ok := s.storeMap[path.storeName]
	if ok {
		errMsg := fmt.Sprintf("Schema already exists: %s", path.storeName)
		http.Error(w, errMsg, http.StatusNotAcceptable)
		return
	}

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	labelCount := bytes.Count(body, []byte(",")) + 1

	labels := make(ds.Labels, labelCount)
	err = json.Unmarshal(body, &labels)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotAcceptable)
		return
	}

	s.storeMap[path.storeName] = ds.NewStore(labels)
	w.Write([]byte("OK"))
}

func (s *httpServer) deleteSchema(w http.ResponseWriter, r *http.Request, path parsedPath) {
	_, ok := s.storeMap[path.storeName]
	if !ok {
		http.NotFound(w, r)
		return
	}
	delete(s.storeMap, path.storeName)
	w.Write([]byte("OK"))
}

func (s *httpServer) handleUnkown(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}

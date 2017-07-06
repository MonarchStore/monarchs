package serialization

import (
	"encoding/json"
	"sort"

	ds "bitbucket.org/enticusa/kingdb/docstore"
)

type ReadModel struct {
	ID       ds.ID          `json:"id"`
	ParentID ds.ID          `json:"parent_id"`
	Values   ds.KeyValueMap `json:"values"`
	Children ReadModelSlice `json:"children"`
}

type ReadModelSlice []ReadModel

func (s ReadModelSlice) Len() int {
	return len(s)
}

func (s ReadModelSlice) Less(i, j int) bool {
	return s[i].ID < s[j].ID
}

func (s ReadModelSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func mapDocumentToReadModel(d *ds.Document, depth int) ReadModel {
	var children ReadModelSlice

	if depth > 0 {
		childrenCount := len(d.NestedDocuments)
		children = make(ReadModelSlice, childrenCount)
		childIndex := 0

		for _, c := range d.NestedDocuments {
			children[childIndex] = mapDocumentToReadModel(c, depth-1)
			childIndex++
		}

		sort.Sort(children)
	}

	return ReadModel{
		ID:       d.ID,
		ParentID: d.ParentID,
		Values:   d.KeyValueFields,
		Children: children,
	}
}

func SerializeDocument(d *ds.Document, depth int) ([]byte, error) {
	m := mapDocumentToReadModel(d, depth)
	return json.Marshal(m)
}

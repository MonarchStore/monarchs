package serialization

import (
	"encoding/json"
	"sort"

	ds "github.com/arturom/monarchs/docstore"
)

type ReadModel struct {
	ID       ds.ID          `json:"id"`
	Values   ds.KeyValueMap `json:"values"`
	Children ReadModelSlice `json:"children"`
	Parents  []ReadModel    `json:"parents,omitempty"`
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
		Values:   d.KeyValueFields,
		Children: children,
	}
}

func SerializeDocument(d *ds.Document, depth int, parents ds.DocumentSlice) ([]byte, error) {
	m := mapDocumentToReadModel(d, depth)
	m.Parents = make([]ReadModel, len(parents))
	i := 0

	for _, d := range parents {
		if d == nil {
			break
		}
		m.Parents[i] = mapDocumentToReadModel(d, 0)
		i++
	}

	return json.Marshal(m)
}

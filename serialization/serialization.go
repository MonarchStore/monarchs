package serialization

import ds "github.com/arturom/docdb/docstore"
import "encoding/json"

type ReadModelSlice []ReadModel

type ReadModel struct {
	ID       ds.ID          `json:"id"`
	ParentID ds.ID          `json:"parent_id"`
	Values   ds.KeyValueMap `json:"values"`
	Children ReadModelSlice `json:"children"`
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

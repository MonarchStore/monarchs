package docstore

type ID string

type KeyValueMap map[string]interface{}

type DocumentMap map[ID]*Document

type DocumentSlice []*Document

type Document struct {
	ID              ID
	ParentID        ID
	KeyValueFields  KeyValueMap
	NestedDocuments DocumentMap
}

package docstore

type ID string

type KeyValueMap map[string]string

type DocumentMap map[ID]*Document

type Document struct {
	ID              ID
	ParentID        ID
	KeyValueFields  KeyValueMap
	NestedDocuments DocumentMap
}

package docstore

import "fmt"
import "strconv"

type Store interface {
	GenerateID() ID
	GetHierarchyLabels() Labels
	CreateDocument(documentType Label, doc Document) (ID, error)
	ReadDocument(documentType Label, id ID) (Document, error)
	UpdateDocument(documentType Label, document Document) error
	DeleteDocument(documentType Label, id ID) error
}

type StoreMap map[string]Store

type IDGenerator interface {
	getAndIncrement() ID
}

func NewIDGenerator() IDGenerator {
	return &idGenerator{
		nextID: 0,
	}
}

type idGenerator struct {
	nextID int
}

func (i *idGenerator) getAndIncrement() ID {
	current := i.nextID
	i.nextID++
	return ID(strconv.Itoa(current))
}

func NewStore(labels Labels) Store {
	idGenerator := NewIDGenerator()
	linkMap, linkedListHead := buildHierarchyLinkMap(labels, idGenerator)
	return store{
		idGenerator:    idGenerator,
		linkMap:        linkMap,
		linkedListHead: linkedListHead,
	}
}

type store struct {
	idGenerator    IDGenerator
	linkMap        HierarchyLinkyMap
	linkedListHead *HierarchyLink
}

func buildHierarchyLinkMap(labels Labels, idGenerator IDGenerator) (HierarchyLinkyMap, *HierarchyLink) {
	linkedListHead := &HierarchyLink{
		Label:       "root",
		DocumentMap: make(DocumentMap),
		ParentLink:  nil,
		ChildLink:   nil,
	}
	rootDocument := Document{
		ID:              "root",
		NestedDocuments: make(DocumentMap),
		KeyValueFields:  make(KeyValueMap),
		ParentID:        "root_parent",
	}
	linkedListHead.DocumentMap[rootDocument.ID] = &rootDocument

	lastInserted := linkedListHead
	linkMap := make(HierarchyLinkyMap)

	linkMap[linkedListHead.Label] = linkedListHead

	for _, label := range labels {
		currentLink := &HierarchyLink{
			Label:       label,
			DocumentMap: make(DocumentMap),
			ParentLink:  lastInserted,
			ChildLink:   nil,
		}
		lastInserted.ChildLink = currentLink
		linkMap[label] = currentLink
		lastInserted = currentLink
	}

	return linkMap, linkedListHead
}

func (s store) GenerateID() ID {
	return s.idGenerator.getAndIncrement()
}

func (s store) GetHierarchyLabels() Labels {
	labels := make(Labels, len(s.linkMap))
	current := s.linkedListHead
	i := 0
	for {
		labels[i] = current.Label

		if current.ChildLink == nil {
			break
		}
		i++
		current = current.ChildLink
	}

	return labels
}

func (s store) CreateDocument(documentType Label, doc Document) (ID, error) {
	link, ok := s.linkMap[documentType]
	if !ok {
		return doc.ID, fmt.Errorf("Cannot create document. Invalid document type: %s", documentType)
	}

	parentDocument, ok := link.ParentLink.DocumentMap[doc.ParentID]
	if !ok {
		return doc.ID, fmt.Errorf("Cannot create document. Parent ID not found: %s:%s", link.ParentLink.Label, doc.ParentID)
	}

	_, exists := link.DocumentMap[doc.ID]
	if exists {
		return doc.ID, fmt.Errorf("Cannot create document. Document ID already exists: %s:%s", link.Label, doc.ID)
	}

	doc.NestedDocuments = make(DocumentMap)

	if doc.KeyValueFields == nil {
		doc.KeyValueFields = make(KeyValueMap)
	}

	parentDocument.NestedDocuments[doc.ID] = &doc
	link.DocumentMap[doc.ID] = &doc

	return doc.ID, nil
}

func (s store) ReadDocument(documentType Label, id ID) (Document, error) {
	link, ok := s.linkMap[documentType]
	if !ok {
		return Document{}, fmt.Errorf("Cannot retrieve document. Invalid document type: %s", documentType)
	}
	document, ok := link.DocumentMap[id]
	if !ok {
		return Document{}, fmt.Errorf("Cannot retrieve document. Document ID not found: %s:%s", documentType, id)
	}
	return *document, nil
}

func (s store) UpdateDocument(documentType Label, document Document) error {
	link, ok := s.linkMap[documentType]
	if !ok {
		return fmt.Errorf("Cannot update document. Invalid document type: %s", documentType)
	}
	existingDocument, ok := link.DocumentMap[document.ID]
	if !ok {
		return fmt.Errorf("Cannot update document. Document ID not found: %s:%s", documentType, document.ID)
	}
	existingDocument.KeyValueFields = document.KeyValueFields
	return nil
}

func (s store) DeleteDocument(documentType Label, id ID) error {
	link, ok := s.linkMap[documentType]
	if !ok {
		return fmt.Errorf("Cannot delete document. Invalid document type: %s", documentType)
	}

	document, ok := link.DocumentMap[id]
	if !ok {
		return fmt.Errorf("Cannot delete document. Document ID not found: %s:%s", documentType, document.ID)
	}

	parentDocument, ok := link.ParentLink.DocumentMap[document.ParentID]
	if !ok {
		return fmt.Errorf("Cannot delete document. Parent ID not found: %s:%s", link.ParentLink.Label, document.ParentID)
	}

	delete(parentDocument.NestedDocuments, document.ID)
	delete(link.ParentLink.DocumentMap, document.ID)

	return nil
}

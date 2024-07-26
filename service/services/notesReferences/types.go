package notesReferencesService

import (
	"github.com/alperklc/the-zula/service/infrastructure/db/notes"
	"github.com/alperklc/the-zula/service/infrastructure/db/notesReferences"
)

type NoteReferenceNode struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

type NoteReferenceLink struct {
	Source string `json:"source"`
	Target string `json:"target"`
}

type NoteReferencesResponse struct {
	Links []NoteReferenceLink `json:"links"`
	Nodes []NoteReferenceNode `json:"nodes"`
}

func NewNoteReferencesResponse(l []notesReferences.ReferencesDocument, n []notes.NoteDocument) NoteReferencesResponse {

	var links []NoteReferenceLink
	for _, link := range l {
		links = append(links, NoteReferenceLink{
			Source: link.From,
			Target: link.To,
		})
	}

	var nodes []NoteReferenceNode
	for _, node := range n {
		nodes = append(nodes, NoteReferenceNode{
			Title: node.Title,
			ID:    node.Id,
		})
	}

	return NoteReferencesResponse{
		Nodes: nodes,
		Links: links,
	}
}

package referencesService

import (
	"github.com/alperklc/the-zula/service/infrastructure/db/notes"
	"github.com/alperklc/the-zula/service/infrastructure/db/references"
)

type ReferenceNode struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

type ReferenceLink struct {
	Source string `json:"source"`
	Target string `json:"target"`
}

type ReferencesResponse struct {
	Links []ReferenceLink `json:"links"`
	Nodes []ReferenceNode `json:"nodes"`
}

func NewReferencesResponse(l []references.ReferencesDocument, n []notes.NoteDocument) ReferencesResponse {
	var links []ReferenceLink
	for _, link := range l {
		links = append(links, ReferenceLink{
			Source: link.From,
			Target: link.To,
		})
	}

	var nodes []ReferenceNode
	for _, node := range n {
		nodes = append(nodes, ReferenceNode{
			Title: node.Title,
			ID:    node.ShortId,
		})
	}

	return ReferencesResponse{
		Nodes: nodes,
		Links: links,
	}
}

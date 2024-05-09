package notesReferencesCtrl

import (
	"github.com/alperklc/the-zula/service/infrastructure/db/notesReferences"
)

func GetNoteIdsFromReferences(references []notesReferences.ReferencesDocument) []string {

	nodesMap := make(map[string]bool)
	nodeIds := []string{}
	for _, ref := range references {
		if !nodesMap[ref.From] {
			nodeIds = append(nodeIds, ref.From)
			nodesMap[ref.From] = true
		}

		if !nodesMap[ref.To] {
			nodeIds = append(nodeIds, ref.To)
			nodesMap[ref.To] = true
		}
	}

	return nodeIds
}

package referencesService

import (
	"fmt"
	"math"

	"github.com/alperklc/the-zula/service/utils"

	"github.com/alperklc/the-zula/service/infrastructure/db/notes"
	"github.com/alperklc/the-zula/service/infrastructure/db/references"
)

type ReferencesService interface {
	ListReferencesToNote(userId, noteId string, depth int) (ReferencesResponse, error)
	UpsertReferencesOfNote(noteId, noteContent string) error
	DeleteReferencesOfNote(noteId string) error
}

type datasources struct {
	notes      notes.Collection
	references references.Collection
}

func NewService(n notes.Collection, nr references.Collection) ReferencesService {
	return &datasources{
		notes: n, references: nr,
	}
}

func getPaginationRange(count, page, pageSize int) (int, int, string) {
	if count == 0 {
		return 0, 0, ""
	}

	var from int = (page-1)*pageSize + 1
	var to int = int(math.Min(float64(page*pageSize), float64(count)))

	return from, to, fmt.Sprintf("%d - %d", from, to)
}

func (d *datasources) ListReferencesToNote(userId, noteId string, depth int) (ReferencesResponse, error) {
	note, getNoteErr := d.notes.GetOne(noteId)
	if getNoteErr != nil {
		return ReferencesResponse{}, getNoteErr
	}

	if userId != note.CreatedBy {
		return ReferencesResponse{}, fmt.Errorf("NOT_ALLOWED_TO_GET")
	}

	references, listErr := d.references.ListReferencesOfNoteInDepth(noteId, depth)
	nodeIds := GetNoteIdsFromReferences(references)

	notes, getNotesErr := d.notes.GetNotes(nodeIds, []string{"id", "shortId", "title"})
	if getNotesErr != nil {
		return ReferencesResponse{}, getNotesErr
	}

	return NewReferencesResponse(references, notes), listErr
}

func (d *datasources) UpsertReferencesOfNote(noteId, noteContent string) error {
	idsOfReferences := utils.ParseInternalLinksFromNote(noteContent)

	errDelete := d.references.DeleteAllReferencesFromNote(noteId)
	if errDelete != nil {
		return errDelete
	}

	// check if each note exist
	validReferences, errGetNotes := d.notes.GetNotes(idsOfReferences, []string{"id", "shortId"})
	if errGetNotes != nil {
		return errGetNotes
	}
	var idsOfValidReferences []string
	for _, r := range validReferences {
		idsOfValidReferences = append(idsOfValidReferences, r.ShortId)
	}

	errInsert := d.references.InsertMany(noteId, idsOfValidReferences)
	if errInsert != nil {
		return errInsert
	}

	return nil
}

func (d *datasources) DeleteReferencesOfNote(noteId string) error {
	errDeleteOne := d.references.DeleteAllReferencesFromNote(noteId)
	if errDeleteOne != nil {
		return errDeleteOne
	}

	errDeleteTwo := d.references.DeleteAllReferencesToNote(noteId)
	if errDeleteTwo != nil {
		return errDeleteTwo
	}

	return nil
}

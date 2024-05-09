package notesReferencesCtrl

import (
	"fmt"
	"math"

	"github.com/alperklc/the-zula/service/utils"

	"github.com/alperklc/the-zula/service/infrastructure/db/notes"
	"github.com/alperklc/the-zula/service/infrastructure/db/notesReferences"
)

type NotesReferencesController interface {
	ListReferencesToNote(userId, noteId string, depth int) (NoteReferencesResponse, error)
	UpsertReferencesOfNote(noteId, noteContent string) error
	DeleteReferencesOfNote(noteId string) error
}

type datasources struct {
	notes           notes.Collection
	notesReferences notesReferences.Collection
}

func NewNotesReferencesController(n notes.Collection, nr notesReferences.Collection) NotesReferencesController {
	return &datasources{
		notes: n, notesReferences: nr,
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

func (d *datasources) ListReferencesToNote(userId, noteId string, depth int) (NoteReferencesResponse, error) {
	note, getNoteErr := d.notes.GetOne(noteId)
	if getNoteErr != nil {
		return NoteReferencesResponse{}, getNoteErr
	}

	if userId != note.CreatedBy {
		return NoteReferencesResponse{}, fmt.Errorf("NOT_ALLOWED_TO_GET")
	}

	references, listErr := d.notesReferences.ListReferencesOfNoteInDepth(noteId, depth)
	nodeIds := GetNoteIdsFromReferences(references)

	notes, getNotesErr := d.notes.GetNotes(nodeIds, []string{"id", "title"})
	if getNotesErr != nil {
		return NoteReferencesResponse{}, getNotesErr
	}

	return NewNoteReferencesResponse(references, notes), listErr
}

func (d *datasources) UpsertReferencesOfNote(noteId, noteContent string) error {
	idsOfReferences := utils.ParseInternalLinksFromNote(noteContent)

	errDelete := d.notesReferences.DeleteAllReferencesFromNote(noteId)
	if errDelete != nil {
		return errDelete
	}

	// check if each note exist
	validReferences, errGetNotes := d.notes.GetNotes(idsOfReferences, []string{"id"})
	if errGetNotes != nil {
		return errGetNotes
	}
	var idsOfValidReferences []string
	for _, r := range validReferences {
		idsOfValidReferences = append(idsOfValidReferences, r.Id)
	}

	errInsert := d.notesReferences.InsertMany(noteId, idsOfValidReferences)
	if errInsert != nil {
		return errInsert
	}

	return nil
}

func (d *datasources) DeleteReferencesOfNote(noteId string) error {
	errDeleteOne := d.notesReferences.DeleteAllReferencesFromNote(noteId)
	if errDeleteOne != nil {
		return errDeleteOne
	}

	errDeleteTwo := d.notesReferences.DeleteAllReferencesToNote(noteId)
	if errDeleteTwo != nil {
		return errDeleteTwo
	}

	return nil
}

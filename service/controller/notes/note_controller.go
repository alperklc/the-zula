package notectrl

import (
	"fmt"
	"math"

	"github.com/alperklc/the-zula/service/infrastructure/db/notes"
	"github.com/alperklc/the-zula/service/infrastructure/db/notesDrafts"
	"github.com/alperklc/the-zula/service/utils"
)

type NoteController interface {
	SearchTags(userId, searchKeyword string, limit int) ([]notes.TagsResult, error)
	GetStatistics(userId string) (Statistics, error)
	ListNotes(userId string, searchKeyword *string, page, pageSize *int, sortBy, sortDirection *string, tags *[]string) (NotesPage, error)
	CreateNote(userId, clientId string, title, content *string, tags *[]string) (Note, error)
	UpdateNote(noteId, userId, clientId string, update map[string]interface{}) error
	GetNote(noteId, userId, clientId string) (Note, error)
	GetNotes(userId string, noteIds, fields []string) ([]Note, error)
	DeleteNote(noteId, userId, clientId string) error
	GetDraftOfNote(userId, noteId string) (Note, error)
	UpdateDraft(userId, noteId, title, content string, tags []string) error
	DeleteDraft(userId, noteId string) error
}

type datasources struct {
	notes       notes.Collection
	notesDrafts notesDrafts.Collection
}

func NewNotesController(n notes.Collection, nd notesDrafts.Collection) NoteController {
	return &datasources{
		notes: n, notesDrafts: nd,
	}
}

func getPaginationRange(count, page, pageSize int) string {
	if count == 0 {
		return ""
	}

	var from int = (page-1)*pageSize + 1
	var to int = int(math.Min(float64(page*pageSize), float64(count)))

	return fmt.Sprintf("%d - %d", from, to)
}

func (d *datasources) SearchTags(userId, searchKeyword string, limit int) ([]notes.TagsResult, error) {
	return d.notes.SearchTags(userId, searchKeyword, limit)
}

func (d *datasources) GetStatistics(userId string) (Statistics, error) {
	count, err := d.notes.Count(userId)
	stats := Statistics{Count: count}

	return stats, err
}

func (d *datasources) ListNotes(userId string, q *string, p, ps *int, sb, sd *string, t *[]string) (NotesPage, error) {
	searchKeyword := ""
	if q != nil {
		searchKeyword = *q
	}

	page := 1
	if p != nil {
		page = *p
	}

	pageSize := 10
	if ps != nil {
		pageSize = *ps
	}

	sortBy := "updatedAt"
	if sb != nil {
		sortBy = *sb
	}

	sortDirection := "desc"
	if sd != nil {
		sortDirection = *sd
	}

	tags := []string{}
	if t != nil {
		tags = *t
	}

	notes, count, listErr := d.notes.List(userId, searchKeyword, page, pageSize, sortBy, sortDirection, tags)

	// query drafts of note entries on the page
	idsOfNotes := make([]string, len(notes))
	for i, v := range notes {
		idsOfNotes[i] = v.Id
	}
	draftsOnPage, _ := d.notesDrafts.CheckExistence(idsOfNotes)

	var items []Note = make([]Note, 0, len(notes))
	for _, b := range notes {
		note := Note{b.Id, b.UpdatedAt, b.UpdatedBy, b.CreatedBy, b.CreatedAt, b.Title, "", b.Tags, draftsOnPage[b.Id]}
		items = append(items, note)
	}

	return NotesPage{
		Items: items,
		Meta: PaginationMeta{
			Count:         count,
			Query:         searchKeyword,
			Page:          page,
			PageSize:      pageSize,
			SortBy:        sortBy,
			SortDirection: sortDirection,
			Range:         getPaginationRange(count, page, pageSize),
		},
	}, listErr
}

func (d *datasources) CreateNote(userId, clientId string, title, content *string, t *[]string) (Note, error) {
	tags := []string{}
	if t != nil {
		tags = *t
	}

	createdNote, err := d.notes.InsertOne(userId, *title, *content, tags)

	return Note{
		ShortId:   createdNote.ShortId,
		UpdatedAt: createdNote.UpdatedAt,
		UpdatedBy: createdNote.UpdatedBy,
		CreatedBy: createdNote.CreatedBy,
		CreatedAt: createdNote.CreatedAt,
		Title:     createdNote.Title,
		Content:   createdNote.Content,
		Tags:      createdNote.Tags,
	}, err
}

func (d *datasources) UpdateNote(noteId, userId, clientId string, update map[string]interface{}) error {
	note, getNoteErr := d.notes.GetOne(noteId)
	if getNoteErr != nil {
		return getNoteErr
	}

	if note.CreatedBy != userId {
		return fmt.Errorf("NOT_ALLOWED_TO_UPDATE")
	}

	// only following fields are allowed for update
	allowedFields := []string{"title", "content", "tags"}
	updatedFields := utils.FilterFieldsOfObject(allowedFields, update)
	err := d.notes.UpdateOne(userId, noteId, updatedFields)
	if err == nil {
		d.notesDrafts.DeleteOne(noteId)
	}

	return err
}

func (d *datasources) GetNote(noteId, userId, clientId string) (Note, error) {
	note, getNoteErr := d.notes.GetOne(noteId)
	if getNoteErr != nil || note.Id == "" {
		return Note{}, getNoteErr
	}

	if userId != note.CreatedBy {
		return Note{}, fmt.Errorf("NOT_ALLOWED_TO_GET")
	}

	draftExist, _ := d.notesDrafts.CheckExistence([]string{noteId})

	return Note{
		ShortId:   note.ShortId,
		UpdatedAt: note.UpdatedAt,
		UpdatedBy: note.UpdatedBy,
		CreatedBy: note.CreatedBy,
		CreatedAt: note.CreatedAt,
		Title:     note.Title,
		Content:   note.Content,
		Tags:      note.Tags,
		HasDraft:  draftExist[note.ShortId],
	}, getNoteErr
}

func (d *datasources) GetNotes(userId string, noteIds, fields []string) ([]Note, error) {
	notesFound, err := d.notes.GetNotes(noteIds, fields)
	var items []Note = make([]Note, 0, len(notesFound))

	draftExist, _ := d.notesDrafts.CheckExistence(noteIds)

	for _, b := range notesFound {
		if b.CreatedBy == userId {
			note := Note{
				ShortId:   b.ShortId,
				UpdatedAt: b.UpdatedAt,
				UpdatedBy: b.UpdatedBy,
				CreatedBy: b.CreatedBy,
				CreatedAt: b.CreatedAt,
				Title:     b.Title,
				Content:   b.Content,
				Tags:      b.Tags,
				HasDraft:  draftExist[b.ShortId],
			}

			items = append(items, note)
		}
	}

	return items, err
}

func (d *datasources) DeleteNote(noteId, userId, clientId string) error {
	note, getNoteErr := d.notes.GetOne(noteId)
	if getNoteErr != nil {
		return getNoteErr
	}

	if note.CreatedBy != userId {
		return fmt.Errorf("NOT_ALLOWED_TO_DELETE")
	}

	err := d.notes.DeleteOne(noteId)
	if err == nil {
		d.notesDrafts.DeleteOne(noteId)
	}
	return err
}

func (d *datasources) GetDraftOfNote(userId, noteId string) (Note, error) {
	note, getNoteErr := d.notes.GetOne(noteId)
	if getNoteErr != nil {
		return Note{}, getNoteErr
	}

	if userId != note.CreatedBy {
		return Note{}, fmt.Errorf("NOT_ALLOWED_TO_GET")
	}

	noteDraft, getDraftErr := d.notesDrafts.GetOne(noteId)
	if noteDraft.ShortId == "" {
		return Note{}, fmt.Errorf("NOT_FOUND")
	}

	return Note{
		ShortId:   note.ShortId,
		UpdatedAt: note.UpdatedAt,
		UpdatedBy: note.UpdatedBy,
		CreatedBy: note.CreatedBy,
		CreatedAt: note.CreatedAt,
		Title:     noteDraft.Title,
		Content:   noteDraft.Content,
		Tags:      noteDraft.Tags,
	}, getDraftErr
}

func (d *datasources) UpdateDraft(userId, noteId, title, content string, tags []string) error {
	note, getNoteErr := d.notes.GetOne(noteId)
	if getNoteErr != nil {
		return getNoteErr
	}

	if userId != note.CreatedBy {
		return fmt.Errorf("NOT_ALLOWED_TO_UPDATE")
	}

	if note.IsDifferent(notes.NoteDocument{Title: title, Content: content, Tags: tags}) {
		return d.notesDrafts.UpsertOne(noteId, title, content, tags)
	}
	return d.notesDrafts.DeleteOne(noteId)
}

func (d *datasources) DeleteDraft(userId, noteId string) error {
	note, getNoteErr := d.notes.GetOne(noteId)
	if getNoteErr != nil {
		return getNoteErr
	}

	if userId != note.CreatedBy {
		return fmt.Errorf("NOT_ALLOWED_TO_DELETE")
	}

	return d.notesDrafts.DeleteOne(noteId)
}

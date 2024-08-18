package notesService

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"

	"github.com/alperklc/the-zula/service/infrastructure/db/notes"
	"github.com/alperklc/the-zula/service/infrastructure/db/notesDrafts"
	mqpublisher "github.com/alperklc/the-zula/service/infrastructure/messageQueue/publisher"
	referencesService "github.com/alperklc/the-zula/service/services/references"
	usersService "github.com/alperklc/the-zula/service/services/users"
	"github.com/alperklc/the-zula/service/utils"
)

type NoteService interface {
	SearchTags(userId, searchKeyword string, limit int) ([]notes.TagsResult, error)
	GetStatistics(userId string) (Statistics, error)
	ListNotes(userId string, searchKeyword *string, page, pageSize *int, sortBy, sortDirection *string, tags *[]string) (NotesPage, error)
	CreateNote(userId, clientId string, title, content *string, tags *[]string) (Note, error)
	UpdateNote(noteId, userId, clientId string, update map[string]interface{}) error
	GetNote(noteId, userId, clientId string, params GetNoteParams) (Note, error)
	GetNotes(userId string, noteIds, fields []string) ([]Note, error)
	DeleteNote(noteId, userId, clientId string) error
	GetDraftOfNote(userId, noteId string) (Note, error)
	UpdateDraft(userId, noteId, title, content string, tags []string) error
	DeleteDraft(userId, noteId string) error
}

type datasources struct {
	users       usersService.UsersService
	notes       notes.Collection
	notesDrafts notesDrafts.Collection
	references  referencesService.ReferencesService
	mqpublisher mqpublisher.MessagePublisher
}

func NewService(u usersService.UsersService, n notes.Collection, nd notesDrafts.Collection, nrs referencesService.ReferencesService, mqp mqpublisher.MessagePublisher) NoteService {
	return &datasources{
		users: u, notes: n, notesDrafts: nd, references: nrs, mqpublisher: mqp,
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
		idsOfNotes[i] = v.ShortId
	}
	draftsOnPage, _ := d.notesDrafts.CheckExistence(idsOfNotes)

	var items []Note = make([]Note, 0, len(notes))
	for _, b := range notes {
		note := Note{b.ShortId, b.UpdatedAt, b.UpdatedBy, b.CreatedBy, b.CreatedAt, b.Title, "", b.Tags, draftsOnPage[b.ShortId], nil, nil}
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

	if err == nil {
		createdNoteBuffer := new(bytes.Buffer)
		json.NewEncoder(createdNoteBuffer).Encode(createdNote)
		createdNoteBytes := createdNoteBuffer.Bytes()

		go d.mqpublisher.Publish(mqpublisher.NoteCreated(userId, clientId, createdNote.ShortId, &createdNoteBytes))
	}

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

	if err == nil {
		updateContentBuffer := new(bytes.Buffer)
		json.NewEncoder(updateContentBuffer).Encode(update)
		updateContentBytes := updateContentBuffer.Bytes()

		go d.mqpublisher.Publish(mqpublisher.NoteUpdated(userId, clientId, noteId, &updateContentBytes))
	}

	return err
}

func (d *datasources) GetNote(noteId, userId, clientId string, params GetNoteParams) (Note, error) {
	note, getNoteErr := d.notes.GetOne(noteId)
	if getNoteErr != nil || note.Id == "" {
		return Note{}, getNoteErr
	}

	if userId != note.CreatedBy {
		return Note{}, fmt.Errorf("NOT_ALLOWED_TO_GET")
	}

	creator, errGetCreatorUser := d.users.GetUser(note.CreatedBy)
	if errGetCreatorUser != nil {
		return Note{}, errGetCreatorUser
	}
	updater, errGetUpdaterUser := d.users.GetUser(note.UpdatedBy)
	if errGetUpdaterUser != nil {
		return Note{}, errGetUpdaterUser
	}

	draftExist, _ := d.notesDrafts.CheckExistence([]string{noteId})
	response := Note{
		ShortId:   note.ShortId,
		UpdatedAt: note.UpdatedAt,
		UpdatedBy: updater.DisplayName,
		CreatedBy: creator.DisplayName,
		CreatedAt: note.CreatedAt,
		Title:     note.Title,
		Content:   note.Content,
		Tags:      note.Tags,
		HasDraft:  draftExist[note.ShortId],
	}

	if params.LoadDraft {
		draftOfNote, errGetDraft := d.notesDrafts.GetOne(noteId)
		if errGetDraft != nil {
			return response, nil
		}

		response.Content = draftOfNote.Content
		response.Title = draftOfNote.Title
		response.Tags = draftOfNote.Tags
	}

	if params.GetHistory {
		// todo
	}

	if params.GetReferences {
		references, errGetReferences := d.references.ListReferencesToNote(userId, noteId, 1)
		if errGetReferences != nil {
			return Note{}, errGetReferences
		}

		response.References = &references
	}

	go d.mqpublisher.Publish(mqpublisher.NoteRead(userId, clientId, noteId, nil))

	return response, getNoteErr
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

	if err == nil {
		go d.mqpublisher.Publish(mqpublisher.NoteDeleted(userId, clientId, noteId, nil))
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

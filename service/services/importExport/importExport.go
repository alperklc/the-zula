package importExportService

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/alperklc/the-zula/service/infrastructure/db/bookmarks"
	"github.com/alperklc/the-zula/service/infrastructure/db/notes"
	"github.com/alperklc/the-zula/service/infrastructure/db/notesChanges"
	"github.com/alperklc/the-zula/service/infrastructure/db/pageContent"
	"github.com/alperklc/the-zula/service/infrastructure/db/references"
	useractivity "github.com/alperklc/the-zula/service/infrastructure/db/userActivity"
	"github.com/samber/lo"
)

type ImportExportService interface {
	MakeZipFile(userId string) (*os.File, error)
	ProcessIncomingZipFile(zipData []byte) (ImportResult, error)
}

type datasources struct {
	notes        notes.Collection
	notesChanges notesChanges.Collection
	references   references.Collection
	bookmarks    bookmarks.Collection
	pageContent  pageContent.Collection
	useractivity useractivity.Collection
}

func NewService(n notes.Collection, nc notesChanges.Collection, nr references.Collection, b bookmarks.Collection, pc pageContent.Collection, ua useractivity.Collection) ImportExportService {
	return &datasources{
		notes: n, notesChanges: nc, references: nr, bookmarks: b, pageContent: pc, useractivity: ua,
	}
}

func SaveDocumentsToJSONFiles[T WithID](userId, collectionName string, documents []T) error {
	// Create directory for the collection
	dirPath := filepath.Join("exports/zula"+"-"+userId, collectionName)
	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		return err
	}

	// Write each document into a separate JSON file
	for _, doc := range documents {
		filePath := filepath.Join(dirPath, fmt.Sprintf("%s.json", doc.GetId()))
		file, err := os.Create(filePath)
		if err != nil {
			return err
		}

		encoder := json.NewEncoder(file)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(doc); err != nil {
			file.Close()
			return err
		}
		file.Close()
	}
	return nil
}

func ZipFolder(source, target string) (*os.File, error) {
	zipFile, err := os.Create(target)
	if err != nil {
		return nil, err
	}
	defer zipFile.Close()

	writer := zip.NewWriter(zipFile)
	defer writer.Close()

	errWalk := filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Get the relative path of the file
		relPath, err := filepath.Rel(source, path)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		zipWriter, err := writer.Create(relPath)
		if err != nil {
			return err
		}

		_, err = io.Copy(zipWriter, file)
		return err
	})

	return zipFile, errWalk
}

func (d *datasources) MakeZipFile(userId string) (*os.File, error) {
	notesOfUser, errExportNotes := d.notes.ExportForUser(userId)
	if errExportNotes != nil {
		return nil, errExportNotes
	}
	if err := SaveDocumentsToJSONFiles(userId, "notes", notesOfUser); err != nil {
		return nil, fmt.Errorf("could not export notes, %s", err)
	}

	bookmarksOfUser, errExportBookmarks := d.bookmarks.ExportForUser(userId)
	if errExportBookmarks != nil {
		return nil, errExportBookmarks
	}
	if err := SaveDocumentsToJSONFiles(userId, "bookmarks", bookmarksOfUser); err != nil {
		return nil, fmt.Errorf("could not export bookmarks, %s", err)
	}

	notesIds := lo.Map(notesOfUser, func(note notes.NoteDocument, index int) string {
		return note.Id
	})
	uniqNotesIds := lo.Uniq(notesIds)
	notesChangesOfUser, errExportChanges := d.notesChanges.Export(uniqNotesIds)
	if errExportChanges != nil {
		return nil, errExportChanges
	}
	if err := SaveDocumentsToJSONFiles(userId, "notes_changes", notesChangesOfUser); err != nil {
		return nil, fmt.Errorf("could not export notes_changes, %s", err)
	}

	referencesOfUser, errExportReferences := d.references.Export(uniqNotesIds)
	if errExportReferences != nil {
		return nil, errExportReferences
	}
	if err := SaveDocumentsToJSONFiles(userId, "references", referencesOfUser); err != nil {
		return nil, fmt.Errorf("could not export references, %s", err)
	}

	urls := lo.Map(bookmarksOfUser, func(b bookmarks.BookmarkDocument, index int) string {
		return b.URL
	})
	pageContentsOfUser, errExportContent := d.pageContent.ExportContent(urls)
	if errExportContent != nil {
		return nil, errExportContent
	}
	if err := SaveDocumentsToJSONFiles(userId, "page-content", pageContentsOfUser); err != nil {
		return nil, fmt.Errorf("could not export page-content, %s", err)
	}

	activitiesOfUser, errExportActivity := d.useractivity.ExportForUser(userId)
	if errExportActivity != nil {
		return nil, errExportActivity
	}
	if err := SaveDocumentsToJSONFiles(userId, "users-activity", activitiesOfUser); err != nil {
		return nil, fmt.Errorf("could not export users-activity, %s", err)
	}

	return ZipFolder("exports/zula"+"-"+userId, "exports/zula"+"-"+userId+".zip")
}

func (d *datasources) ProcessIncomingZipFile(zipData []byte) (ImportResult, error) {
	result := ImportResult{}

	reader, err := zip.NewReader(bytes.NewReader(zipData), int64(len(zipData)))
	if err != nil {
		return result, fmt.Errorf("failed to read ZIP file: %w", err)
	}

	var notesItems []notes.NoteDocument
	var notesChangesItems []notesChanges.NotesChangesDocument
	var referencesItems []references.ReferencesDocument
	var bookmarksItems []bookmarks.BookmarkDocument
	var pageContentItems []pageContent.PageContentDocument
	var useractivityItems []useractivity.UserActivityDocument

	for _, file := range reader.File {
		if file.FileInfo().IsDir() {
			continue // Skip directories
		}

		// Extract folder name (which will be the collection name)
		folderName := filepath.Dir(file.Name)
		if folderName == "." || folderName == "" {
			continue
		}

		zippedFile, err := file.Open()
		if err != nil {
			return result, fmt.Errorf("failed to open file %s: %w", file.Name, err)
		}
		defer zippedFile.Close()

		content, err := io.ReadAll(zippedFile)
		if err != nil {
			return result, fmt.Errorf("failed to read content of %s: %w", file.Name, err)
		}

		switch folderName {
		case "notes":
			var note notes.NoteDocument
			err := json.Unmarshal(content, &note)
			if err == nil {
				notesItems = append(notesItems, note)
			}
		case "notes_changes":
			var notesChange notesChanges.NotesChangesDocument
			err := json.Unmarshal(content, &notesChange)
			if err == nil {
				notesChangesItems = append(notesChangesItems, notesChange)
			}
		case "references":
			var ref references.ReferencesDocument
			err := json.Unmarshal(content, &ref)
			if err == nil {
				referencesItems = append(referencesItems, ref)
			}
		case "bookmarks":
			var bookmark bookmarks.BookmarkDocument
			err := json.Unmarshal(content, &bookmark)
			if err == nil {
				bookmarksItems = append(bookmarksItems, bookmark)
			}
		case "page-content":
			var pc pageContent.PageContentDocument
			err := json.Unmarshal(content, &pc)
			if err == nil {
				pageContentItems = append(pageContentItems, pc)
			}
		case "users-activity":
			var ua useractivity.UserActivityDocument
			err := json.Unmarshal(content, &ua)
			if err == nil {
				useractivityItems = append(useractivityItems, ua)
			}
		}
	}

	if len(notesItems) > 0 {
		result.Notes.Total = len(notesItems)
		importedCount, errImportNotes := d.notes.ImportMany(notesItems)
		if errImportNotes != nil {
			return result, fmt.Errorf("failed to import notes %s", errImportNotes)
		}

		result.Notes.ImportedCount = importedCount
	}

	if len(notesChangesItems) > 0 {
		result.NotesChanges.Total = len(notesChangesItems)
		importedCount, errImportNotesChanges := d.notesChanges.ImportMany(notesChangesItems)
		if errImportNotesChanges != nil {
			return result, fmt.Errorf("failed to import notes changes %s", errImportNotesChanges)
		}
		result.NotesChanges.ImportedCount = importedCount
	}

	if len(referencesItems) > 0 {
		result.References.Total = len(referencesItems)
		importedCount, errImportReferences := d.references.ImportMany(referencesItems)
		if errImportReferences != nil {
			return result, fmt.Errorf("failed to import references %s", errImportReferences)
		}

		result.References.ImportedCount = importedCount
	}

	if len(bookmarksItems) > 0 {
		result.Bookmarks.Total = len(bookmarksItems)
		importedCount, errImportBookmarks := d.bookmarks.ImportMany(bookmarksItems)
		if errImportBookmarks != nil {
			return result, fmt.Errorf("failed to import bookmarks %s", errImportBookmarks)
		}

		result.Bookmarks.ImportedCount = importedCount
	}

	if len(pageContentItems) > 0 {
		result.PageContent.Total = len(pageContentItems)
		importedCount, errImportPageContent := d.pageContent.ImportMany(pageContentItems)
		if errImportPageContent != nil {
			return result, fmt.Errorf("failed to import page content %s", errImportPageContent)
		}

		result.PageContent.ImportedCount = importedCount
	}

	if len(useractivityItems) > 0 {
		result.Useractivity.Total = len(useractivityItems)
		importedCount, errImportActivities := d.useractivity.ImportMany(useractivityItems)
		if errImportActivities != nil {
			return result, fmt.Errorf("failed to import user activities %s", errImportActivities)
		}

		result.Useractivity.ImportedCount = importedCount
	}

	return result, nil
}

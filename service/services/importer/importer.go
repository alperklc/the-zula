package importerService

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"path/filepath"

	"github.com/alperklc/the-zula/service/infrastructure/db/bookmarks"
	"github.com/alperklc/the-zula/service/infrastructure/db/notes"
	"github.com/alperklc/the-zula/service/infrastructure/db/pageContent"
	"github.com/alperklc/the-zula/service/infrastructure/db/references"
	useractivity "github.com/alperklc/the-zula/service/infrastructure/db/userActivity"
)

type ImporterService interface {
	ProcessZipFile(zipData []byte) error
}

type datasources struct {
	notes        notes.Collection
	references   references.Collection
	bookmarks    bookmarks.Collection
	pageContent  pageContent.Collection
	useractivity useractivity.Collection
}

func NewService(n notes.Collection, nr references.Collection, b bookmarks.Collection, pc pageContent.Collection, ua useractivity.Collection) ImporterService {
	return &datasources{
		notes: n, references: nr, bookmarks: b, pageContent: pc, useractivity: ua,
	}
}

func (d *datasources) ProcessZipFile(zipData []byte) error {
	reader, err := zip.NewReader(bytes.NewReader(zipData), int64(len(zipData)))
	if err != nil {
		return fmt.Errorf("failed to read ZIP file: %w", err)
	}

	var notesItems []notes.NoteDocument
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
			return fmt.Errorf("failed to open file %s: %w", file.Name, err)
		}
		defer zippedFile.Close()

		content, err := io.ReadAll(zippedFile)
		if err != nil {
			return fmt.Errorf("failed to read content of %s: %w", file.Name, err)
		}

		switch folderName {
		case "notes":
			var note notes.NoteDocument
			err := json.Unmarshal(content, &note)
			if err == nil {
				notesItems = append(notesItems, note)
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
		case "pageContent":
			var pc pageContent.PageContentDocument
			err := json.Unmarshal(content, &pc)
			if err == nil {
				pageContentItems = append(pageContentItems, pc)
			}
		case "useractivity":
			var ua useractivity.UserActivityDocument
			err := json.Unmarshal(content, &ua)
			if err == nil {
				useractivityItems = append(useractivityItems, ua)
			}
		}
	}

	if len(notesItems) > 0 {
		d.notes.ImportMany(notesItems)
	}

	if len(referencesItems) > 0 {
		d.references.ImportMany(referencesItems)
	}

	if len(bookmarksItems) > 0 {
		d.bookmarks.ImportMany(bookmarksItems)
	}

	if len(pageContentItems) > 0 {
		d.pageContent.ImportMany(pageContentItems)
	}

	if len(useractivityItems) > 0 {
		d.useractivity.ImportMany(useractivityItems)
	}

	return nil
}

package importExportService

type CollectionResult struct {
	Total         int `json:"total"`
	ImportedCount int `json:"importedCount"`
}

type ImportResult struct {
	Notes        CollectionResult `json:"notes"`
	NotesChanges CollectionResult `json:"notesChanges"`
	References   CollectionResult `json:"references"`
	Bookmarks    CollectionResult `json:"bookmarks"`
	PageContent  CollectionResult `json:"pageContent"`
	Useractivity CollectionResult `json:"useractivity"`
}

type WithID interface {
	GetId() string
}

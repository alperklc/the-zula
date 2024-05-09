package notesReferences

type NotesReferencesDocument struct {
	From string `bson:"from"`
	To   string `bson:"to"`
}

type ReferencesDocument struct {
	From string `bson:"from" json:"from"`
	To   string `bson:"to" json:"to"`
}

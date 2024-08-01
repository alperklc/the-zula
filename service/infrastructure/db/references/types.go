package references

type ReferencesDocument struct {
	From string `bson:"from" json:"from"`
	To   string `bson:"to" json:"to"`
}

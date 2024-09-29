package references

type ReferencesDocument struct {
	Id   string `bson:"_id,omitempty"`
	From string `bson:"from" json:"from"`
	To   string `bson:"to" json:"to"`
}

func (n ReferencesDocument) GetId() string {
	return n.Id
}

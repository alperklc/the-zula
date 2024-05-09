package notesDrafts

type NoteDraftDocument struct {
	UserId  string   `bson:"userId"`
	ShortId string   `bson:"shortId"`
	Title   string   `bson:"title"`
	Content string   `bson:"content"`
	Tags    []string `bson:"tags"`
}

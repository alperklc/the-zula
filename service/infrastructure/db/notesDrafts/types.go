package notesDrafts

type NoteDraftDocument struct {
	UserId  string   `bson:"userId"`
	Id      string   `bson:"id"`
	Title   string   `bson:"title"`
	Content string   `bson:"content"`
	Tags    []string `bson:"tags"`
}

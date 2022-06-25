package Entity

type MongoBenchEntity struct {
	Name    string `bson:"Name" json:"Name"`
	Counter int    `bson:"Counter" json:"Counter"`
}

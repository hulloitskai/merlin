package models

// A Note is an annotation that can be associated with an Item.
type Note struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Link  string `json:"link"`
}

// Notes are a set of Note objects.
type Notes []*Note

package balance

// A Note is an annotation that can be associated with an Item.
type Note struct {
	ID int `json:"id"`
}

// Notes are a set of Note objects.
type Notes []*Note

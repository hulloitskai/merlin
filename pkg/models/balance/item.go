package balance

import "fmt"

// An Item is one record on a Sheet. It contains a name, a value, and may
// contain additional notes.
type Item struct {
	Name  string `json:"name"`
	Value string `json:"value"`
	Notes []int  `json:"notes,omitempty"`
}

func (i *Item) String() string {
	return fmt.Sprintf("Item{Name: %s, Value: %s, Notes: %v}", i.Name,
		i.Value, i.Notes)
}

// Items are a set of Item objects.
type Items []*Item

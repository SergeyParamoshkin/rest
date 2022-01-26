package model

// Article data model. I suggest looking at https://upper.io for an easy
// and powerful data persistence adapter.
type Article struct {
	ID     string `json:"id"`
	UserID int64  `json:"userId"` // the author
	Title  string `json:"title"`
	Slug   string `json:"slug"`
}

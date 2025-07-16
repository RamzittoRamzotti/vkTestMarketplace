package models

type Ad struct {
	ID             int     `db:"id" json:"id"`
	Title          string  `db:"title" json:"title"`
	Description    string  `db:"description" json:"description"`
	ImageURL       string  `db:"image_url" json:"image_url"`
	Price          float64 `db:"price" json:"price"`
	AuthorID       int     `db:"author_id" json:"-"`
	AuthorUsername string  `db:"author_username" json:"author,omitempty"`
}

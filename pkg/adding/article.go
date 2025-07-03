package adding

type Article struct {
	ID       int    `json:"id"`
	AuthorID int    `json:"author_id" validate:"required"`
	Title    string `json:"title" validate:"required"`
	Body     string `json:"body" validate:"required"`
}

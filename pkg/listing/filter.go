package listing

type FilterGetArticle struct {
	AuthorName string `schema:"author_name"`
	Search     string `schema:"search"`

	Limit int `schema:"limit"`
	Page  int `schema:"page"`
}

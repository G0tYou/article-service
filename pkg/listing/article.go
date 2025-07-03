package listing

type Article struct {
	ID     int    `json:"id"`
	Author Author `json:"author"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

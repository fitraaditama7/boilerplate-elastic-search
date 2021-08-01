package post

type CreateRequest struct {
	Title string   `json:"title"`
	Text  string   `json:"text"`
	Tags  []string `json:"tags"`
}

type UpdateRequest struct {
	ID    string   `json:"id"`
	Title string   `json:"title"`
	Text  string   `json:"text"`
	Tags  []string `json:"tags"`
}

type DeleteRequest struct {
	ID string `json:"id"`
}

type FindRequest struct {
	ID string `json:"id"`
}

type SearchRequest struct {
	Keyword string `json:"keyword"`
}

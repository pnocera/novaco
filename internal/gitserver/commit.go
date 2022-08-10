package gitserver

type Commit struct {
	Message    string `json:"message"`
	ID         string `json:"id"`
	ObjectType string `json:"object_type"`
	Author     string `json:"author"`
}

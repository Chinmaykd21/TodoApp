package customDataStructs

type Todo struct {
	TodoId      int    `json:"todoId"`
	Title       string `json:"title"`
	Body        string `json:"body"`
	IsCompleted bool   `json:"isCompleted"`
}

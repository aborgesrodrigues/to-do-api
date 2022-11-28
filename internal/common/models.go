package common

type TaskState string

type User struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
}

type Task struct {
	Id          string    `json:"id"`
	UserId      string    `json:"user_id"`
	Description string    `json:"description"`
	State       TaskState `json:"state"`
	User        *User     `json:"user,omitempty"`
}

type Metadata struct {
	Name  string
	Value interface{}
}

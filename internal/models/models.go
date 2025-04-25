package models

type NewTask struct {
	ID int `json:"id,omitempty"`
}

type Task struct {
	ID      string `json:"id,omitempty"`
	Date    string `json:"date,omitempty"` //*time.Time
	Title   string `json:"title,omitempty"`
	Comment string `json:"comment,omitempty"`
	Repeat  string `json:"repeat,omitempty"`
}

type TasksResp struct {
	Tasks []*Task `json:"tasks"`
}

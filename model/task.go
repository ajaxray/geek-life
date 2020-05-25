package model

type Task struct {
	ID        int64  `storm:"id,increment",json:"id"`
	ProjectID int64  `storm:"index",json:"project_id"`
	UUID      string `storm:"unique",json:"uuid,omitempty"`
	Title     string `json:"text"`
	Details   string `json:"notes"`
	Completed bool   `storm:"index",json:"completed"`
	DueDate   int64  `storm:"index",json:"due_date,omitempty"`
}

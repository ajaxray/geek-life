package model

// Task represent a task - the building block of the TaskManager app
type Task struct {
	ID        int64  `storm:"id,increment",json:"id"`
	ProjectID int64  `storm:"index",json:"project_id"`
	Title     string `json:"text"`
	Details   string `json:"notes"`
	Completed bool   `storm:"index",json:"completed"`
	DueDate   int64  `storm:"index",json:"due_date,omitempty"`

	// Related to integration
	ModifiedAt    int64
	Integration   string
	IntegrationID string `storm:"unique",storm:"index"`
}

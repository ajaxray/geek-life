package model

// Project represent a collection of related tasks (tags of Habitica)
type Project struct {
	ID    int64  `storm:"id,increment",json:"id"`
	Title string `storm:"index",json:"title"`

	// Related to integration
	ModifiedAt    int64
	Integration   string
	IntegrationID string `storm:"unique",storm:"index"`
}

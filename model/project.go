package model

type Project struct {
	ID    int64  `storm:"id,increment",json:"id"`
	Title string `storm:"index",json:"title"`
	UUID  string `storm:"unique",json:"uuid,omitempty"`
}

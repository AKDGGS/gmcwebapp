package model

type File struct {
	ID   int32  `json:"file_id"`
	Name string `json:"file_name"`
	Size string `json:"file_size"`
}

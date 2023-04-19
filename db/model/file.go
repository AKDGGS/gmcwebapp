package model

type File struct {
	ID          int32  `json:"file_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Size        string `json:"size"`
	Type        string `json:"mimetype"`
}

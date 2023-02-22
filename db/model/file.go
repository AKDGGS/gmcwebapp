package model

type File struct {
	ID   int32  `json:"file_id"`
	Name string `json:"name"`
	Size string `json:"size"`
	Type string `json:"mimetype"`
}

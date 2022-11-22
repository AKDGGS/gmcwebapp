package model

type File struct {
	ID       int32  `json:"id"`
	FileName string `json:"file_name"`
	FileSize string `json:"file_size"`
}

package model

type File struct {
	ID          int     `json:"file_id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Size        int64   `json:"size"`
	Type        string  `json:"mimetype"`
	MD5         string  `json:"content_md5"`
	WellIDs     []int32 `json:"well_ids"`
}

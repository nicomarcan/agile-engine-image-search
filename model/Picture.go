package model

type PicturesResponse struct {
	Pictures  []Picture `json:"pictures,omitempty"`
	Page      int       `json:"page"`
	PageCount int       `json:"pageCount"`
	HasMore   bool      `json:"hasMore"`
}

type Picture struct {
	Id             string `json:"id"`
	CroppedPicture string `json:"cropped_picture"`
	Author         string `json:"author,omitempty"`
	Camera         string `json:"camera,omitempty"`
	Tags           string `json:"tags,omitempty"`
	FullPicture    string `json:"full_picture,omitempty"`
}

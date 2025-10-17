package api

type ResBase struct {
	IsSuccess bool   `json:"is_success"`
	Err       string `json:"err"`
}

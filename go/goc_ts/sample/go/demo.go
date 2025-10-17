package api

const URI_ListUser = "/user/list"

type Pagination struct {
	PageNum  int `json:"page_num"`
	PageSize int `json:"page_size"`
}

type ListUserReq struct {
	Operator string     `json:"operator"`
	Page     Pagination `json:"page"`
}

type ListUserRes struct {
	Res     ResBase  `json:"res"`
	Summary int64    `json:"summary"`
	Users   []string `json:"users"`
}

const URI_CreateUser = "/user/create"

type CreateUserReq struct{}

type CreateUserRes struct{}

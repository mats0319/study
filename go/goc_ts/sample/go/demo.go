package api

type UserIdentify int8

const (
	UserIdentify_Value0 UserIdentify = 0
	UserIdentify_Value1 UserIdentify = 1
	UserIdentify_Value2 UserIdentify = 2
)

const URI_ListUser = "/user/list"

type Pagination struct {
	PageNum  int `json:"page_num"`
	PageSize int `json:"page_size"`
}

type ListUserReq struct {
	Operator     string       `json:"operator"`
	ListIdentify UserIdentify `json:"list_identify"`
	Page         Pagination   `json:"page"`
}

type ListUserRes struct {
	Res     ResBase  `json:"res"`
	Summary int64    `json:"summary"`
	Users   []string `json:"users"`
}

const URI_CreateUser = "/user/create"

type CreateUserReq struct{}

type CreateUserRes struct{}

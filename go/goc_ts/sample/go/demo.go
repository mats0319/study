package api

const URI_ListUser = "/user/list"

type UserIdentify int8

const (
	UserIdentify_Administrator UserIdentify = 0
	UserIdentify_VIP           UserIdentify = 1
	UserIdentify_Visitor       UserIdentify = 2
)

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

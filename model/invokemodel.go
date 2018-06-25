package model

type LoginReq struct {
	Username string `json:"username"`
	Password string	`json:"password"`
}

type LoginResponse struct {
	Code int `json:"code"`
	Data LoginDataResponse `json:"data"`
}

type LoginDataResponse struct {
	Token    string `json:"token"`
	OrgExist bool  `json:"org_exist"`
}

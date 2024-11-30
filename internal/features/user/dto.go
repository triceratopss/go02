package user

type CreateUpdateUserRequest struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type GetUserListRequest struct {
	Limit  int `query:"limit"`
	Offset int `query:"offset"`
}

type GetUserListResponse struct {
	Users []GetUserResponse `json:"users"`
}

type GetUserResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

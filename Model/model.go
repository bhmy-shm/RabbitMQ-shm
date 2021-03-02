package Model

type UserModel struct {
	UserID   int    `json:"uid"`
	UserName string `json:"uname" binding:"required"`
}

func NewUserModel() *UserModel {
	return &UserModel{}
}

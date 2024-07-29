package domain

type SignupRequest struct {
	Username       string `json:"username" binding:"required,min=5,max=20"`
	Password       string `json:"password" binding:"required,min=5,max=20"`
	SecondPassword string `json:"second_password" binding:"required,min=6,max=10"`
}

type RetrieveRequest struct {
	Username       string `json:"username" binding:"required,min=5,max=20"`
	SecondPassword string `json:"second_password" binding:"required,min=6,max=10"`
	NewPassword    string `json:"new_password" binding:"required,min=5,max=20"`
}

type GetAccountResponse struct {
	Username        string `json:"username" binding:"required,min=5,max=20"`
	IsForeverBanned bool   `json:"is_forever_banned"`
	CashPoint       uint32 `json:"cash_point"`
	MaplePoint      uint32 `json:"maple_point"`
}

type UpdatePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required,min=5,max=20"`
	NewPassword string `json:"new_password" binding:"required,min=5,max=20"`
}

type UpdateSecondPasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required,min=6,max=10"`
	NewPassword string `json:"new_password" binding:"required,min=6,max=10"`
}

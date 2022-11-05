package model

type User struct {
	Id      int
	Balance float64
}

type UserDto struct {
	Id      int     `json:"id"`
	Balance float64 `json:"balance"`
}

type UserDepositRequest struct {
	Id      int     `json:"id" binding:"required"`
	Deposit float64 `json:"deposit" binding:"required"`
}

func (u *User) ToDto() UserDto {
	return UserDto{
		Id:      u.Id,
		Balance: u.Balance,
	}
}

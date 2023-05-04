package domain

type UserRepository interface {
	CreateUser(request TransferUser)
}

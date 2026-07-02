package users_transport_http

import "github.com/vmkzy/todoapp-go/internal/core/domain"

type UserDTORepsonse struct {
	ID          int     `json:"id"`
	Version     int     `json:"version"`
	FullName    string  `json:"full_name"`
	PhoneNumber *string `json:"phone_number"`
}

func userDTOFromDomain(user domain.User) UserDTORepsonse {
	return UserDTORepsonse{
		ID:          user.ID,
		Version:     user.Version,
		FullName:    user.FullName,
		PhoneNumber: user.PhoneNumber,
	}
}
func usersDTOFromDomains(users []domain.User) []UserDTORepsonse {
	usersDTO := make([]UserDTORepsonse, len(users))
	for i, user := range users {
		usersDTO[i] = userDTOFromDomain(user)
	}
	return usersDTO
}

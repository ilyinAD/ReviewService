package domain

type User struct {
	UserID   string
	IsActive bool
	Username string
	TeamName string
}

func NewUser(userID string, isActive bool, username string, teamName string) *User {
	return &User{
		UserID:   userID,
		IsActive: isActive,
		Username: username,
		TeamName: teamName,
	}
}

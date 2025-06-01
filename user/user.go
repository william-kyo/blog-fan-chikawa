package user

import "time"

type User struct {
	ID        int
	Nickname  string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

var UserList = []User{
	{
		ID:        1,
		Nickname:  "fanchiikawa",
		Email:     "fanchiikawa@fanchiikawa.com",
		CreatedAt: time.Now(),
	},
	{
		ID:        2,
		Nickname:  "kyo",
		Email:     "kyo@fanchiikawa.com",
		CreatedAt: time.Now(),
	},
	{
		ID:        3,
		Nickname:  "yuki",
		Email:     "yuki@fanchiikawa.com",
		CreatedAt: time.Now(),
	},
}

func GetUserList() []User {
	return UserList
}

func SaveUser(user User) {
	UserList = append(UserList, user)
}

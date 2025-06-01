package user

import "time"

type User struct {
	ID        int64     `json:"id"`
	Nickname  string    `json:"nickname"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
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

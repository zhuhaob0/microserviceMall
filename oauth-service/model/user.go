package model

type UserDetails struct {
	// 用户标识
	UserId int64
	// 用户名 唯一
	Username string
	// 用户密码
	Password string
	// 用户具有的权限
	Authorities []string // 具备的权限
}

func (u *UserDetails) IsMatch(username string, password string) bool {
	return u.Username == username && u.Password == password
}

// 密码授权时使用
type UserInfo struct {
	Username string `json:"username"`
	Password string `json:"password"`
	IdType   int    `json:"id_type"`
}

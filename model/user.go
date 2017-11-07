package model

type User struct {
	ID         int64  `gorm:primary key;not_nil`
	Login      string `gorm:not_nil`
	Pass       string `gorm:not_nil`
	WorkNumber int32
}

func init() {
	if DBConn == nil {
		GormInit()
	}
}

func (u *User) Get(login string, pass string) error {
	return DBConn.Where("login = ? and pass = ?", login, pass).First(u).Error
}

func (u *User) Save() error {
	return DBConn.Save(u).Error
}

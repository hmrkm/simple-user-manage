package usecase

type User struct {
	Id             string `json:"id" gorm:"column:id"`
	Email          string `json:"email" gorm:"column:email"`
	HashedPassword string `json:"hashed_password" gorm:"column:hashed_password"`
}

func (User) TableName() string {
	return "users"
}

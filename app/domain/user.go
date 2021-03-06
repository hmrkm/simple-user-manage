package domain

type User struct {
	Id             string `json:"id" gorm:"column:id"`
	Email          string `json:"email" gorm:"column:email"`
	HashedPassword string `json:"password" gorm:"column:password"`
}

func (User) TableName() string {
	return "users"
}

func (u User) Update(s Store, value User) error {
	return s.Update(&u, map[string]interface{}{
		"email": value.Email,
	})
}

func (u User) UpdateWithPassword(s Store, value User) error {
	return s.Update(&u, map[string]interface{}{
		"email":    value.Email,
		"password": value.HashedPassword,
	})
}

func (u User) Delete(s Store) error {
	return s.Delete(&u)
}

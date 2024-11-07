package service

import (
	"seed-sync/driver/db"
	"time"

	"gorm.io/gorm"
)

const (
	TwoFactorTypeEmail = "email"
	TwoFactorTypeAuth  = "auth"
)

// 用户表
type User struct {
	Id            int       `json:"id" gorm:"column:id"`
	Username      string    `json:"username" gorm:"column:username"`
	Password      string    `json:"password" gorm:"column:password"`
	Token         string    `json:"token" gorm:"column:token"`
	Status        string    `json:"status" gorm:"column:status"`
	IsTwoFactor   bool      `json:"is_two_factor" gorm:"column:is_two_factor"`
	TwoFactorType string    `json:"two_factor_type" gorm:"column:two_factor_type"`
	CreateTime    time.Time `json:"create_time" gorm:"column:create_time"`
	UpdateTime    time.Time `json:"update_time" gorm:"column:update_time"`
}

func (u *User) TableName() string {
	return "seed_sync_user"
}

// 根据用户名获取用户
func GetUserByUsername(username string) (*User, error) {
	user := &User{}
	if err := db.DB.Where("username = ?", username).First(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

// 创建或更新用户
func CreateOrUpdateUser(user *User) error {
	userData, err := GetUserByUsername(user.Username)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return db.DB.Create(user).Error
		}
		return err
	}
	// 更新
	user.Id = userData.Id
	return db.DB.Save(user).Error
}

package user

import (
	"seed-sync/db"
	"time"

	"gorm.io/gorm"
)

type UserDAO struct {
	db *gorm.DB
}

var userDAO = &UserDAO{
	db: db.DB,
}

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
func (dao *UserDAO) GetUserByUsername(username string) *User {
	user := &User{}
	if err := dao.db.Where("username = ?", username).First(user).Error; err != nil {
		return nil
	}
	return user
}

// 创建或更新用户
func (dao *UserDAO) CreateOrUpdateUser(user *User) error {
	userData := dao.GetUserByUsername(user.Username)
	if userData == nil {
		return dao.db.Create(user).Error
	}
	// 更新
	user.Id = userData.Id
	return dao.db.Save(user).Error
}

// 统计用户数量
func (dao *UserDAO) CountUser() int64 {
	var count int64
	dao.db.Model(&User{}).Count(&count)
	return count
}

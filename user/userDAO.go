package user

import (
	"seed-sync/db"
	"time"

	"gorm.io/gorm"
)

type UserDAO struct {
	db *gorm.DB
}

var userDAO *UserDAO

func InitUserDAO() {
	userDAO = &UserDAO{
		db: db.DB,
	}
}

type UserTable struct {
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

func (u *UserTable) TableName() string {
	return "seed_sync_user"
}

// 查询用户
func (dao *UserDAO) GetUser() *UserTable {
	user := &UserTable{}
	if err := dao.db.First(user).Error; err != nil {
		return nil
	}
	return user
}

// 创建或更新用户
func (dao *UserDAO) CreateOrUpdateUser(user *UserTable) error {
	userData := dao.GetUser()
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
	dao.db.Model(&UserTable{}).Count(&count)
	return count
}

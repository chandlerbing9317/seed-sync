package user

import (
	"sync"
	"time"
)

// 用户状态
const (
	NORMAL         = "normal"
	NOT_AUTHORIZED = "not_authorized"
	BAN            = "ban"
)

const (
	DEFAULT_USERNAME = "admin"
	DEFAULT_PASSWORD = "seed-sync"
)

type UserServiceType struct {
	userDAO *UserDAO
	lock    sync.Mutex
}

var UserService = &UserServiceType{
	userDAO: userDAO,
	lock:    sync.Mutex{},
}

// 初始化，如果数据库中不存在用户就创建一个默认的用户名和密码
func init() {
	UserService.lock.Lock()
	defer UserService.lock.Unlock()
	count := UserService.userDAO.CountUser()
	if count > 0 {
		return
	}
	//创建用户
	UserService.CreateUser(&CreateUserRequest{
		Username: DEFAULT_USERNAME,
		Password: DEFAULT_PASSWORD,
		Status:   NOT_AUTHORIZED,
	})
}

func (service *UserServiceType) CreateUser(user *CreateUserRequest) error {
	userData := &User{
		Username:      user.Username,
		Password:      user.Password,
		Status:        user.Status,
		IsTwoFactor:   false,
		TwoFactorType: "",
		CreateTime:    time.Now(),
		UpdateTime:    time.Now(),
	}
	return service.userDAO.CreateOrUpdateUser(userData)
}

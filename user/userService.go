package user

import (
	"errors"
	"seed-sync/log"
	"seed-sync/scheduler"
	"seed-sync/seedSyncServer"
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

type userService struct {
	userDAO *UserDAO
	lock    sync.Mutex
}

var UserService *userService

// 初始化
var once sync.Once

func InitUserService() {
	once.Do(func() {
		UserService = &userService{
			userDAO: userDAO,
			lock:    sync.Mutex{},
		}
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
	})
}

func (service *userService) CreateUser(user *CreateUserRequest) error {
	userData := &UserTable{
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

// 查询用户，当前系统仅支持单用户
func (service *userService) GetUser() *UserTable {
	return service.userDAO.GetUser()
}

// 定时任务，定时从服务器同步用户是否可用
func (service *userService) CheckUserScheduler(schedulerTask *scheduler.SchedulerTaskTable) error {
	return service.CheckUser()
}

func (service *userService) CheckUser() error {
	service.lock.Lock()
	defer service.lock.Unlock()
	user := service.GetUser()
	if user == nil {
		log.Fatal("当前系统不存在用户")
		return errors.New("当前系统不存在用户")
	}
	userStatus, err := seedSyncServer.SeedSyncServerService.GetUserStatus(user.Token)
	if err != nil {
		return err
	}
	// 更新用户状态
	user.Status = userStatus
	service.userDAO.CreateOrUpdateUser(user)

	return nil
}

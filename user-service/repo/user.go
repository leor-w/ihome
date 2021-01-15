package repo

import (
	"github.com/leor-w/ihome/user-service/model"
	"gorm.io/gorm"
)

// UserRepositoryInterface 用户数据库接口
type UserRepositoryInterface interface {
	Create(user *model.User) error
	Get(id uint) (*model.User, error)
	GetByEmail(email string) (*model.User, error)
	GetAll() ([]*model.User, error)
	Update(user *model.User) error
}

// UserRepository 用户数据库接口实现
type UserRepository struct {
	Db *gorm.DB
}

// Create 数据库操作 创建用户
func (repo *UserRepository) Create(user *model.User) error {
	err := repo.Db.Create(user).Error
	if err != nil {
		return err
	}
	return nil
}

// Get 数据库操作 获取指定 id 的用户
func (repo *UserRepository) Get(id uint) (*model.User, error) {
	user := &model.User{Id: id}
	err := repo.Db.First(&user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetByEmail 数据库操作 获取指定 email 的用户
func (repo *UserRepository) GetByEmail(email string) (*model.User, error) {
	user := &model.User{}
	err := repo.Db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetAll 数据库操作 获取所有用户数据
func (repo *UserRepository) GetAll() ([]*model.User, error) {
	users := []*model.User{}
	err := repo.Db.Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

// Update 数据库操作 修改用户
func (repo *UserRepository) Update(user *model.User) error {
	err := repo.Db.Save(user).Error
	if err != nil {
		return err
	}
	return nil
}

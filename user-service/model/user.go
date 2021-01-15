package model

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	pb "github.com/leor-w/ihome/user-service/proto/user"
)

type User struct {
	Id            uint   `gorm:"primary_key;autoIncrement;<-:create"`
	Name          string `gorm:"type:varchar(100)"`
	Email         string `gorm:"type:varchar(100)"`
	Password      string
	Status        uint8 `gorm:"default:0"`
	RememberToken string
	CreatedAt     time.Time `gorm:"<-:create"`
	UpdatedAt     time.Time
}

// ToModel 将 gRPC 模型转换为数据库模型
func (model User) ToModel(req *pb.User) error {
	if req.Id == "" {
		id, err := strconv.ParseUint(req.Id, 10, 64)
		if err != nil {
			return err
		}
		model.Id = uint(id)
	}
	if req.Name == "" {
		return errors.New("user name cannot empty")
	}
	model.Name = req.Name

	if req.Email == "" {
		return errors.New("user email cannot empty")
	}
	model.Email = req.Email

	if req.Password == "" {
		return errors.New("user password cannot empty")
	}
	model.Password = req.Password

	if req.Status != "" {
		status, err := strconv.ParseUint(req.Status, 10, 64)
		if err != nil {
			return err
		}
		model.Status = uint8(status)
	}

	if req.RememberToken != "" {
		model.RememberToken = req.RememberToken
	}
	return nil
}

// ToProto 将数据库模型转换为 gRPC 模型
func (model *User) ToProto() *pb.User {
	user := &pb.User{}
	user.Id = fmt.Sprintf("%d", model.Id)
	user.Name = model.Name
	user.Password = model.Password
	user.Email = model.Email
	user.Status = fmt.Sprintf("%d", model.Status)
	user.RememberToken = model.RememberToken
	user.CreatedAt = model.CreatedAt.Format("2006-01-02 15:04:05")
	user.UpdatedAt = model.UpdatedAt.Format("2006-01-02 15:04:05")
	return user
}

// UserProtoToModels 将数据库模型数组转换为 gRPC 数组模型
func UserProtoToModels(users []*User) []*pb.User {
	protoUsers := []*pb.User{}
	for _, user := range users {
		protoUsers = append(protoUsers, user.ToProto())
	}
	return protoUsers
}

package handler

import (
	"context"
	"errors"
	"strconv"

	"github.com/sirupsen/logrus"

	"golang.org/x/crypto/bcrypt"

	"github.com/leor-w/ihome/user-service/model"
	"gorm.io/gorm"

	pb "github.com/leor-w/ihome/user-service/proto/user"
	"github.com/leor-w/ihome/user-service/repo"
	"github.com/leor-w/ihome/user-service/service"
)

type UserService struct {
	Repo  repo.UserRepositoryInterface
	Token service.AuthInterface
}

// Create 创建用户
func (svc *UserService) Create(_ context.Context, req *pb.User, resp *pb.Response) error {
	user := &model.User{}
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	req.Password = string(hashedPwd)
	err = user.ToModel(req)
	if err != nil {
		return err
	}
	err = svc.Repo.Create(user)
	if err != nil {
		return err
	}
	resp.User = user.ToProto()
	return nil
}

// Get 获取用户
func (svc *UserService) Get(_ context.Context, req *pb.User, resp *pb.Response) error {
	var (
		user     *model.User
		respUser = &pb.User{}
		id       uint64
		err      error
	)

	if req.Id != "" {
		id, err = strconv.ParseUint(req.Id, 10, 64)
		user, err = svc.Repo.Get(uint(id))
	} else if req.Email != "" {
		user, err = svc.Repo.GetByEmail(req.Email)
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	if user != nil {
		respUser = user.ToProto()
	}

	resp.User = respUser
	return nil
}

// GetAll 获取所有用户
func (svc *UserService) GetAll(_ context.Context, _ *pb.Request, resp *pb.Response) error {
	users, err := svc.Repo.GetAll()
	if err != nil {
		return err
	}
	resp.Users = model.UserProtoToModels(users)
	return nil
}

// Update 用户修改密码
func (svc *UserService) Update(_ context.Context, req *pb.User, resp *pb.Response) error {
	if req.Id == "" {
		return errors.New("user id cannot empty")
	}
	if req.Password == "" {
		return errors.New("user password cannot empty")
	}
	hashPwd, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	req.Password = string(hashPwd)
	user := &model.User{}
	err = user.ToModel(req)
	if err != nil {
		return err
	}
	err = svc.Repo.Update(user)
	if err != nil {
		return err
	}
	resp.User = user.ToProto()
	return nil
}

// Auth 登陆并生成 jwt
func (svc *UserService) Auth(_ context.Context, req *pb.User, resp *pb.Token) error {
	logrus.Infof("loging with : [%s] [%s]", req.Email, req.Password)

	user, err := svc.Repo.GetByEmail(req.Email)
	if err != nil {
		logrus.Errorf("GetByEmail faild: %s", err.Error())
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(req.Password), []byte(user.Password))
	if err != nil {
		logrus.Errorf("password fault")
		return err
	}

	token, err := svc.Token.Encode(user)
	if err != nil {
		logrus.Errorf("encode token string failed: user [%v], error [%s]", user, err)
		return err
	}

	resp.Token = token
	return nil
}

// ValidateToken 验证 jwt 有效性
func (svc *UserService) ValidateToken(_ context.Context, req *pb.Token, resp *pb.Token) error {
	if req.Token == "" {
		logrus.Error("token can't empty")
		return errors.New("token can't empty")
	}
	claims, err := svc.Token.Decode(req.Token)
	if err != nil {
		logrus.Errorf("decode token failed: token [%s], error [%s]", req.Token, err.Error())
		return err
	}

	if claims.User.Id <= 0 {
		logrus.Errorf("invalid token: decoded value [%v]", claims.User)
		return errors.New("invalid token")
	}
	resp.Valid = true
	return nil
}

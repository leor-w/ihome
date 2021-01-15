package main

import (
	"github.com/leor-w/ihome/user-service/database"
	"github.com/leor-w/ihome/user-service/handler"
	"github.com/leor-w/ihome/user-service/model"
	pb "github.com/leor-w/ihome/user-service/proto/user"
	repository "github.com/leor-w/ihome/user-service/repo"
	"github.com/leor-w/ihome/user-service/service"
	"github.com/micro/go-micro/v2"
	"github.com/sirupsen/logrus"
)

func main() {
	// init database
	db, err := database.CreateConnection()
	if err != nil {
		logrus.Fatalf("Could not connect to databse: %s", err.Error())
	}

	// migrate model to database
	err = db.AutoMigrate(&model.User{})
	if err != nil {
		logrus.Fatalf("Database migrate failed: %s", err.Error())
	}

	repo := &repository.UserRepository{Db: db}
	tokenSvc := &service.AuthService{
		Conf: service.Config{
			TokenKey: "iHomeUserTokenKeySecret",
			Issuer:   "iHome.user.service",
		},
		Repo: repo,
	}

	// register micro service to go-Micro
	svc := micro.NewService(
		micro.Name("iHome.service.user"),
		micro.Version("v0.0.1"),
	)
	svc.Init()

	// register handler to gRPC
	pb.RegisterUserServiceHandler(
		svc.Server(),
		&handler.UserService{
			Repo:  repo,
			Token: tokenSvc,
		},
	)

	// run go-micro service
	if err := svc.Run(); err != nil {
		logrus.Fatalf("Run service failed: %s", err.Error())
	}
}

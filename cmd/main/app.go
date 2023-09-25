package main

import (
	"context"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net"
	"net/http"
	"os"
	"osipovPetRestApi/internal/config"
	"osipovPetRestApi/internal/user"
	"osipovPetRestApi/internal/user/db"
	"osipovPetRestApi/pkg/client/mongodb"
	"osipovPetRestApi/pkg/logging"
	"path"
	"path/filepath"
	"time"
)

func main() {
	logger := logging.GetLogger()
	logger.Info("create router")
	router := httprouter.New()

	cfg := config.GetConfig()

	cfgMongo := cfg.MongoDB
	mongoDBClient, err := mongodb.NewClient(context.Background(), cfgMongo.Host, cfgMongo.Port, cfgMongo.Username,
		cfgMongo.Password, cfgMongo.Database, cfgMongo.AuthDB)
	if err != nil {
		panic(err)
	}
	storage := db.NewStorage(mongoDBClient, cfgMongo.Collection, logger)

	user1 := user.User{
		Id:           "",
		Name:         "Dmitry",
		PasswordHash: "123456",
		Email:        "os_dimay@mail.ru",
	}
	user1Id, err := storage.Create(context.Background(), user1)
	if err != nil {
		panic(err)
	}
	logger.Info(user1Id)

	user2 := user.User{
		Id:           "",
		Name:         "Dmitry2",
		PasswordHash: "123456",
		Email:        "os_dimay@mail.ru2",
	}
	user2Id, err := storage.Create(context.Background(), user2)
	if err != nil {
		panic(err)
	}
	logger.Info(user2Id)

	user2Found, err := storage.FindOne(context.Background(), user2Id)
	if err != nil {
		panic(err)
	}
	fmt.Println(user2Found)

	user2Found.Email = "newEmail@here.ok"
	err = storage.Update(context.Background(), user2Found)
	if err != nil {
		panic(err)
	}

	err = storage.Delete(context.Background(), user2Id)
	if err != nil {
		panic(err)
	}

	_, err = storage.FindOne(context.Background(), user2Id)
	if err != nil {
		panic(err)
	}

	logger.Info("register user handler")
	handler := user.NewHandler(logger)
	handler.Register(router)

	start(router, cfg)
}

func start(router *httprouter.Router, cfg *config.Config) {
	logger := logging.GetLogger()
	logger.Info("start application")

	var listener net.Listener
	var listenErr error

	if cfg.Listen.Type == "sock" {
		logger.Info("detect app path")
		appDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			logger.Fatal(err)
		}

		logger.Info("create socket")
		socketPath := path.Join(appDir, "app.sock")

		logger.Info("listen unix socket")
		listener, listenErr = net.Listen("unix", socketPath)
		logger.Infof("server is listening unix socket %s", socketPath)
	} else {
		logger.Info("listen tcp")
		listener, listenErr = net.Listen("tcp", fmt.Sprintf("%s:%s", cfg.Listen.BindIp, cfg.Listen.Port))
		logger.Infof("server is listening port %s:%s", cfg.Listen.BindIp, cfg.Listen.Port)
	}

	if listenErr != nil {
		logger.Fatal(listenErr)
	}

	server := &http.Server{
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	logger.Fatal(server.Serve(listener))
}

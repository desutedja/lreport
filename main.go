package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"

	"github.com/desutedja/lreport/internal/config"
	"github.com/desutedja/lreport/internal/handler/bonus"
	"github.com/desutedja/lreport/internal/handler/category"
	"github.com/desutedja/lreport/internal/handler/ping"
	"github.com/desutedja/lreport/internal/handler/transaction"
	"github.com/desutedja/lreport/pkg/database"
	"github.com/desutedja/lreport/pkg/log"
	"github.com/desutedja/lreport/pkg/router"
	"github.com/desutedja/lreport/pkg/token"

	psql "github.com/desutedja/lreport/pkg/database/mysql"

	"github.com/desutedja/lreport/internal/handler/user"

	bonusStore "github.com/desutedja/lreport/internal/repository/mysql/bonus"
	categoryStore "github.com/desutedja/lreport/internal/repository/mysql/category"
	transactionStore "github.com/desutedja/lreport/internal/repository/mysql/transaction"
	userStore "github.com/desutedja/lreport/internal/repository/mysql/user"

	bonusService "github.com/desutedja/lreport/internal/service/bonus"
	categoryService "github.com/desutedja/lreport/internal/service/category"
	transactionService "github.com/desutedja/lreport/internal/service/transaction"
	userService "github.com/desutedja/lreport/internal/service/user"
)

func main() {
	log.Init()
	config.Load()

	routes := setupRouter()
	handler := cors.AllowAll().Handler(routes)

	serverAddr := fmt.Sprintf("0.0.0.0:%d", config.GetHTTPConfig().Port)

	srv := &http.Server{
		Addr:         serverAddr,
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
		IdleTimeout:  60 * time.Second,
		Handler:      handler,
	}

	go func() {
		logrus.Infof("Starting application server on :%s", serverAddr)
		if err := srv.ListenAndServe(); err != nil {
			logrus.WithError(err).Error("failed start HTTP Server")
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), config.GetHTTPConfig().GracefulShutdownTimeout)
	defer cancel()

	srv.Shutdown(ctx)
	logrus.Warn("shutting down")
	os.Exit(0)
}

func setupRouter() *mux.Router {
	handler := setupHandler()
	r := mux.NewRouter()
	baseURL := ""

	// rewrite url if necessary
	if len(baseURL) > 0 && baseURL != "/" {
		r.PathPrefix(baseURL).HandlerFunc(router.RewriteURL(r, baseURL))
	}
	r.StrictSlash(true)
	r.HandleFunc("/health", ping.Ping).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/user/login", handler.user.Login).Methods(http.MethodPost)
	r.HandleFunc("/user/register/bypass", handler.user.CreateUser).Methods(http.MethodOptions, http.MethodPost)

	internal := r.NewRoute().Subrouter()
	internal.Use(handler.tokenStore.MiddlewareJWTAuthorization)
	internal.HandleFunc("/user", handler.user.UserList).Methods(http.MethodOptions, http.MethodGet)
	internal.HandleFunc("/user/login/history", handler.user.LoginHistory).Methods(http.MethodOptions, http.MethodGet)
	internal.HandleFunc("/user/register", handler.user.CreateUser).Methods(http.MethodOptions, http.MethodPost)
	internal.HandleFunc("/user/password/reset", handler.user.ResetPassword).Methods(http.MethodOptions, http.MethodPost)
	internal.HandleFunc("/user/password/change", handler.user.ChangePassword).Methods(http.MethodOptions, http.MethodPut)

	internal.HandleFunc("/category", handler.category.CreateCategory).Methods(http.MethodOptions, http.MethodPost)
	internal.HandleFunc("/category", handler.category.GetCategory).Methods(http.MethodOptions, http.MethodGet)

	internal.HandleFunc("/transaction", handler.transaction.CreateTransaction).Methods(http.MethodOptions, http.MethodPost)
	internal.HandleFunc("/transaction", handler.transaction.GetTransaction).Methods(http.MethodOptions, http.MethodGet)
	internal.HandleFunc("/transaction/statistic", handler.transaction.GetTransactionStatistic).Methods(http.MethodOptions, http.MethodGet)

	internal.HandleFunc("/bonus", handler.bonus.CreateBonus).Methods(http.MethodOptions, http.MethodPost)
	internal.HandleFunc("/bonus", handler.bonus.GetBonus).Methods(http.MethodOptions, http.MethodGet)

	return r
}

type handler struct {
	user        *user.Handler
	tokenStore  *token.TokenGenerator
	category    *category.Handler
	transaction *transaction.Handler
	bonus       *bonus.Handler
}

func setupHandler() *handler {
	ctx := context.Background()
	dbCfg := config.GetDatabaseConfig()

	db, err := psql.Connect(ctx, &database.Config{
		Host:     dbCfg.Host,
		Port:     dbCfg.Port,
		User:     dbCfg.Username,
		Password: dbCfg.Password,
		DBName:   dbCfg.DBName,
	})
	if err != nil {
		logrus.Errorf("host:%s;port:%d;user:%s;password:%s;dbname:%s", dbCfg.Host, dbCfg.Port, dbCfg.Username, dbCfg.Password, dbCfg.DBName)
		logrus.Panic(err)
	}

	tokenStore := token.NewTokenGenerator("S4L7K3Y!", 3600)
	userStore := userStore.NewUserStore(db)
	userService := userService.NewService(userStore, tokenStore)
	userHandler := user.NewHandler(userService)

	categoryStore := categoryStore.NewCategoryStore(db)
	categoryService := categoryService.NewService(categoryStore)
	categoryHandler := category.NewHandler(categoryService)

	transactionStore := transactionStore.NewTransactionStore(db)
	transactionService := transactionService.NewService(transactionStore)
	transactionHandler := transaction.NewHandler(transactionService)

	bonusStore := bonusStore.NewBonusStore(db)
	bonusService := bonusService.NewService(bonusStore)
	bonusHandler := bonus.NewHandler(bonusService)

	return &handler{
		user:        userHandler,
		tokenStore:  tokenStore,
		category:    categoryHandler,
		transaction: transactionHandler,
		bonus:       bonusHandler,
	}
}

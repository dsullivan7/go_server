package app

import (
	"fmt"
	"go_server/internal/auth/auth0"
	"go_server/internal/captcha/twocaptcha"
	"go_server/internal/config"
	goServerRodCrawler "go_server/internal/crawler/rod"
	"go_server/internal/db"
	goServerZapLogger "go_server/internal/logger/zap"
	"go_server/internal/server"
	"go_server/internal/store/gorm"
	"log"
	"net/http"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/launcher/flags"

	"go.uber.org/zap"

	"github.com/go-chi/chi"
)

const callerSkip = 2

func Run() {
	config, configErr := config.NewConfig()

	if configErr != nil {
		log.Fatal(configErr)
	}

	zapConfig := zap.NewProductionConfig()
	zapConfig.OutputPaths = []string{"stdout"}
	zapLogger, errZap := zapConfig.Build(zap.AddCallerSkip(callerSkip))

	if errZap != nil {
		log.Fatal(errZap)
	}

	logger := goServerZapLogger.NewLogger(zapLogger)

	connection, errConnection := db.NewSQLConnection(
		config.DBHost,
		config.DBName,
		config.DBPort,
		config.DBUser,
		config.DBPassword,
		config.DBSSL,
	)

	if errConnection != nil {
		log.Fatal(errConnection)
	}

	db, errDatabase := db.NewGormDB(connection)

	if errDatabase != nil {
		log.Fatal(errDatabase)
	}

	store := gorm.NewStore(db)

	// initialize 2captcha
	captchaKey := config.TwoCaptchaKey
	path, _ := launcher.LookPath()
	u := launcher.New().Set(flags.NoSandbox).Bin(path).MustLaunch()
	browser := rod.New().ControlURL(u)
	captcha := twocaptcha.NewTwoCaptcha(captchaKey, logger)
	crawler := goServerRodCrawler.NewCrawler(browser, captcha)

	auth := auth0.NewAuth(config.Auth0Domain, config.Auth0Audience, logger)
	auth.Init()

	router := chi.NewRouter()
	handler := server.NewChiServer(config, router, store, crawler, auth, logger)

	httpServer := http.Server{
		Addr:    fmt.Sprintf(":%s", config.Port),
		Handler: handler.Init(),
	}

	logger.Info(fmt.Sprintf("started on port: %s", config.Port))
	log.Fatal(httpServer.ListenAndServe())
}

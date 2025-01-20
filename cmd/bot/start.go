package bot

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"quiz-mod/internal/config"
	"quiz-mod/internal/controller"
	"quiz-mod/internal/server"
	"quiz-mod/internal/service/question"
	storage "quiz-mod/internal/storage/postgres/quiz"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	tele "gopkg.in/telebot.v3"
)

func Start(envFile, confName, path string) {
	logrus.Info("starting")
	defer logrus.Info("stopped")

	logLvl := os.Getenv("LOG_LEVEL")

	switch logLvl {
	case "info":
		logrus.SetLevel(logrus.InfoLevel)
	case "warn":
		logrus.SetLevel(logrus.WarnLevel)
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
	case "error":
		logrus.SetLevel(logrus.ErrorLevel)
	case "trace":
		logrus.SetLevel(logrus.TraceLevel)
	case "panic":
		logrus.SetLevel(logrus.PanicLevel)
	case "fatal":
		logrus.SetLevel(logrus.PanicLevel)
	default:
		logrus.SetLevel(logrus.InfoLevel)
	}

	logrus.Infof("log level: %+v", logrus.GetLevel())

	cfg, err := config.LoadConfig(confName, path)
	if err != nil {
		logrus.Fatalf("unable to load config: %v", err)
	}

	err = godotenv.Load()
	if err != nil {
		logrus.Fatal(err)
	}

	token := os.Getenv("TOKEN")
	if len(token) == 0 {
		logrus.Fatalf("bot token is not set")
	}

	timeout := os.Getenv("TIMEOUT")
	if len(timeout) == 0 {
		timeout = "5ms"
	}

	duration, err := time.ParseDuration(timeout)
	if err != nil {
		logrus.Fatalf("failed to parse timeout '%s': %+v", timeout, err)
	}

	channelIDStr := os.Getenv("CHANNEL_ID")
	if len(channelIDStr) == 0 {
		logrus.Fatalf("channel ID is not set")
	}

	channelID, err := strconv.Atoi(channelIDStr)
	if err != nil {
		logrus.Fatalf("cannot parse channel ID %s: %+v", channelIDStr, err)
	}

	// bot
	bot, err := tele.NewBot(tele.Settings{
		Token:     token,
		Poller:    &tele.LongPoller{Timeout: duration},
		ParseMode: "html",
	})
	if err != nil {
		logrus.Fatalf("cannot create a bot: %v", err)
	}

	logrus.Info("successfully created bot")

	dbUser := os.Getenv("POSTGRES_USER")
	if len(dbUser) == 0 {
		logrus.Fatalf("POSTGRES_USER not set")
	}

	dbPass := os.Getenv("POSTGRES_PASSWORD")
	if len(dbPass) == 0 {
		logrus.Fatalf("POSTGRES_PASSWORD not set")
	}

	dbName := os.Getenv("POSTGRES_DB")
	if len(dbName) == 0 {
		logrus.Fatalf("POSTGRES_DB not set")
	}

	dbHost := os.Getenv("POSTGRES_HOST")
	if len(dbHost) == 0 {
		logrus.Fatalf("POSTGRES_HOST not set")
	}

	dbPort := os.Getenv("POSTGRES_PORT")
	if len(dbPort) == 0 {
		logrus.Fatalf("POSTGRES_PORT not set")
	}

	dbAddr := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPass, dbHost, dbPort, dbName)

	logrus.Infof("connecting db: %s", dbAddr)

	storage, err := storage.New(dbAddr)
	if err != nil {
		logrus.Fatalf("failed to connect db: %+v", err)
	}

	questionSrv := question.New(cfg, storage)

	controller := controller.New(bot, channelID, cfg, questionSrv)

	server := server.New(bot, controller)

	logrus.Info("starting server...")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	server.Start(ctx)

	go func() {
		_, msgErr := bot.Send(&tele.Chat{ID: int64(channelID)}, "Бот запущен")
		if msgErr != nil {
			logrus.Errorf("Error while sending message 'Бот запущен': %v\n", msgErr)
		}

		bot.Start()
	}()

	notifyCtx, notify := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	defer notify()

	<-notifyCtx.Done()
	logrus.Info("shutdown")

	var wg sync.WaitGroup

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()

		_, msgErr := bot.Send(&tele.Chat{ID: int64(channelID)}, "Бот выключается")
		if msgErr != nil {
			logrus.Errorf("Error while sending message 'Бот выключается': %v\n", msgErr)
		}
		logrus.Info("gently shutdown")

		bot.Stop()

		storage.Close()

	}(&wg)

	wg.Wait()

	notify()
}

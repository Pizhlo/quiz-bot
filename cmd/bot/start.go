package bot

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"quiz-bot/internal/config"
	"quiz-bot/internal/controller"
	"quiz-bot/internal/server"
	"quiz-bot/internal/service/question"
	"quiz-bot/internal/storage/minio"
	storage "quiz-bot/internal/storage/postgres/quiz"
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
		logrus.Error(err)
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

	minioEndpoint := os.Getenv("MINIO_ENDPOINT")
	if len(dbPort) == 0 {
		logrus.Fatalf("MINIO_ENDPOINT not set")
	}

	minioAccessKey := os.Getenv("MINIO_ACCESS_KEY")
	if len(dbPort) == 0 {
		logrus.Fatalf("MINIO_ACCESS_KEY not set")
	}

	minioSecretAccessKey := os.Getenv("SECRET_ACCESS_KEY")
	if len(dbPort) == 0 {
		logrus.Fatalf("SECRET_ACCESS_KEY not set")
	}

	minioUseSSLStr := os.Getenv("MINIO_USE_SSL")
	if len(dbPort) == 0 {
		logrus.Fatalf("MINIO_USE_SSL not set")
	}

	minioUSeSSL, err := strconv.ParseBool(minioUseSSLStr)
	if err != nil {
		logrus.Fatalf("error parsing string MINIO_USE_SSL '%s': %+v", minioUseSSLStr, err)
	}

	bucket := os.Getenv("MINIO_BUCKET")
	if len(bucket) == 0 {
		logrus.Fatalf("MINIO_BUCKET not set")
	}

	minio, err := minio.New(minioEndpoint, minioAccessKey, minioSecretAccessKey, minioUSeSSL, bucket)
	if err != nil {
		logrus.Fatalf("failed to connect minio: %+v", err)
	}

	endRoundPic := os.Getenv("END_ROUND_PICTURE")
	if len(endRoundPic) == 0 {
		logrus.Fatalf("END_ROUND_PICTURE not set")
	}

	questionSrv := question.New(cfg, storage, minio, endRoundPic)

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

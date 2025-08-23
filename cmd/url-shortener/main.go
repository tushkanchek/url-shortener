package main

import (
	"back/back/urlShortner/internal/config"
	"back/back/urlShortner/internal/config/http-server/middleware/logger"
	"back/back/urlShortner/internal/config/lib/logger/handlers/slogpretty"
	"back/back/urlShortner/internal/config/lib/logger/sl"
	"back/back/urlShortner/internal/config/storage/postgres"
	"log/slog"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	envLocal = "local"
	envDev = "dev"
	envProd = "prod"
)


func main(){
	 
	cfg:= config.MustLoad()

	log:= setupLogger(cfg.Env)

	

	log.Info("starting url-shortener", slog.String("env",cfg.Env))  //first start shows our environment
	log.Debug("debug messages are enabled")
	log.Error("error message are enabled")

	

	storage, err:= postgres.New(cfg.DBConfig)
	if err!=nil{
		log.Error("failed to init storage",sl.Err(err))
		os.Exit(1)
	}



	_ = storage



	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(logger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	// TODO: run server
}

func setupLogger(env string) *slog.Logger{
	var log *slog.Logger
	switch env{
	case envLocal:
		log = setupPrettySlog()
	case envDev:
		log  = slog.New(
			slog.NewJSONHandler(os.Stdout,&slog.HandlerOptions{Level: slog.LevelDebug}),
		)	
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)	
	}

	return log
}


func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
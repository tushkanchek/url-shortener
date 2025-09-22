package main
import (
	"urlShortner/internal/config"
	"urlShortner/internal/http-server/handlers/url/redirect"
	"urlShortner/internal/http-server/handlers/url/save"
	"urlShortner/internal/http-server/middleware/logger"
	"urlShortner/internal/lib/logger/handlers/slogpretty"
	"urlShortner/internal/lib/logger/sl"
	"urlShortner/internal/storage/postgres"
	"log/slog"
	"net/http"
	"os"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"urlShortner/internal/http-server/handlers/url/delete"
	
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)



func main() {

	initConfig()

	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info(
		"starting url-shortener", 
		slog.String("env", cfg.Env),
		slog.String("version", "123"),
	) //first start shows our environment
	log.Debug("debug messages are enabled")

	storage, err := postgres.New(cfg.DBConfig)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
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


	router.Route("/url", func(r chi.Router) {
		r.Use(middleware.BasicAuth("url-shortener", map[string]string{
			cfg.HTTPServer.User: cfg.HTTPServer.Password, 
		}))

		r.Post("/", save.New(log, storage))
		r.Delete("/{alias}", delete.New(log,storage))
	})

	
	router.Get("/{alias}", redirect.New(log, storage))
	
	
	log.Info("starting server", slog.String("adress", cfg.Address))

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}

	log.Error("server stopped")
	// TODO: run server
}

//
func initConfig(){
	if err:=godotenv.Load(); err!=nil{
		os.Exit(1)
	}
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = setupPrettySlog()
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
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

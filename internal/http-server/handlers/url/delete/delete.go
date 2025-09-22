package delete

import (
	resp "urlShortner/internal/lib/api/response"
	"urlShortner/internal/lib/logger/sl"
	"urlShortner/internal/storage"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

//go:generate go run github.com/vektra/mockery/v2@v2.53.5 --name=URLDeleter
type URLDeleter interface{
	DeleteURL(alias string) error
}


func New(log *slog.Logger, urlDeleter URLDeleter) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){
		const op = "handlers.url.delete.New"

		log = log.With(
			slog.String("op",op),
			slog.String("request_id", middleware.GetReqID(r.Context())),

		)

		alias:= chi.URLParam(r, "alias")
		if alias==""{
			log.Info("alias is empty")

			render.JSON(w,r,resp.Error("invalid request"))

			return 
		}

		err:=urlDeleter.DeleteURL(alias)
		if errors.Is(err, storage.ErrURLNotFound){
			log.Info("url not found", "alias",alias)

			render.JSON(w, r, resp.Error("not found"))

			return 
		}
		if err!=nil{
			log.Info("failed to get url", sl.Err(err))

			render.JSON(w, r, resp.Error("internal error"))

			return 
		}
		
		log.Info("deleted url", slog.String("alias",alias))

		

	}

}
	


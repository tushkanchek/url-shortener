package save

import (
	resp "back/back/urlShortner/internal/config/lib/api/response"
	"back/back/urlShortner/internal/config/lib/logger/sl"
	"back/back/urlShortner/internal/config/lib/random"

	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type Request struct{
	URL string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct{
	resp.Response
	Alias string `json:"alias,omitempty"`
}


//TODO: move to cfg
const aliasLength = 6

type URLSaver interface{
	SaveURL(urlToSave string,alias string) (int64, error)
}

func New(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){
		const op = "handlers.url.save.New"

		log = log.With(
			slog.String("op",op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err:= render.DecodeJSON(r.Body, &req)

		if err!=nil{
			log.Error("failed to decode request body", sl.Err(err))
			
			render.JSON(w,r,resp.Error("failed to decode request"))

			return 
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err:=validator.New().Struct(req);err!=nil{
			validateErr := err.(validator.ValidationErrors)

			log.Error("Invalid request", sl.Err(err))

			render.JSON(w,r,resp.ValidationError(validateErr))

			return 
		}

		alias:= req.Alias

		if alias == ""{
			alias = random.NewRandomString(aliasLength)
		}
		
	}
}
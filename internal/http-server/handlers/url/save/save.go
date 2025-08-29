package save

import (
	resp "back/back/urlShortner/internal/lib/api/response"
	"back/back/urlShortner/internal/lib/logger/sl"
	"back/back/urlShortner/internal/lib/random"
	"back/back/urlShortner/internal/storage"
	"errors"

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


//go:generate go run github.com/vektra/mockery/v2@v2.53.5 --name=URLSaver
type URLSaver interface{
	SaveURL(urlToSave string,alias string) error
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
		
		err = urlSaver.SaveURL(req.URL,alias)
		if errors.Is(err,storage.ErrURLExists){
			log.Info("url already exists", slog.String("url",req.URL))

			render.JSON(w, r, resp.Error("url already exists"))

			return 	
		}

		if err != nil{
			log.Error("failed to add url", sl.Err(err))

			render.JSON(w,r, resp.Error("failed to add url"))

			return 
		}

		log.Info("url added")

		responceOK(w,r,alias)
		
	}
}


func responceOK(w http.ResponseWriter, r *http.Request, alias string){
	render.JSON(w, r, Response{
		Response: resp.OK(),
		Alias: alias,
	})
}
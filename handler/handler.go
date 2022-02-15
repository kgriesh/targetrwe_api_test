package handler

import (
	"net/http"
	"net/url"
	"strconv"

	"github.com/targetrw-api/db"
	database "github.com/targetrw-api/db"
	model "github.com/targetrw-api/model"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
)

func NewHandler(db db.DbService) http.Handler {
	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))
	router.Route("/products", func(r chi.Router) {
		r.Get("/", getAllProducts(db))
		r.Get("/{itemId}", GetProduct(db))
		r.Post("/", CreateProduct(db))
	})
	return router
}

func CreateProduct(db db.DbService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		p := &model.Product{}
		err := render.Bind(r, p)
		if err != nil {
			render.Render(w, r, BadRequest(err))
			return
		}
		err = db.AddProduct(p)
		if err != nil {
			render.Render(w, r, ServerError(err))
			return
		}
		err = render.Render(w, r, p)
		if err != nil {
			render.Render(w, r, ServerError(err))
			return
		}
	}
}

func getAllProducts(db db.DbService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		pl, err := db.GetAllProducts()
		if err != nil {
			render.Render(w, r, ServerError(err))
			return
		}
		err = render.Render(w, r, pl)
		if err != nil {
			render.Render(w, r, ServerError(err))
			return
		}
	}
}

func GetProduct(db db.DbService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		prodId, _ := url.PathUnescape(chi.URLParam(r, "itemId"))
		id, err := strconv.Atoi(prodId)
		if err != nil {
			render.Render(w, r, BadRequest(err))
			return
		}
		product, err := db.GetProductById(id)
		if err != nil {
			if err == database.ErrNoMatch {
				render.Render(w, r, BadRequest(err))
			} else {
				render.Render(w, r, ServerError(err))
			}
			return
		}
		err = render.Render(w, r, &product)
		if err != nil {
			render.Render(w, r, ServerError(err))
			return
		}
	}
}

package main

import (
	"github.com/gorilla/mux"
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"fmt"
	"net/http"
	"strconv"
	"encoding/json"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

func (a *App) Initialize(user, password, dbname string) {
	connectionString :=
		fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", user, password, dbname)

	var err error
	a.DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

func (a *App)initializeRoutes()  {
	a.Router.HandleFunc("/product/{id:[0-9]+}", a.getProduct).Methods("GET")
	a.Router.HandleFunc("/product", a.createProduct).Methods("POST")
	a.Router.HandleFunc("/products", a.getProducts).Methods("GET")
	a.Router.HandleFunc("/product/{id:[0-9]+}", a.updateProduct).Methods("PUT")
	a.Router.HandleFunc("/product/{id:[0-9]+}", a.deleteProduct).Methods("DELETE")

}

func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}

func (a *App) getProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, error := strconv.Atoi(vars["id"])
	if error != nil {
		respondWithError(w, http.StatusBadRequest, "id doesnt exist")
	}

	p := product{ID: id}

	if err := p.getProduct(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Product not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return

	}
	respondWithJSON(w, http.StatusOK, p)
}
func (a *App) getProducts(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var count, start int
	count, _ = strconv.Atoi(vars["count"])
	start, _ = strconv.Atoi(vars["start"])
	if count < 1 || count > 10 {
		count = 10
	}
	if start < 0 {
		start = 0
	}

	products, err := getProducts(a.DB, start, count)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, products)
}

func (a *App) createProduct(w http.ResponseWriter, r *http.Request) {
	var prod product
	err := json.NewDecoder(r.Body).Decode(&prod)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	if err := prod.createProduct(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusCreated, prod)
}

func (a *App) updateProduct(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id, error := strconv.Atoi(vars["id"])
	if error != nil {
		respondWithError(w, http.StatusBadRequest, "id doesnt exist")
	}

	var prod product
	err := json.NewDecoder(r.Body).Decode(&prod)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	prod.ID = id
	defer r.Body.Close()

	if err := prod.updateProduct(a.DB); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, prod)

}

func (a *App) deleteProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
	}

	var prod product
	prod.ID = id
	if err := prod.deleteProduct(a.DB); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
	}
	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})

}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

package router

import (
	"encoding/json"
	"net/http"

	"../dao"
	"../models"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"gopkg.in/mgo.v2/bson"
)

type AccountService struct {
	Dao dao.IAccountDao
}

var (
	getReqProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "get_request_total",
		Help: "The total number of 'get' request processed events",
	})

	postReqProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "post_request_total",
		Help: "The total number of 'post' request processed events",
	})

	deleteReqProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "delete_request_total",
		Help: "The total number of 'delete' request processed events",
	})

	putReqProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "put_request_total",
		Help: "The total number of 'put' request processed events",
	})
)

func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJson(w, code, map[string]string{"error": msg})
}

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func (accountService *AccountService) GetByID(w http.ResponseWriter, r *http.Request) {
	defer getReqProcessed.Inc()
	params := mux.Vars(r)
	account, err := accountService.Dao.GetByID(params["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid account ID")
		return
	}
	respondWithJson(w, http.StatusOK, account)
}

func (accountService *AccountService) CreateAccount(w http.ResponseWriter, r *http.Request) {
	defer postReqProcessed.Inc()
	defer r.Body.Close()
	var account models.Account
	if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	account.ID = bson.NewObjectId()
	if err := accountService.Dao.CreateAccount(account); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusCreated, account)
}

func (accountService *AccountService) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	defer postReqProcessed.Inc()
	defer r.Body.Close()
	var transaction models.Transaction
	if err := json.NewDecoder(r.Body).Decode(&transaction); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// validate acocunt id
	_, err := accountService.Dao.GetByID(transaction.AccountID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid account ID")
		return
	}

	transaction.ID = bson.NewObjectId()
	if err := accountService.Dao.CreateTransaction(transaction); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJson(w, http.StatusCreated, transaction)
}

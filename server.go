package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	dbconfig "./config"
	"./dao"
	"./router"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var details = dbconfig.Config{}
var accountDao = &dao.AccountDao{}
var accountRouter = router.AccountService{}

func init() {
	details.Read()
	accountDao.Server = details.Server
	accountDao.Database = details.Database
	accountDao.Connect()
	accountRouter.Dao = accountDao
}

func main() {
	r := mux.NewRouter()
	r.Handle("/metrics", promhttp.Handler())
	r.HandleFunc("/api/v1/accounts", accountRouter.CreateAccount).Methods("POST")
	r.HandleFunc("/api/v1/accounts/{id}", accountRouter.GetByID).Methods("GET")
	r.HandleFunc("/api/v1/transactions", accountRouter.CreateTransaction).Methods("POST")

	var _, portErr = strconv.ParseInt(os.Getenv("PORT"), 10, 32)
	if portErr != nil {
		log.Fatal("Missing port number environment variable value")
	}
	if _, err := strconv.Atoi(os.Getenv("PORT")); err != nil {
		log.Fatal("Bad format port number environment variable value")
	}
	var readTimeout, readTimeoutErr = strconv.ParseInt(os.Getenv("READ_TIMEOUT"), 10, 64)
	if readTimeoutErr != nil {
		log.Fatal("Missing read timeout environment variable value READ_TIMEOUT")
	}

	var writeTimeout, writeTimeoutErr = strconv.ParseInt(os.Getenv("WRITE_TIMEOUT"), 10, 64)
	if writeTimeoutErr != nil {
		log.Fatal("Missing write timeout environment variable value WRITE_TIMEOUT")
	}

	fmt.Println("Server running in port:", os.Getenv("PORT"))

	srv := &http.Server{
		Addr:         ":" + os.Getenv("PORT"),
		Handler:      r,
		ReadTimeout:  time.Second * time.Duration(readTimeout),
		WriteTimeout: time.Second * time.Duration(writeTimeout),
	}

	log.Fatal(srv.ListenAndServe())
}

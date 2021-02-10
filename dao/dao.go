package dao

import (
	"errors"
	"log"

	"../models"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type AccountDao struct {
	Server   string
	Database string
}

var db *mgo.Database
var accountCollection *mgo.Collection
var transCollection *mgo.Collection

const (
	ACCOUNT_COLLECTION     = "accounts"
	TRANSACTION_COLLECTION = "transactions"
)

type IAccountDao interface {
	GetByID(id string) (models.Account, error)
	CreateAccount(account models.Account) error
	CreateTransaction(account models.Transaction) error
}

func (m *AccountDao) Connect() {
	session, err := mgo.Dial(m.Server)
	if err != nil {
		log.Fatal(err)
	}
	db = session.DB(m.Database)
	accountCollection = db.C(ACCOUNT_COLLECTION)
	transCollection = db.C(TRANSACTION_COLLECTION)
}

func (m *AccountDao) GetByID(id string) (models.Account, error) {
	var account models.Account
	accoErr := accountCollection.FindId(bson.ObjectIdHex(id)).One(&account)
	var transactions []models.Transaction
	transErr := transCollection.Find(bson.M{"account_id": id}).All(&transactions)

	var err error
	if accoErr != nil || transErr != nil {
		err = errors.New("Error retrieving account and transactions")
	}
	account.Transactions = transactions
	return account, err
}

func (m *AccountDao) CreateAccount(account models.Account) error {
	err := accountCollection.Insert(&account)
	return err
}

func (m *AccountDao) CreateTransaction(transaction models.Transaction) error {
	err := transCollection.Insert(&transaction)
	return err
}

package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Account struct {
	ID             bson.ObjectId `bson:"_id"`
	DocumentNumber string        `bson:"document_number" json:"document_number"`
	Transactions   []Transaction
}

type Transaction struct {
	ID              bson.ObjectId `bson:"_id"`
	AccountID       string        `bson:"account_id" json:"account_id"`
	OperationTypeID int16         `bson:"operation_type_id" json:"operation_type_id"`
	Amount          float32       `bson:"amount" json:"amount"`
	EventDate       time.Time     `bson:"eventDate"`
}

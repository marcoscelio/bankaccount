package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"../models"
)

type MockAccountDao struct {
}

func (m *MockAccountDao) GetByID(id string) (models.Account, error) {
	var account = models.Account{Name: "Marcos", Email: "marcos@email.com", Active: true}
	return account, nil
}

func (m *MockAccountDao) CreateAccount(account models.Account) error {
	return nil
}

func (m *MockAccountDao) CreateTransaction(account models.Account) error {
	return nil
}

func TestAccontRouterGetByID(t *testing.T) {

	req, err := http.NewRequest("GET", "/account", nil)
	if err != nil {
		t.Fatal(err)
	}

	var accountRouter = AccountService{Dao: &MockAccountDao{}}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(accountRouter.GetByID)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `{"id":"","name":"Marcos","email":"marcos@email.com","active":true}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

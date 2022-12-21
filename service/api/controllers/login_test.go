package controllers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/lucaronca/wasa-homework/service/api/models"
	"github.com/lucaronca/wasa-homework/service/api/repositories"
	"github.com/lucaronca/wasa-homework/service/api/reqcontext"
	"github.com/lucaronca/wasa-homework/service/api/services"
)

type authRepositoryMock struct{}

func (am *authRepositoryMock) GetToken(relations ...repositories.Relation) (string, error) {
	return "mock_token", nil
}
func (am *authRepositoryMock) SetToken(userId int, token string) error {
	return nil
}
func (am *authRepositoryMock) WithTokens() repositories.Relation {
	return nil
}
func (am *authRepositoryMock) FilterByToken(token string) repositories.Relation {
	return nil
}

type usersRepositoryMock struct {
	hasUser bool
}

func (um *usersRepositoryMock) GetUserById(id int) (*models.BaseUser, error) {
	return nil, nil
}
func (um *usersRepositoryMock) GetUser(relations ...repositories.Relation) (*models.BaseUser, error) {
	if um.hasUser {
		return &models.BaseUser{Id: 1, Username: "Mario"}, nil
	}
	return nil, nil
}
func (um *usersRepositoryMock) GetFullUser(relations ...repositories.Relation) (*models.FullUser, error) {
	return nil, nil
}
func (um *usersRepositoryMock) GetUsers(relations ...repositories.Relation) (*[]models.BaseUser, error) {
	return nil, nil
}
func (um *usersRepositoryMock) CreateUser(user *models.BaseUser) (userId int, err error) {
	return 0, nil
}
func (um *usersRepositoryMock) UpdateUser(user *models.BaseUser) (err error) {
	return nil
}
func (um *usersRepositoryMock) WithUsers() repositories.Relation {
	return nil
}
func (um *usersRepositoryMock) FilterByUserId(userId int) repositories.Relation {
	return nil
}
func (um *usersRepositoryMock) FilterByUsername(username string, strict bool) repositories.Relation {
	return nil
}

func TestDoLogin_CreateUser(t *testing.T) {
	var jsonStr = []byte(`{"name": "Mario"}`)
	req, err := http.NewRequest(http.MethodPost, "/session", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}

	authService := services.NewAuthService(&authRepositoryMock{}, &usersRepositoryMock{})
	lci := NewLoginController(authService)
	lc, _ := lci.(*loginController)

	res := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lc.DoLogin(w, r, httprouter.Params{}, reqcontext.RequestContext{})
	})

	handler.ServeHTTP(res, req)

	if http.StatusCreated != res.Code {
		t.Error("expected", http.StatusCreated, "got:", res.Code)
	}

	if !strings.Contains(res.Body.String(), "\"identifier\":") {
		t.Error("expected JSON response with \"identifier\" got:", res.Body.String())
	}
}

func TestDoLogin_UseExistingUser(t *testing.T) {
	var jsonStr = []byte(`{"name": "Mario"}`)
	req, err := http.NewRequest(http.MethodPost, "/session", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}

	authService := services.NewAuthService(&authRepositoryMock{}, &usersRepositoryMock{true})
	lci := NewLoginController(authService)
	lc, _ := lci.(*loginController)

	res := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lc.DoLogin(w, r, httprouter.Params{}, reqcontext.RequestContext{})
	})

	handler.ServeHTTP(res, req)

	if http.StatusOK != res.Code {
		t.Error("expected", http.StatusOK, "got:", res.Code)
	}

	if !strings.Contains(res.Body.String(), "\"identifier\":") {
		t.Error("expected JSON response with \"identifier\" got:", res.Body.String())
	}
}

func TestDoLogin_PayloadError(t *testing.T) {
	var jsonStr = []byte(`{"not": "valid"}`)
	req, err := http.NewRequest(http.MethodPost, "/session", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}

	authService := services.NewAuthService(&authRepositoryMock{}, &usersRepositoryMock{})
	lci := NewLoginController(authService)
	lc, _ := lci.(*loginController)

	res := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lc.DoLogin(w, r, httprouter.Params{}, reqcontext.RequestContext{})
	})

	handler.ServeHTTP(res, req)

	if http.StatusBadRequest != res.Code {
		t.Error("expected", http.StatusOK, "got:", res.Code)
	}

	if res.Body.String() != "Payload not valid" {
		t.Error("expected \"Payload not valid\" got:", res.Body.String())
	}
}

func TestDoLogin_NameNotValid(t *testing.T) {
	var jsonStr = []byte(`{"name": "M"}`)
	req, err := http.NewRequest(http.MethodPost, "/session", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}

	authService := services.NewAuthService(&authRepositoryMock{}, &usersRepositoryMock{})
	lci := NewLoginController(authService)
	lc, _ := lci.(*loginController)

	res := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lc.DoLogin(w, r, httprouter.Params{}, reqcontext.RequestContext{})
	})

	handler.ServeHTTP(res, req)

	if http.StatusBadRequest != res.Code {
		t.Error("expected", http.StatusOK, "got:", res.Code)
	}

	if res.Body.String() != "Name should be at least 3 characters long" {
		t.Error("expected \"Name should be at least 3 characters long\" got:", res.Body.String())
	}
}

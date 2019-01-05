package webserver

import (
	"TicketSystem/config"
	"TicketSystem/utils"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path"
	"strings"
	"testing"
)

func setup() {
	config.DataPath = path.Join("..", "datatest")
	config.TemplatePath = path.Join("..", "templates")
	Setup()
}

func teardown() {
	err := os.RemoveAll(config.DataPath)
	if err != nil {
		log.Println(err)
	}
}

func TestServeNewTicket(t *testing.T) {
	setup()
	defer teardown()

	req := httptest.NewRequest(http.MethodPost, "/tickets/new", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ServeNewTicket)

	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestServeIndex(t *testing.T) {
	setup()
	defer teardown()

	req := httptest.NewRequest(http.MethodPost, "/", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ServeIndex)

	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestServeErrorPage(t *testing.T) {
	setup()
	defer teardown()

	req := httptest.NewRequest(http.MethodPost, "/error/1", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ServeErrorPage)
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "/error/1", req.URL.Path)

	req = httptest.NewRequest(http.MethodPost, "/error/1000", nil)
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(ServeErrorPage)
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusMovedPermanently, rr.Code)
	resultURL, err := rr.Result().Location()
	assert.Nil(t, err)
	assert.Equal(t, "/error/0", resultURL.Path)

	req = httptest.NewRequest(http.MethodPost, "/error", nil)
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(ServeErrorPage)
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusMovedPermanently, rr.Code)
	resultURL, err = rr.Result().Location()
	assert.Nil(t, err)
	assert.Equal(t, "/error/0", resultURL.Path)

	req = httptest.NewRequest(http.MethodPost, "/error/1test", nil)
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(ServeErrorPage)
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusMovedPermanently, rr.Code)
	resultURL, err = rr.Result().Location()
	assert.Nil(t, err)
	assert.Equal(t, "/error/0", resultURL.Path)
}

func TestServeSignOut(t *testing.T) {
	setup()
	defer teardown()

	req := httptest.NewRequest(http.MethodPost, "/signOut", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ServeSignOut)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusMovedPermanently, rr.Code)
	resultURL, err := rr.Result().Location()
	assert.Nil(t, err)
	assert.Equal(t, "/", resultURL.Path)

	for _, cookie := range rr.Result().Cookies() {
		if cookie.Name == "session-id" || cookie.Name == "requested-url-while-not-authenticated" {
			assert.Equal(t, -1, cookie.MaxAge)
		}
	}
}

func TestServeUserRegistrationShowTemplate(t *testing.T) {
	setup()
	defer teardown()

	form := url.Values{}
	form.Add("username", "Test123")
	form.Add("password1", "")
	form.Add("password2", "123")

	req := httptest.NewRequest(http.MethodPost, "/signUp", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	req.Form = form

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ServeUserRegistration)
	handler.ServeHTTP(rr, req)
	assert.Equal(t, rr.Code, http.StatusOK)
}

func TestServeUserRegistrationPasswordsDontMatch(t *testing.T) {
	setup()
	defer teardown()

	form := url.Values{}
	form.Add("username", "Test123")
	form.Add("password1", "123456")
	form.Add("password2", "12345")

	req := httptest.NewRequest(http.MethodPost, "/signUp", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	req.Form = form

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ServeUserRegistration)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusFound)
	resultURL, err := rr.Result().Location()
	assert.Nil(t, err)
	assert.Equal(t, utils.ErrorInvalidInputs.ErrorPageUrl(), resultURL.Path)
}

func TestServeUserRegistrationInvalidUsername(t *testing.T) {
	setup()
	defer teardown()

	form := url.Values{}
	form.Add("username", "Tes")
	form.Add("password1", "12345")
	form.Add("password2", "12345")

	req := httptest.NewRequest(http.MethodPost, "/signUp", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	req.Form = form

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ServeUserRegistration)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusFound)
	resultURL, err := rr.Result().Location()
	assert.Nil(t, err)
	assert.Equal(t, utils.ErrorInvalidInputs.ErrorPageUrl(), resultURL.Path)
}

func TestServeUserRegistrationInvalidPassword(t *testing.T) {
	setup()
	defer teardown()

	form := url.Values{}
	form.Add("username", "Test123")
	form.Add("password1", "123")
	form.Add("password2", "123")

	req := httptest.NewRequest(http.MethodPost, "/signUp", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	req.Form = form

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ServeUserRegistration)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusFound)
	resultURL, err := rr.Result().Location()
	assert.Nil(t, err)
	assert.Equal(t, utils.ErrorInvalidInputs.ErrorPageUrl(), resultURL.Path)
}

func TestServeUserRegistrationSuccess(t *testing.T) {
	setup()
	defer teardown()

	form := url.Values{}
	form.Add("username", "Test123")
	form.Add("password1", "Aa!123456")
	form.Add("password2", "Aa!123456")

	req := httptest.NewRequest(http.MethodPost, "/signUp", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	req.Form = form

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ServeUserRegistration)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusMovedPermanently)
	resultURL, err := rr.Result().Location()
	assert.Nil(t, err)
	assert.Equal(t, "/", resultURL.Path)
}

func TestServeSignInShowTemplate(t *testing.T) {
	setup()
	defer teardown()

	form := url.Values{}

	req := httptest.NewRequest(http.MethodPost, "/signIn", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	req.Form = form

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ServeSignIn)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusOK)
}

func createUser(username, password string) {
	form := url.Values{}
	form.Add("username", username)
	form.Add("password1", password)
	form.Add("password2", password)

	req := httptest.NewRequest(http.MethodPost, "/signUp", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	req.Form = form

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ServeUserRegistration)
	handler.ServeHTTP(rr, req)
}

func TestServeSignInInvalidCredentials(t *testing.T) {
	setup()
	defer teardown()

	form := url.Values{}
	form.Add("username", "Tes23")
	form.Add("password", "A3456")

	req := httptest.NewRequest(http.MethodPost, "/signIn", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	req.Form = form

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ServeSignIn)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusFound)
	resultURL, err := rr.Result().Location()
	assert.Nil(t, err)
	assert.Equal(t, utils.ErrorUserLogin.ErrorPageUrl(), resultURL.Path)
}

func TestServeSignInSuccessWithHomeRedirect(t *testing.T) {
	setup()
	defer teardown()

	createUser("Test123", "Aa!123456")

	form := url.Values{}
	form.Add("username", "Test123")
	form.Add("password", "Aa!123456")

	req := httptest.NewRequest(http.MethodPost, "/signIn", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	req.Form = form

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ServeSignIn)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusFound)
	resultURL, err := rr.Result().Location()
	assert.Nil(t, err)
	assert.Equal(t, "/", resultURL.Path)
}

func TestServeSignInSuccessWithRequestedURLRedirect(t *testing.T) {
	setup()
	defer teardown()

	createUser("Test123", "Aa!123456")
	requestedURL := "/ticket/new"

	form := url.Values{}
	form.Add("username", "Test123")
	form.Add("password", "Aa!123456")

	req := httptest.NewRequest(http.MethodPost, "/signIn", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	req.Form = form
	req.AddCookie(&http.Cookie{
		Name:     "requested-url-while-not-authenticated",
		Value:    requestedURL,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   60,
	})

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ServeSignIn)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusFound)
	resultURL, err := rr.Result().Location()
	assert.Nil(t, err)
	assert.Equal(t, requestedURL, resultURL.Path)
}

/*
func TestServeTicketCreationInvalidInput(t *testing.T) {
	setup()
	defer teardown()

	createUser("Test123", "Aa!123456")
	requestedURL := "/ticket/new"

	form := url.Values{}
	form.Add("email", "mustermann@gmail.com")
	form.Add("subject", "PC Issue")
	form.Add("message", "I have issues with my pc...")

	req := httptest.NewRequest(http.MethodPost, "/signIn", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	req.Form = form

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ServeTicketCreation)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusFound)
	resultURL, err := rr.Result().Location()
	assert.Nil(t, err)
	assert.Equal(t, requestedURL, resultURL.Path)
}*/

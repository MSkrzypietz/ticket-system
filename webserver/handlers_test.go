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
	"strconv"
	"strings"
	"testing"
)

// TODO: In jeder Testdatei anstelle assert.NotNil(t, err) -> assert.NoError(t, err)

func setup() {
	config.DataPath = "datatest"
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
	assert.Equal(t, utils.ErrorInvalidInputs.ErrorPageURL(), resultURL.Path)
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
	assert.Equal(t, utils.ErrorInvalidInputs.ErrorPageURL(), resultURL.Path)
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
	assert.Equal(t, utils.ErrorInvalidInputs.ErrorPageURL(), resultURL.Path)
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

func TestServeAuthenticationShowTemplate(t *testing.T) {
	setup()
	defer teardown()

	form := url.Values{}

	req := httptest.NewRequest(http.MethodPost, "/signIn", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	req.Form = form

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ServeAuthentication)
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

func TestServeAuthenticationInvalidCredentials(t *testing.T) {
	setup()
	defer teardown()

	form := url.Values{}
	form.Add("username", "Tes23")
	form.Add("password", "A3456")

	req := httptest.NewRequest(http.MethodPost, "/signIn", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	req.Form = form

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ServeAuthentication)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusFound)
	resultURL, err := rr.Result().Location()
	assert.Nil(t, err)
	assert.Equal(t, utils.ErrorUserLogin.ErrorPageURL(), resultURL.Path)
}

func TestServeAuthenticationSuccessWithHomeRedirect(t *testing.T) {
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
	handler := http.HandlerFunc(ServeAuthentication)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusFound)
	resultURL, err := rr.Result().Location()
	assert.Nil(t, err)
	assert.Equal(t, "/", resultURL.Path)
}

func TestServeAuthenticationSuccessWithRequestedURLRedirect(t *testing.T) {
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
	handler := http.HandlerFunc(ServeAuthentication)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusFound)
	resultURL, err := rr.Result().Location()
	assert.Nil(t, err)
	assert.Equal(t, requestedURL, resultURL.Path)
}

func TestServeTicketCreationInvalidInputs(t *testing.T) {
	setup()
	defer teardown()

	form := url.Values{}
	form.Add("email", "mustermann@gmail.com")
	form.Add("subject", "PC Issue")

	req := httptest.NewRequest(http.MethodPost, "/tickets/new", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	req.Form = form

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ServeTicketCreation)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusFound)
	resultURL, err := rr.Result().Location()
	assert.Nil(t, err)
	assert.Equal(t, utils.ErrorInvalidInputs.ErrorPageURL(), resultURL.Path)
}

func TestServeTicketCreationSuccess(t *testing.T) {
	setup()
	defer teardown()

	form := url.Values{}
	form.Add("email", "mustermann@gmail.com")
	form.Add("subject", "PC Issue")
	form.Add("message", "I have issues with my pc...")

	req := httptest.NewRequest(http.MethodPost, "/tickets/new", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	req.Form = form

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ServeTicketCreation)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusMovedPermanently)
	resultURL, err := rr.Result().Location()
	assert.Nil(t, err)
	assert.Equal(t, "/", resultURL.Path)
}

func TestServeAddCommentUnauthorized(t *testing.T) {
	setup()
	defer teardown()

	form := url.Values{}
	form.Add("comment", "Test Comment")

	req := httptest.NewRequest(http.MethodPost, "/addComment", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	req.Form = form

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ServeAddComment)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusFound)
	resultURL, err := rr.Result().Location()
	assert.Nil(t, err)
	assert.Equal(t, utils.ErrorUnauthorized.ErrorPageURL(), resultURL.Path)
}

func TestServeAddCommentInvalidInput(t *testing.T) {
	setup()
	defer teardown()

	form := url.Values{}
	form.Add("comment", "")

	req := httptest.NewRequest(http.MethodPost, "/addComment", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	req.Form = form

	uuid := utils.CreateUUID(64)
	req.AddCookie(&http.Cookie{
		Name:     "session-id",
		Value:    uuid,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   60 * 60,
	})

	rr := httptest.NewRecorder()
	createUser("Test123", "Aa!123456")
	assert.Nil(t, loginUser(rr, "Test123", "Aa!123456", uuid))

	handler := http.HandlerFunc(ServeAddComment)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusFound)
	resultURL, err := rr.Result().Location()
	assert.Nil(t, err)
	assert.Equal(t, utils.ErrorInvalidInputs.ErrorPageURL(), resultURL.Path)
}

func loginUser(rr *httptest.ResponseRecorder, username, password, uuid string) error {
	err := utils.LoginUser(username, password, uuid)
	CreateSessionCookie(rr, uuid)
	return err
}

func TestServeAddCommentInvalidURL(t *testing.T) {
	setup()
	defer teardown()

	form := url.Values{}
	form.Add("comment", "My comment")

	req := httptest.NewRequest(http.MethodPost, "/addComment", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	req.Header.Set("Referer", "/tickets/MyTicket")
	req.Form = form

	uuid := utils.CreateUUID(64)
	req.AddCookie(&http.Cookie{
		Name:     "session-id",
		Value:    uuid,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   60 * 60,
	})

	rr := httptest.NewRecorder()
	createUser("Test123", "Aa!123456")
	assert.Nil(t, loginUser(rr, "Test123", "Aa!123456", uuid))

	handler := http.HandlerFunc(ServeAddComment)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusFound)
	resultURL, err := rr.Result().Location()
	assert.Nil(t, err)
	assert.Equal(t, utils.ErrorURLParsing.ErrorPageURL(), resultURL.Path)
}

func TestServeAddCommentInvalidTicketID(t *testing.T) {
	setup()
	defer teardown()

	form := url.Values{}
	form.Add("comment", "My comment")

	req := httptest.NewRequest(http.MethodPost, "/addComment", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	req.Header.Set("Referer", "/tickets/1337")
	req.Form = form

	uuid := utils.CreateUUID(64)
	req.AddCookie(&http.Cookie{
		Name:     "session-id",
		Value:    uuid,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   60 * 60,
	})

	rr := httptest.NewRecorder()
	createUser("Test123", "Aa!123456")
	assert.Nil(t, loginUser(rr, "Test123", "Aa!123456", uuid))

	handler := http.HandlerFunc(ServeAddComment)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusFound)
	resultURL, err := rr.Result().Location()
	assert.Nil(t, err)
	assert.Equal(t, utils.ErrorInvalidTicketID.ErrorPageURL(), resultURL.Path)
}

func TestServeAddCommentSuccessComment(t *testing.T) {
	setup()
	defer teardown()

	testTicket, err := createDummyTicket()
	assert.Nil(t, err)

	form := url.Values{}
	form.Add("comment", "My comment")
	form.Add("sendoption", "comments")

	req := httptest.NewRequest(http.MethodPost, "/addComment", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	req.Header.Set("Referer", "/tickets/"+strconv.Itoa(testTicket.Id))
	req.Form = form

	uuid := utils.CreateUUID(64)
	req.AddCookie(&http.Cookie{
		Name:     "session-id",
		Value:    uuid,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   60 * 60,
	})

	rr := httptest.NewRecorder()
	createUser("Test123", "Aa!123456")
	assert.Nil(t, loginUser(rr, "Test123", "Aa!123456", uuid))

	handler := http.HandlerFunc(ServeAddComment)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusMovedPermanently)
	resultURL, err := rr.Result().Location()
	assert.Nil(t, err)
	assert.Equal(t, req.Referer(), resultURL.Path)
}

func TestServeAddCommentSuccessEmail(t *testing.T) {
	setup()
	defer teardown()

	testTicket, err := createDummyTicket()
	assert.Nil(t, err)

	form := url.Values{}
	form.Add("comment", "My comment")
	form.Add("sendoption", "customer")

	req := httptest.NewRequest(http.MethodPost, "/addComment", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	req.Header.Set("Referer", "/tickets/"+strconv.Itoa(testTicket.Id))
	req.Form = form

	uuid := utils.CreateUUID(64)
	req.AddCookie(&http.Cookie{
		Name:     "session-id",
		Value:    uuid,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   60 * 60,
	})

	rr := httptest.NewRecorder()
	createUser("Test123", "Aa!123456")
	assert.Nil(t, loginUser(rr, "Test123", "Aa!123456", uuid))

	handler := http.HandlerFunc(ServeAddComment)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusMovedPermanently)
	resultURL, err := rr.Result().Location()
	assert.Nil(t, err)
	assert.Equal(t, req.Referer(), resultURL.Path)
}

func createDummyTicket() (utils.Ticket, error) {
	return utils.CreateTicket("test@gmail.com", "Subject Dummy", "Message dummy")
}

func TestServeTicketAssignmentUnauthorized(t *testing.T) {
	setup()
	defer teardown()

	form := url.Values{}

	req := httptest.NewRequest(http.MethodPost, "/assignTicket", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	req.Form = form

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ServeTicketAssignment)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusFound)
	resultURL, err := rr.Result().Location()
	assert.Nil(t, err)
	assert.Equal(t, utils.ErrorUnauthorized.ErrorPageURL(), resultURL.Path)
}

func TestServeTicketAssignmentInvalidURL(t *testing.T) {
	setup()
	defer teardown()

	form := url.Values{}
	form.Add("editor", "Test123")

	req := httptest.NewRequest(http.MethodPost, "/assignTicket", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	req.Header.Set("Referer", "/tickets/MyTicket")
	req.Form = form

	uuid := utils.CreateUUID(64)
	req.AddCookie(&http.Cookie{
		Name:     "session-id",
		Value:    uuid,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   60 * 60,
	})

	rr := httptest.NewRecorder()
	createUser("Test123", "Aa!123456")
	assert.Nil(t, loginUser(rr, "Test123", "Aa!123456", uuid))

	handler := http.HandlerFunc(ServeTicketAssignment)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusFound)
	resultURL, err := rr.Result().Location()
	assert.Nil(t, err)
	assert.Equal(t, utils.ErrorURLParsing.ErrorPageURL(), resultURL.Path)
}

func TestServeTicketAssignmentInvalidEditor(t *testing.T) {
	setup()
	defer teardown()

	testTicket, err := createDummyTicket()
	assert.Nil(t, err)

	form := url.Values{}
	form.Add("editor", "Test")

	req := httptest.NewRequest(http.MethodPost, "/assignTicket", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	req.Header.Set("Referer", "/tickets/"+strconv.Itoa(testTicket.Id))
	req.Form = form

	uuid := utils.CreateUUID(64)
	req.AddCookie(&http.Cookie{
		Name:     "session-id",
		Value:    uuid,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   60 * 60,
	})

	rr := httptest.NewRecorder()
	createUser("Test123", "Aa!123456")
	assert.Nil(t, loginUser(rr, "Test123", "Aa!123456", uuid))

	handler := http.HandlerFunc(ServeTicketAssignment)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusFound)
	resultURL, err := rr.Result().Location()
	assert.Nil(t, err)
	assert.Equal(t, utils.ErrorInvalidInputs.ErrorPageURL(), resultURL.Path)
}

func TestServeTicketAssignmentInvalidTicketID(t *testing.T) {
	setup()
	defer teardown()

	form := url.Values{}
	form.Add("editor", "Test123")

	req := httptest.NewRequest(http.MethodPost, "/assignTicket", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	req.Header.Set("Referer", "/tickets/1337")
	req.Form = form

	uuid := utils.CreateUUID(64)
	req.AddCookie(&http.Cookie{
		Name:     "session-id",
		Value:    uuid,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   60 * 60,
	})

	rr := httptest.NewRecorder()
	createUser("Test123", "Aa!123456")
	assert.Nil(t, loginUser(rr, "Test123", "Aa!123456", uuid))

	handler := http.HandlerFunc(ServeTicketAssignment)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusFound)
	resultURL, err := rr.Result().Location()
	assert.Nil(t, err)
	assert.Equal(t, utils.ErrorDataStoring.ErrorPageURL(), resultURL.Path)
}

func TestServeTicketAssignmentSuccessWithoutRedirect(t *testing.T) {
	setup()
	defer teardown()

	testTicket, err := createDummyTicket()
	assert.Nil(t, err)

	form := url.Values{}
	form.Add("editor", "Test123")

	req := httptest.NewRequest(http.MethodPost, "/assignTicket", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	req.Header.Set("Referer", "/tickets/"+strconv.Itoa(testTicket.Id))
	req.Form = form

	uuid := utils.CreateUUID(64)
	req.AddCookie(&http.Cookie{
		Name:     "session-id",
		Value:    uuid,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   60 * 60,
	})

	rr := httptest.NewRecorder()
	createUser("Test123", "Aa!123456")
	assert.Nil(t, loginUser(rr, "Test123", "Aa!123456", uuid))

	handler := http.HandlerFunc(ServeTicketAssignment)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusFound)
	resultURL, err := rr.Result().Location()
	assert.Nil(t, err)
	assert.Equal(t, req.Referer(), resultURL.Path)
}

func TestServeTicketAssignmentSuccessWithRedirect(t *testing.T) {
	setup()
	defer teardown()

	testTicket, err := createDummyTicket()
	assert.Nil(t, err)

	form := url.Values{}
	form.Add("editor", "Test123456")

	req := httptest.NewRequest(http.MethodPost, "/assignTicket", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	req.Header.Set("Referer", "/tickets/"+strconv.Itoa(testTicket.Id))
	req.Form = form

	uuid := utils.CreateUUID(64)
	req.AddCookie(&http.Cookie{
		Name:     "session-id",
		Value:    uuid,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   60 * 60,
	})

	rr := httptest.NewRecorder()
	createUser("Test123", "Aa!123456")
	createUser("Test123456", "Aa!123456")
	assert.Nil(t, loginUser(rr, "Test123", "Aa!123456", uuid))

	handler := http.HandlerFunc(ServeTicketAssignment)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusFound)
	resultURL, err := rr.Result().Location()
	assert.Nil(t, err)
	assert.Equal(t, "/tickets/", resultURL.Path)
}

func TestServeTicketReleaseUnauthorized(t *testing.T) {
	setup()
	defer teardown()

	form := url.Values{}

	req := httptest.NewRequest(http.MethodPost, "/releaseTicket", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	req.Form = form

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ServeTicketRelease)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusFound)
	resultURL, err := rr.Result().Location()
	assert.Nil(t, err)
	assert.Equal(t, utils.ErrorUnauthorized.ErrorPageURL(), resultURL.Path)
}

func TestServeTicketReleaseInvalidURL(t *testing.T) {
	setup()
	defer teardown()

	form := url.Values{}

	req := httptest.NewRequest(http.MethodPost, "/releaseTicket", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	req.Header.Set("Referer", "/tickets/MyTicket")
	req.Form = form

	uuid := utils.CreateUUID(64)
	req.AddCookie(&http.Cookie{
		Name:     "session-id",
		Value:    uuid,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   60 * 60,
	})

	rr := httptest.NewRecorder()
	createUser("Test123", "Aa!123456")
	assert.Nil(t, loginUser(rr, "Test123", "Aa!123456", uuid))

	handler := http.HandlerFunc(ServeTicketRelease)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusFound)
	resultURL, err := rr.Result().Location()
	assert.Nil(t, err)
	assert.Equal(t, utils.ErrorURLParsing.ErrorPageURL(), resultURL.Path)
}

func TestServeTicketReleaseInvalidTicketID(t *testing.T) {
	setup()
	defer teardown()

	form := url.Values{}

	req := httptest.NewRequest(http.MethodPost, "/releaseTicket", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	req.Header.Set("Referer", "/tickets/1337")
	req.Form = form

	uuid := utils.CreateUUID(64)
	req.AddCookie(&http.Cookie{
		Name:     "session-id",
		Value:    uuid,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   60 * 60,
	})

	rr := httptest.NewRecorder()
	createUser("Test123", "Aa!123456")
	assert.Nil(t, loginUser(rr, "Test123", "Aa!123456", uuid))

	handler := http.HandlerFunc(ServeTicketRelease)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusFound)
	resultURL, err := rr.Result().Location()
	assert.Nil(t, err)
	assert.Equal(t, utils.ErrorDataFetching.ErrorPageURL(), resultURL.Path)
}

func TestServeTicketReleaseInvalidUser(t *testing.T) {
	setup()
	defer teardown()

	createUser("Test1234567", "Aa!123456")
	testTicket, err := createDummyTicket()
	assert.Nil(t, utils.ChangeEditor(testTicket.Id, "Test1234567"))
	assert.Nil(t, err)

	form := url.Values{}

	req := httptest.NewRequest(http.MethodPost, "/releaseTicket", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	req.Header.Set("Referer", "/tickets/"+strconv.Itoa(testTicket.Id))
	req.Form = form

	uuid := utils.CreateUUID(64)
	req.AddCookie(&http.Cookie{
		Name:     "session-id",
		Value:    uuid,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   60 * 60,
	})

	rr := httptest.NewRecorder()
	createUser("Test123", "Aa!123456")
	assert.Nil(t, loginUser(rr, "Test123", "Aa!123456", uuid))

	handler := http.HandlerFunc(ServeTicketRelease)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusFound)
	resultURL, err := rr.Result().Location()
	assert.Nil(t, err)
	assert.Equal(t, utils.ErrorUnauthorized.ErrorPageURL(), resultURL.Path)
}

func TestServeTicketReleaseSuccess(t *testing.T) {
	setup()
	defer teardown()

	testTicket, err := createDummyTicket()
	assert.Nil(t, utils.ChangeEditor(testTicket.Id, "Test123"))
	assert.Nil(t, err)

	form := url.Values{}

	req := httptest.NewRequest(http.MethodPost, "/releaseTicket", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	req.Header.Set("Referer", "/tickets/"+strconv.Itoa(testTicket.Id))
	req.Form = form

	uuid := utils.CreateUUID(64)
	req.AddCookie(&http.Cookie{
		Name:     "session-id",
		Value:    uuid,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   60 * 60,
	})

	rr := httptest.NewRecorder()
	createUser("Test123", "Aa!123456")
	assert.Nil(t, loginUser(rr, "Test123", "Aa!123456", uuid))

	handler := http.HandlerFunc(ServeTicketRelease)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusFound)
	resultURL, err := rr.Result().Location()
	assert.Nil(t, err)
	assert.Equal(t, req.Referer(), resultURL.Path)
}

func TestServeTicketsUnauthorized(t *testing.T) {
	setup()
	defer teardown()

	req := httptest.NewRequest(http.MethodPost, "/tickets/", nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ServeTickets)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusFound)
	resultURL, err := rr.Result().Location()
	assert.Nil(t, err)
	assert.Equal(t, utils.ErrorUnauthorized.ErrorPageURL(), resultURL.Path)
}

func TestServeTicketsShowTemplate(t *testing.T) {
	setup()
	defer teardown()

	req := httptest.NewRequest(http.MethodPost, "/ticket/", nil)

	uuid := utils.CreateUUID(64)
	req.AddCookie(&http.Cookie{
		Name:     "session-id",
		Value:    uuid,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   60 * 60,
	})

	rr := httptest.NewRecorder()
	createUser("Test123", "Aa!123456")
	assert.Nil(t, loginUser(rr, "Test123", "Aa!123456", uuid))
	handler := http.HandlerFunc(ServeTickets)

	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestServeTicketsInvalidTicketID(t *testing.T) {
	setup()
	defer teardown()

	req := httptest.NewRequest(http.MethodPost, "/ticket/1337", nil)

	uuid := utils.CreateUUID(64)
	req.AddCookie(&http.Cookie{
		Name:     "session-id",
		Value:    uuid,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   60 * 60,
	})

	rr := httptest.NewRecorder()
	createUser("Test123", "Aa!123456")
	assert.Nil(t, loginUser(rr, "Test123", "Aa!123456", uuid))
	handler := http.HandlerFunc(ServeTickets)

	handler.ServeHTTP(rr, req)
	resultURL, err := rr.Result().Location()
	assert.Nil(t, err)
	assert.Equal(t, utils.ErrorInvalidTicketID.ErrorPageURL(), resultURL.Path)
}

func TestServeTicketsSuccess(t *testing.T) {
	setup()
	defer teardown()

	createUser("Test124563", "Aa!123456")

	testTicket, err := createDummyTicket()
	assert.Nil(t, err)

	req := httptest.NewRequest(http.MethodPost, "/ticket/"+strconv.Itoa(testTicket.Id), nil)

	uuid := utils.CreateUUID(64)
	req.AddCookie(&http.Cookie{
		Name:     "session-id",
		Value:    uuid,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   60 * 60,
	})

	rr := httptest.NewRecorder()
	createUser("Test123", "Aa!123456")
	assert.Nil(t, loginUser(rr, "Test123", "Aa!123456", uuid))
	handler := http.HandlerFunc(ServeTickets)

	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestServeCloseTicketUnauthorized(t *testing.T) {
	setup()
	defer teardown()

	req := httptest.NewRequest(http.MethodPost, "/closeTicket", nil)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ServeCloseTicket)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusFound)
	resultURL, err := rr.Result().Location()
	assert.Nil(t, err)
	assert.Equal(t, utils.ErrorUnauthorized.ErrorPageURL(), resultURL.Path)
}

func TestServeCloseTicketInvalidURL(t *testing.T) {
	setup()
	defer teardown()

	req := httptest.NewRequest(http.MethodPost, "/closeTicket", nil)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	req.Header.Set("Referer", "/tickets/MyTicket")

	uuid := utils.CreateUUID(64)
	req.AddCookie(&http.Cookie{
		Name:     "session-id",
		Value:    uuid,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   60 * 60,
	})

	rr := httptest.NewRecorder()
	createUser("Test123", "Aa!123456")
	assert.Nil(t, loginUser(rr, "Test123", "Aa!123456", uuid))

	handler := http.HandlerFunc(ServeCloseTicket)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusFound)
	resultURL, err := rr.Result().Location()
	assert.Nil(t, err)
	assert.Equal(t, utils.ErrorURLParsing.ErrorPageURL(), resultURL.Path)
}

func TestServeCloseTicketSuccess(t *testing.T) {
	setup()
	defer teardown()

	testTicket, err := createDummyTicket()
	assert.Nil(t, err)

	req := httptest.NewRequest(http.MethodPost, "/closeTicket", nil)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	req.Header.Set("Referer", "/ticket/"+strconv.Itoa(testTicket.Id))

	uuid := utils.CreateUUID(64)
	req.AddCookie(&http.Cookie{
		Name:     "session-id",
		Value:    uuid,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   60 * 60,
	})

	rr := httptest.NewRecorder()
	createUser("Test123", "Aa!123456")
	assert.Nil(t, loginUser(rr, "Test123", "Aa!123456", uuid))

	handler := http.HandlerFunc(ServeCloseTicket)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusMovedPermanently)
	resultURL, err := rr.Result().Location()
	assert.Nil(t, err)
	assert.Equal(t, "/tickets/", resultURL.Path)
}
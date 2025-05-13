package controller

import (
	"bytes"
	"encoding/json"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/hielkefellinger/go-dnd/app/helpers"
	"github.com/hielkefellinger/go-dnd/app/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestLoginPageWithoutUser(t *testing.T) {
	// Arrange
	w := httptest.NewRecorder()
	ctx := createGinTestContext(w)

	// Act
	LoginPage(ctx)

	// Assert
	assert.EqualValues(t, http.StatusOK, w.Code)
	bodyString := w.Body.String()
	assert.Equal(t, strings.Contains(bodyString, "<b>LOGIN</b>"), true)
	assert.Equal(t, strings.Contains(bodyString, "/u/login"), true)
	assert.Equal(t, strings.Contains(bodyString, "/u/register"), true)
	assert.Equal(t, strings.Contains(bodyString, "Register a new user?"), true)
}

func TestRegisterPageWithoutUser(t *testing.T) {
	// Arrange
	w := httptest.NewRecorder()
	ctx := createGinTestContext(w)

	// Act
	RegisterPage(ctx)

	// Assert
	assert.EqualValues(t, http.StatusOK, w.Code)
	bodyString := w.Body.String()
	assert.Equal(t, strings.Contains(bodyString, "<b>REGISTER</b>"), true)
	assert.Equal(t, strings.Contains(bodyString, "/u/login"), true)
	assert.Equal(t, strings.Contains(bodyString, "/u/register"), true)
	assert.Equal(t, strings.Contains(bodyString, "Already have a user?"), true)
}

func TestRegisterUser(t *testing.T) {
	// Arrange
	w := httptest.NewRecorder()
	ctx := createGinTestContext(w)
	jsonValue, _ := json.Marshal(models.User{Name: "TEST_USER", Password: "PASSWORD", PasswordCheck: "PASSWORD!"})
	reqReg, _ := http.NewRequest("POST", "/u/register", bytes.NewBuffer(jsonValue))
	ctx.Request = reqReg
	ctx.Request.Header = http.Header{"Content-Type": []string{"application/json"}}

	// Act
	Register(ctx)

	// Assert
	assert.EqualValues(t, http.StatusBadRequest, w.Code)
	bodyString := w.Body.String()
	assert.Equal(t, strings.Contains(bodyString, "<b>REGISTER</b>"), true)
	assert.Equal(t, strings.Contains(bodyString, "/u/register"), true)
	assert.Equal(t, strings.Contains(bodyString, "passwords do not match"), true)
	assert.Equal(t, strings.Contains(bodyString, "Already have a user?"), true)
}

func TestLoginUser(t *testing.T) {
	// Arrange
	form := url.Values{}
	form.Add("username", "TEST_USER")
	form.Add("password", "PASSWORD")

	reqLogin, _ := http.NewRequest("POST", "/u/login", strings.NewReader(form.Encode()))
	reqLogin.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()
	ctx := createGinTestContext(w)
	ctx.Request = reqLogin

	mockDb, mock, _ := sqlmock.New()
	dialector := postgres.New(postgres.Config{
		Conn:       mockDb,
		DriverName: "postgres",
	})
	db, _ := gorm.Open(dialector, &gorm.Config{})
	models.DB = db

	t.Setenv("CRYPT_COST", "2")
	hashByteArray, _ := helpers.HashPassword("PASSWORD")
	rows := sqlmock.NewRows([]string{"Id", "Name", "Password"}).AddRow("0", "TEST_USER", string(hashByteArray))
	mock.ExpectQuery(`SELECT`).WillReturnRows(rows)

	// Act
	Login(ctx)

	// Assert
	assert.EqualValues(t, http.StatusOK, w.Code)
}

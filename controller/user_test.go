package controller

import (
	"encoding/json"
	"github.com/DATA-DOG/go-sqlmock"
	G "github.com/NoCLin/douyin-backend-go/config/global"
	"github.com/appleboy/gofight/v2"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
)

func assertBasicResponse(t *testing.T, r gofight.HTTPResponse, statusCode int) map[string]interface{} {

	assert.Equal(t, r.Code, http.StatusOK)
	obj := make(map[string]interface{})
	err := json.Unmarshal(r.Body.Bytes(), &obj)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, float64(statusCode), obj["StatusCode"])
	return obj
}

func TestRegister(t *testing.T) {

	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)
	engine.GET("/", Register)
	r := gofight.New()

	sqlQueryUser := "SELECT * FROM `users` WHERE name = ? AND `users`.`deleted_at` IS NULL LIMIT 1"
	sqlCreateUser := "INSERT INTO `users` (`created_at`,`updated_at`,`deleted_at`,`name`,`password`,`password_hashed`) VALUES (?,?,?,?,?,?)"
	t.Run("OK", func(t *testing.T) {

		username := "test_username"
		password := "test_password"
		q := gofight.H{
			"username": username,
			"password": password,
		}
		gfDebug := false

		mock.ExpectQuery(regexp.QuoteMeta(sqlQueryUser)).WithArgs(username).WillReturnRows(sqlmock.NewRows([]string{"id"}))
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(sqlCreateUser)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), nil, username, password, "").
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		r.GET("/").SetQuery(q).SetDebug(gfDebug).
			Run(engine, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
				obj := assertBasicResponse(t, r, 0)
				assert.Contains(t, obj["Token"], "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9") // jwt header
			})

	})

	t.Run("UsernameTooLong", func(t *testing.T) {

		username := "test_username_long_long_long_long_long_long_long_long"
		password := "test_password"
		q := gofight.H{
			"username": username,
			"password": password,
		}
		gfDebug := false

		r.GET("/").SetQuery(q).SetDebug(gfDebug).
			Run(engine, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
				_ = assertBasicResponse(t, r, -1)
			})

	})

	// we make sure that all expectations were met
	if err := G.DBMock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections:\n%s", err)
	}

}

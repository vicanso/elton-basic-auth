package basicauth

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vicanso/cod"
	"github.com/vicanso/hes"
)

func TestNoVildatePanic(t *testing.T) {
	assert := assert.New(t)
	defer func() {
		r := recover()
		assert.NotNil(r)
		assert.Equal(r.(error), errRequireValidateFunction)
	}()

	New(Config{})
}

func TestBasicAuth(t *testing.T) {
	m := New(Config{
		Validate: func(account, pwd string, c *cod.Context) (bool, error) {
			if account == "tree.xie" && pwd == "password" {
				return true, nil
			}
			if account == "n" {
				return false, hes.New("account is invalid")
			}
			return false, nil
		},
	})
	req := httptest.NewRequest("GET", "https://aslant.site/", nil)

	t.Run("skip", func(t *testing.T) {
		assert := assert.New(t)
		done := false
		mSkip := New(Config{
			Validate: func(account, pwd string, c *cod.Context) (bool, error) {
				return false, nil
			},
			Skipper: func(c *cod.Context) bool {
				return true
			},
		})
		d := cod.New()
		d.Use(mSkip)
		d.GET("/", func(c *cod.Context) error {
			done = true
			return nil
		})
		resp := httptest.NewRecorder()
		d.ServeHTTP(resp, req)
		assert.True(done)
	})

	t.Run("no auth header", func(t *testing.T) {
		assert := assert.New(t)
		d := cod.New()
		d.Use(m)
		d.GET("/", func(c *cod.Context) error {
			return nil
		})
		resp := httptest.NewRecorder()
		d.ServeHTTP(resp, req)
		assert.Equal(resp.Code, http.StatusUnauthorized)
		assert.Equal(resp.Header().Get(cod.HeaderWWWAuthenticate), `basic realm="basic auth tips"`)
	})

	t.Run("auth validate fail", func(t *testing.T) {
		assert := assert.New(t)
		d := cod.New()
		d.Use(m)
		d.GET("/", func(c *cod.Context) error {
			return nil
		})
		req.Header.Set(cod.HeaderAuthorization, "basic YTpi")
		resp := httptest.NewRecorder()
		d.ServeHTTP(resp, req)
		assert.Equal(resp.Code, http.StatusUnauthorized)
		assert.Equal(resp.Body.String(), "category=cod-basic-auth, message=unAuthorized")

		req.Header.Set(cod.HeaderAuthorization, "basic bjph")
		resp = httptest.NewRecorder()
		d.ServeHTTP(resp, req)
		assert.Equal(resp.Code, http.StatusBadRequest)
		assert.Equal(resp.Body.String(), "message=account is invalid")
	})

	t.Run("validate error", func(t *testing.T) {
		assert := assert.New(t)
		mValidateFail := New(Config{
			Validate: func(account, pwd string, c *cod.Context) (bool, error) {
				return false, errors.New("abcd")
			},
		})
		d := cod.New()
		d.Use(mValidateFail)
		d.GET("/", func(c *cod.Context) error {
			return nil
		})
		resp := httptest.NewRecorder()
		d.ServeHTTP(resp, req)
		assert.Equal(resp.Code, http.StatusBadRequest)
		assert.Equal(resp.Body.String(), "category=cod-basic-auth, message=abcd")
	})

	t.Run("auth success", func(t *testing.T) {
		assert := assert.New(t)
		d := cod.New()
		d.Use(m)
		done := false
		d.GET("/", func(c *cod.Context) error {
			done = true
			return nil
		})
		req.Header.Set(cod.HeaderAuthorization, "basic dHJlZS54aWU6cGFzc3dvcmQ=")
		resp := httptest.NewRecorder()
		d.ServeHTTP(resp, req)
		assert.True(done)
	})
}

// https://stackoverflow.com/questions/50120427/fail-unit-tests-if-coverage-is-below-certain-percentage
func TestMain(m *testing.M) {
	// call flag.Parse() here if TestMain uses flags
	rc := m.Run()

	// rc 0 means we've passed,
	// and CoverMode will be non empty if run with -cover
	if rc == 0 && testing.CoverMode() != "" {
		c := testing.Coverage()
		if c < 0.9 {
			fmt.Println("Tests passed but coverage failed at", c)
			rc = -1
		}
	}
	os.Exit(rc)
}

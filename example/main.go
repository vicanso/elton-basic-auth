package main

import (
	"bytes"

	"github.com/vicanso/elton"
	basicauth "github.com/vicanso/elton-basic-auth"
	"github.com/vicanso/hes"
)

func main() {
	d := elton.New()

	d.Use(basicauth.New(basicauth.Config{
		Validate: func(account, pwd string, c *elton.Context) (bool, error) {
			if account == "tree.xie" && pwd == "password" {
				return true, nil
			}
			if account == "n" {
				return false, hes.New("account is invalid")
			}
			return false, nil
		},
	}))

	d.GET("/", func(c *elton.Context) (err error) {
		c.BodyBuffer = bytes.NewBufferString("hello world")
		return
	})

	err := d.ListenAndServe(":3000")
	if err != nil {
		panic(err)
	}
}

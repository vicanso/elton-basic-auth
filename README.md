# elton-basic-auth

[![Build Status](https://img.shields.io/travis/vicanso/elton-basic-auth.svg?label=linux+build)](https://travis-ci.org/vicanso/elton-basic-auth)


Basic auth middleware for elton.

```go
package main

import (
	"bytes"

	"github.com/vicanso/elton"
	basicauth "github.com/vicanso/elton-basic-auth"
	"github.com/vicanso/hes"
)

func main() {
	e := elton.New()

	e.Use(basicauth.New(basicauth.Config{
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

	e.GET("/", func(c *elton.Context) (err error) {
		c.BodyBuffer = bytes.NewBufferString("hello world")
		return
	})

	e.ListenAndServe(":3000")
}

```
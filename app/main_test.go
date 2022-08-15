package main

import (
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestSetupHandler(t *testing.T) {
	e := echo.New()
	SetupHandler(e, nil)
	assert.Greater(t, len(e.Routes()), 0)
}

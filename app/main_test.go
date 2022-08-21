package main

import (
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestSetupLogger(t *testing.T) {
	e := echo.New()
	assert.Panics(t, func() { setupLogger(e) })
}

func TestSetupHandler(t *testing.T) {
	e := echo.New()
	setupHandler(e, nil)
	assert.Greater(t, len(e.Routes()), 0)
}

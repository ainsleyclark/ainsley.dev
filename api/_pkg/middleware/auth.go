// Copyright 2023 ainsley.dev LTD. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/ainsleyclark/ainsley.dev/api/_pkg/environment"
	"github.com/ainsleyclark/ainsley.dev/api/_sdk"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const (
	// AuthHeader is the header used for authentication via an
	// API Key
	AuthHeader = "X-API-Key"
	// OriginURL is the URL of the site to compare against
	// in production.
	OriginURL = "https://ainsley.dev"
)

// Auth validates API request and determines if the AuthHeader value is of equal
// to the header send in the request.
// Returns errors.INVALID if there is a mismatch.
func Auth(cfg *environment.Config) echo.MiddlewareFunc {
	return middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		KeyLookup:              "header:" + AuthHeader,
		ContinueOnIgnoredError: false,
		Validator: func(auth string, ctx echo.Context) (bool, error) {
			if !cfg.IsProduction() {
				return auth == cfg.APIKey, nil
			}
			origin := ctx.Request().Header.Get("Origin")
			if !strings.Contains(origin, cfg.URL) || origin == "" {
				return false, fmt.Errorf("bad origin: %s, with comparison: %s", origin, cfg.URL)
			}
			return auth == cfg.APIKey, nil
		},
		ErrorHandler: func(err error, ctx echo.Context) error {
			return ctx.JSON(http.StatusUnauthorized, sdk.HTTPError{
				Code:      "<unauthorized>",
				Error:     err.Error(),
				Message:   "Not authorised",
				Operation: "API.Auth",
			})
		},
	})
}

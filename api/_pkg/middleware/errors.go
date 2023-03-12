// Copyright 2023 ainsley.dev LTD. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package middleware

import (
	"github.com/ainsleyclark/ainsley.dev/api/_pkg/logger"
	"github.com/ainsleyclark/ainsley.dev/api/_sdk"
	"github.com/ainsleyclark/errors"
	"github.com/labstack/echo/v4"
	"net/http"
)

// ErrorHandler writes the json-encoded error message to the response.
// If the error is of the validationErrors, it will be marshalled
// for correct formatting.
func ErrorHandler(err error, ctx echo.Context) {
	code := http.StatusInternalServerError
	resp := sdk.HTTPError{
		Message: err.Error(),
	}
	e, ok := err.(*errors.Error)
	if ok {
		code = e.HTTPStatusCode()
		resp = sdk.HTTPError{
			Code:      e.Code,
			Error:     err.Error(),
			Message:   e.Message,
			Operation: e.Operation,
		}
	}
	err = ctx.JSON(code, resp)
	if err != nil {
		logger.WithError(err).Error("Error sending payload")
	}
}

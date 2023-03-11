// Copyright 2023 ainsley.dev LTD. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package middleware

import (
	"encoding/json"
	sdk "github.com/ainsleyclark/ainsley.dev/gen/sdk/go"
	"github.com/ainsleyclark/errors"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestErrorHandler(t *testing.T) {
	tt := map[string]struct {
		input error
		want  sdk.HTTPError
	}{
		"Default": {
			input: errors.New("error"),
			want: sdk.HTTPError{
				Message: "error",
			},
		},
		"With Error": {
			input: &errors.Error{
				Code:      errors.INTERNAL,
				Message:   "message",
				Operation: "op",
				Err:       errors.New("error"),
			},
			want: sdk.HTTPError{
				Code:      errors.INTERNAL,
				Error:     "<internal> op: error, message",
				Message:   "message",
				Operation: "op",
			},
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/", nil)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)
			e.POST("/", func(c echo.Context) error {
				return test.input
			})
			ErrorHandler(test.input, ctx)
			response := sdk.HTTPError{}
			err := json.NewDecoder(rec.Body).Decode(&response)
			require.NoError(t, err)
			assert.Equal(t, test.want, response)
		})
	}
}

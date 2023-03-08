// Copyright 2023 ainsley.dev LTD. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"github.com/ainsleyclark/ainsley.dev/api/_pkg/app"
	"log"
	"net/http"

	"github.com/ainsleyclark/ainsley.dev/api/_pkg/environment"
	"github.com/ainsleyclark/ainsley.dev/api/_pkg/gateway/mail"
	"github.com/ainsleyclark/ainsley.dev/api/_pkg/gateway/slack"
	"github.com/ainsleyclark/ainsley.dev/api/_pkg/logger"

	"github.com/ainsleyclark/ainsley.dev/api/_pkg/httpservice"
	sdk "github.com/ainsleyclark/ainsley.dev/gen/sdk/go"
	"github.com/labstack/echo/v4"
)

var (
	// app is the main Echo application handler.
	e *echo.Echo
	// handler is the service to handle the incoming routes.
	handler *httpservice.Handler
)

// init bootstraps the main application by creating a new Echo instance
// and registering the API routes along with middleware and
// any configuration.
func init() {
	logger.Bootstrap()
	config, err := environment.New()
	if err != nil {
		log.Fatalln(err.Error())
	}
	e = echo.New()
	app.Bootstrap(e, config)
	mailer, err := mail.New(config)
	if err != nil {
		log.Fatalln(err.Error())
	}
	handler = &httpservice.Handler{
		Slack:  slack.New(config.SlackToken),
		Mailer: mailer,
	}
	sdk.RegisterHandlersWithBaseURL(e, handler, "/api")
}

// Handler is the main entrypoint to the application.
// Vercel detects this http.HandlerFunc signature to use
// within serverless functions.
func Handler(w http.ResponseWriter, r *http.Request) {
	e.ServeHTTP(w, r)
}

// Copyright 2023 ainsley.dev LTD. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"log"
	"net/http"

	"github.com/ainsleyclark/ainsley.dev/api/_pkg/logger"
	"github.com/sirupsen/logrus"

	"github.com/ainsleyclark/ainsley.dev/api/_pkg/environment"
	"github.com/ainsleyclark/ainsley.dev/api/_pkg/gateway/mail"
	"github.com/ainsleyclark/ainsley.dev/api/_pkg/gateway/slack"

	"github.com/ainsleyclark/ainsley.dev/api/_pkg/api"
	"github.com/ainsleyclark/ainsley.dev/api/_pkg/httpservice"
	sdk "github.com/ainsleyclark/ainsley.dev/gen/sdk/go"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	// app is the main Echo application handler.
	app *echo.Echo
	// handler is the service to handle the incoming routes.
	handler *httpservice.Handler
)

// init bootstraps the main application by creating a new Echo instance
// and registering the API routes along with middleware and
// any configuration.
func init() {
	app = echo.New()
	config, err := environment.New()
	if err != nil {
		log.Fatalln(err.Error())
	}
	mailer, err := mail.New(config)
	if err != nil {
		log.Fatalln(err.Error())
	}
	handler = &httpservice.Handler{
		Slack:  slack.New(config.SlackToken),
		Mailer: mailer,
	}
	echo.NotFoundHandler = api.NotFoundHandler
	app.HTTPErrorHandler = api.ErrorHandler
	app.Use(api.Auth(config))
	app.Use(api.CORS(config))
	app.Use(api.RequestID())
	app.Use(api.Logger())
	app.Use(middleware.GzipWithConfig(middleware.GzipConfig{Level: 5}))
	app.Pre(middleware.AddTrailingSlash())

	// TODO, move to log package and don't export formatter
	logger.DefaultLogger.SetFormatter(&logrus.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
	})
	sdk.RegisterHandlersWithBaseURL(app, handler, "/api")
}

// Handler is the main entrypoint to the application.
// Vercel detects this http.HandlerFunc signature to use
// within serverless functions.
func Handler(w http.ResponseWriter, r *http.Request) {
	app.ServeHTTP(w, r)
}

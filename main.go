// Copyright 2023 ainsley.dev LTD. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"

	"github.com/ainsleyclark/ainsley.dev/api"
	"github.com/ainsleyclark/ainsley.dev/api/_pkg/httpservice"
	"github.com/ainsleyclark/ainsley.dev/api/_pkg/logger"
	"github.com/ainsleyclark/ainsley.dev/api/_sdk"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Pass server port
	var port string
	flag.StringVar(&port, "port", "3000", "Server listen address")
	flag.Parse()

	// Load env file
	err := godotenv.Load(".env")
	if err != nil {
		logger.Fatal(err)
	}

	// Bootstrap Server
	e := echo.New()
	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Root:   "public",
		Index:  "index.html",
		Browse: false,
		HTML5:  true,
	}))
	handler := api.Bootstrap(e)
	sdk.RegisterHandlers(e.Group(httpservice.BasePath), handler)

	// Start server
	logger.Fatal(e.Start(":" + port))
}

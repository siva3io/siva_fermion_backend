package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"

	cache "fermion/backend_core/controllers/cache"
	scheduler "fermion/backend_core/controllers/scheduler/app"
	"fermion/backend_core/db"

	cmiddleware "fermion/backend_core/middleware"
	"fermion/backend_core/pkg/util"
	"fermion/backend_core/pkg/util/helpers"
	"fermion/integrations"
	"fermion/route"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
	_ "github.com/swaggo/echo-swagger"
)

/*
Copyright (C) 2022 Eunimart Omnichannel Pvt Ltd. (www.eunimart.com)
All rights reserved.
This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.
This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.
You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/
func init() {
	ENV := os.Getenv("ENV")
	env := util.NewEnv()
	env.Load(ENV)
	logrus.Info("choose environment " + ENV)
}

// Define the template registry struct
type TemplateRegistry struct {
	templates *template.Template
}

// Implement e.Renderer interface
func (t *TemplateRegistry) Render(w io.Writer, name string, data interface{}, c echo.Context) error {

	// Add global methods if data is a map
	if viewContext, isMap := data.(map[string]interface{}); isMap {
		viewContext["reverse"] = c.Echo().Reverse
	}
	t.templates.Funcs(template.FuncMap{"mod": func(i, j int) bool { return i%j == 0 }})
	return t.templates.ExecuteTemplate(w, name, data)
}

// @title Eunimart Platform API Docs
// @version 1.0.0
// @description This is the API docs for Eunimart Platform.
// @termsOfService https://eunimart.com/in/terms-conditions/

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

// @contact.name API Support
// @contact.url https://eunimart.com/contact-us-1/
// @contact.email contact@eunimart.com

// @license.name MIT License
// @license.url https://opensource.org/licenses/MIT

// @host siva3.io
// @schemes https
// @BasePath /
func main() {
	var (
		APP  = os.Getenv("APP")
		ENV  = os.Getenv("ENV")
		PORT = os.Getenv("PORT")
		NAME = fmt.Sprintf("%s-%s", APP, ENV)
	)

	// Init
	db.Init()
	db.NoSqlInit()
	// uncomment once you setup elasticsearch
	// elastic.Init()
	// log.Init()
	e := echo.New()

	// Middleware
	e.Use(
		middleware.Logger(),
		middleware.Recover(),
		middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: []string{"*"},
			AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
		}),
		middleware.LoggerWithConfig(middleware.LoggerConfig{
			Format: fmt.Sprintf("\n%s | ${host} | ${time_custom} | ${status} | ${latency_human} | ${remote_ip} | ${method} | ${uri} ",
				NAME,
			),
			CustomTimeFormat: "2006/01/02 15:04:05",
			Output:           os.Stdout,
		}),
	)
	e.HTTPErrorHandler = cmiddleware.ErrorHandler
	e.Validator = &util.CustomValidator{Validator: validator.New()}

	e.Renderer = &TemplateRegistry{
		templates: template.Must(template.New("t").Funcs(template.FuncMap{
			"mod": helpers.Mod,
		}).ParseGlob("backend_core/views/*.html")),
	}

	e.Static("/files", "backend_core/assets")

	// Routes
	route.Init(e.Group(""))
	integrations.UseSubrouteForIntegration(e.Group("integrations"))

	scheduler.Init()
	cache.Init()

	// Start
	e.Logger.Fatal(e.Start(":" + PORT))
}

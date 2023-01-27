package util

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

/*
Copyright (C) 2022 Eunimart Omnichannel Pvt Ltd. (www.eunimart.com)
All rights reserved.
This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Lesser General Public License v3.0 as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.
This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Lesser General Public License v3.0 for more details.
You should have received a copy of the GNU Lesser General Public License v3.0
along with this program.  If not, see <https://www.gnu.org/licenses/lgpl-3.0.html/>.
*/
type Env interface {
	GetString(name string) string
}

type env struct {
	Env
}

type EnvGetter struct{}

func NewEnv() *env {
	return &env{Env: &EnvGetter{}}
}

func (e *env) Load(env string) {

	cwd, _ := os.Getwd()

	var envFile string
	switch env {
	case "STAGING":
		envFile = "staging"
	case "PROD":
		envFile = "production"
	default:
		envFile = "development"
	}

	err := godotenv.Load(cwd + `/.env.` + envFile)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"cause": err,
			"cwd":   cwd,
		}).Fatal("Load .env file error")

		os.Exit(-1)
	}
}

func (r *EnvGetter) GetString(name string) string {
	return os.Getenv(name)
}

func (e *env) GetString(name string) string {
	if nil == e.Env {
		return ""
	}
	return e.Env.GetString(name)
}

func (e *env) GetBool(name string) bool {
	s := e.GetString(name)
	i, err := strconv.ParseBool(s)
	if nil != err {
		return false
	}
	return i
}

func (e *env) GetInt(name string) int {
	s := e.GetString(name)
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return i
}

func (e *env) GetFloat(name string) float64 {
	s := e.GetString(name)
	i, err := strconv.ParseFloat(s, 64)
	if nil != err {
		return 0
	}
	return i
}

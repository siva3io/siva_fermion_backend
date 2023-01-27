package response

import (
	"fmt"
	"net/http"

	"fermion/backend_core/internal/model/pagination"

	"github.com/labstack/echo/v4"
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
type MetaResponse struct {
	Success bool                      `json:"success" default:"true"`
	Message string                    `json:"message" default:"true"`
	Info    *pagination.Paginatevalue `json:"info"`
}

type SuccessResponse struct {
	Meta MetaResponse `json:"meta"`
	Data interface{}  `json:"data"`
}

type ErrorResponse struct {
	Meta  MetaResponse `json:"meta"`
	Error string       `json:"error"`
}

type ValidationErrorResponse struct {
	Meta  MetaResponse `json:"meta"`
	Error interface{}  `json:"error"`
}

type ErrorConstant struct {
	Response     ErrorResponse
	Code         int
	ErrorMessage error
}

func (r *ErrorConstant) Error() string {
	return fmt.Sprintf("error code %d", r.Code)
}

func (r *ErrorConstant) Builder() *ErrorConstant {
	return r
}

const (
	E_DUPLICATE            = "duplicate"
	E_NOT_FOUND            = "not_found"
	E_UNPROCESSABLE_ENTITY = "unprocessable_entity"
	E_UNAUTHORIZED         = "unauthorized"
	E_BAD_REQUEST          = "bad_request"
	E_SERVER_ERROR         = "server_error"
	E_METHOD_NOT_ALLOWED   = "method_not_allowed"
)

var (
	ErrDuplicate = ErrorConstant{
		Response: ErrorResponse{
			Meta: MetaResponse{
				Success: false,
				Message: "Created value already exists",
			},
			Error: E_DUPLICATE,
		},
		Code: http.StatusConflict,
	}
	ErrDataNotFound = ErrorConstant{
		Response: ErrorResponse{
			Meta: MetaResponse{
				Success: false,
				Message: "Data not found",
			},
			Error: E_NOT_FOUND,
		},
		Code: http.StatusNotFound,
	}
	ErrRouteNotFound = ErrorConstant{
		Response: ErrorResponse{
			Meta: MetaResponse{
				Success: false,
				Message: "Route not found",
			},
			Error: E_NOT_FOUND,
		},
		Code: http.StatusNotFound,
	}
	ErrUnprocessableEntity = ErrorConstant{
		Response: ErrorResponse{
			Meta: MetaResponse{
				Success: false,
				Message: "Invalid parameters or payload",
			},
			Error: E_UNPROCESSABLE_ENTITY,
		},
		Code: http.StatusUnprocessableEntity,
	}
	ErrUnauthorized = ErrorConstant{
		Response: ErrorResponse{
			Meta: MetaResponse{
				Success: false,
				Message: "Unauthorized, please login",
			},
			Error: E_UNAUTHORIZED,
		},
		Code: http.StatusUnauthorized,
	}
	ErrBadRequest = ErrorConstant{
		Response: ErrorResponse{
			Meta: MetaResponse{
				Success: false,
				Message: "Bad Request",
			},
			Error: E_BAD_REQUEST,
		},
		Code: http.StatusBadRequest,
	}
	ErrValidation = ErrorConstant{
		Response: ErrorResponse{
			Meta: MetaResponse{
				Success: false,
				Message: "Invalid parameters or payload",
			},
			Error: E_BAD_REQUEST,
		},
		Code: http.StatusBadRequest,
	}
	ErrServerError = ErrorConstant{
		Response: ErrorResponse{
			Meta: MetaResponse{
				Success: false,
				Message: "Something bad happened",
			},
			Error: E_SERVER_ERROR,
		},
		Code: http.StatusInternalServerError,
	}
	ErrMethodNotAllowed = ErrorConstant{
		Response: ErrorResponse{
			Meta: MetaResponse{
				Success: false,
				Message: "Method Not Allowed",
			},
			Error: E_METHOD_NOT_ALLOWED,
		},
		Code: http.StatusMethodNotAllowed,
	}
)

func BuildError(err ErrorConstant, msg error) error {
	err.ErrorMessage = msg
	return &err
}

func BuildErrorCustom(code int, error string, message string) error {
	return &ErrorConstant{
		Response: ErrorResponse{
			Meta: MetaResponse{
				Success: false,
				Message: message,
			},
			Error: error,
		},
		Code: code,
	}
}

func RespSuccess(c echo.Context, message string, data interface{}) error {
	response := SuccessResponse{
		Meta: MetaResponse{
			Success: true,
			Message: message,
		},
		Data: data,
	}
	return c.JSON(http.StatusOK, response)
}

func SuccessWithError(c echo.Context, message string, data interface{}) error {
	response := SuccessResponse{
		Meta: MetaResponse{
			Success: false,
			Message: message,
		},
		Data: data,
	}
	return c.JSON(http.StatusOK, response)
}

func RespSuccessInfo(c echo.Context, message string, data interface{}, pagination *pagination.Paginatevalue) error {
	response := SuccessResponse{
		Meta: MetaResponse{
			Success: true,
			Message: message,
			Info:    pagination,
		},
		Data: data,
	}
	return c.JSON(http.StatusOK, response)
}

func RespError(c echo.Context, err error) error {

	// body, e := ioutil.ReadAll(c.Request().Body)
	// if e != nil {
	// 	logrus.Warn("error read body, message : ", e.Error())
	// }

	// bHeader, e := json.Marshal(c.Request().Header)
	// if e != nil {
	// 	logrus.Warn("error read header, message : ", e.Error())
	// }

	re, ok := err.(*ErrorConstant)

	if ok {

		// log.InsertErrorLog(c.Request().Context(), &log.LogError{
		// 	ID:           shortid.MustGenerate(),
		// 	Header:       string(bHeader),
		// 	Body:         string(body),
		// 	URL:          c.Request().URL.Path,
		// 	HttpMethod:   c.Request().Method,
		// 	ErrorMessage: re.Builder().ErrorMessage.Error(),
		// 	Level:        "Error",
		// 	AppName:      os.Getenv("APP"),
		// 	Version:      os.Getenv("VERSION"),
		// 	Env:          os.Getenv("ENV"),
		// 	CreatedAt:    time.Now().Local().UTC(),
		// })

		return c.JSON(re.Builder().Code, re.Builder().Response)
	} else {
		return c.JSON(ErrServerError.Code, ErrServerError.Response)
	}
}
func RespErr(c echo.Context, err error) error {
	response := ErrorResponse{
		Meta: MetaResponse{
			Success: false,
			Message: err.Error(),
		},
	}
	return c.JSON(http.StatusBadGateway, response)
}

func RespValidationErr(c echo.Context, msg string, err interface{}) error {
	response := ValidationErrorResponse{
		Meta: MetaResponse{
			Success: false,
			Message: msg,
		},
		Error: err,
	}
	return c.JSON(http.StatusBadGateway, response)
}

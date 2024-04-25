package server

import (
	"net/http"

	"github.com/HINKOKO/go-microservice/internal/dberrors"
	"github.com/HINKOKO/go-microservice/internal/models"
	"github.com/labstack/echo/v4"
)

func (s *EchoServer) GetAllCustomers(ctx echo.Context) error {
	emailAddress := ctx.QueryParam("emailAddress")

	customers, err := s.DB.GetAllCustomers(ctx.Request().Context(), emailAddress)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusOK, customers)
}

func (s *EchoServer) AddCustomer(ctx echo.Context) error {
	customer := new(models.Customer)
	if err := ctx.Bind(customer); err != nil {
		return ctx.JSON(http.StatusUnsupportedMediaType, err)
	}
	customer, err := s.DB.AddCustomer(ctx.Request().Context(), customer)
	if err != nil {
		switch err.(type) {
		case *dberrors.ConflictError:
			return ctx.JSON(http.StatusConflict, err)
		default:
			return ctx.JSON(http.StatusInternalServerError, err)
		}
	}
	return ctx.JSON(http.StatusCreated, customer)
}

func (s *EchoServer) GetCustomerById(ctx echo.Context) error {
	ID := ctx.Param("id")
	customer, err := s.DB.GetCustomerById(ctx.Request().Context(), ID)

	if err != nil {
		switch err.(type) {
		case *dberrors.NotFoundError:
			return ctx.JSON(http.StatusNotFound, err)
		default:
			return ctx.JSON(http.StatusInternalServerError, err)
		}
	}
	return ctx.JSON(http.StatusOK, customer)
}

func (s *EchoServer) UpdateCustomer(ctx echo.Context) error {
	ID := ctx.Param("id")
	customer := new(models.Customer)
	if err := ctx.Bind(customer); err != nil {
		return ctx.JSON(http.StatusUnsupportedMediaType, err)
	}
	if ID != customer.CustomerID {
		return ctx.JSON(http.StatusBadRequest, "id on path do not match id on body")
	}
	customer, err := s.DB.UpdateCustomer(ctx.Request().Context(), customer)
	if err != nil {
		switch err.(type) {
		case *dberrors.NotFoundError:
			return ctx.JSON(http.StatusNotFound, err)
		case *dberrors.ConflictError:
			return ctx.JSON(http.StatusConflict, err)
		default:
			return ctx.JSON(http.StatusInternalServerError, err)
		}
	}
	return ctx.JSON(http.StatusOK, customer)
}

func (s *EchoServer) DeleteCustomer(ctx echo.Context) error {
	ID := ctx.Param("id")
	err := s.DB.DeleteCustomer(ctx.Request().Context(), ID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	return ctx.NoContent(http.StatusResetContent)
}

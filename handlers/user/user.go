package handlers

import (
	"fmt"
	"net/http"
	"product-svc/helpers/fault"
	"product-svc/model"
	usecases "product-svc/usecases/user"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	user usecases.UserUsecases
}

func NewHandler(usecase usecases.UserUsecases) *Handler {
	return &Handler{
		user: usecase,
	}
}

func (h *Handler) HandleUserRegister(ctx *gin.Context) {
	var body model.RegisterUser

	if err := ctx.ShouldBindJSON(&body); err != nil {
		fault.ErrorHandler(ctx, fault.Custom(
			http.StatusBadRequest,
			fault.ErrBadRequest,
			fmt.Sprintf("failed to bind JSON: %v", err),
		))
		return
	}

	response, err := h.user.UserRegister(body)
	if err != nil {
		fault.ErrorHandler(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, response)
}

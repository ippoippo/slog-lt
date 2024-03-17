package account

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/ippoippo/slog-lt/sldemo"
	"github.com/ippoippo/slog-lt/sldemo/persistence"
)

type Handler struct {
	storage persistence.CrudStorage[*sldemo.Account]
}

func NewHandler() *Handler {
	return &Handler{
		storage: NewStorage(),
	}
}

func (h *Handler) CreateAccount(c echo.Context) error {
	ctx := c.Request().Context()
	req := &sldemo.AccountRequest{}
	if err := c.Bind(req); err != nil {
		slog.ErrorContext(ctx, err.Error())
		return c.JSON(http.StatusBadRequest, sldemo.ErrWithMsg("unable to parse request"))
	}

	accountId, err := uuid.NewV7()
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return c.JSON(http.StatusInternalServerError, sldemo.ErrWithMsg("unable to create account"))
	}

	account := sldemo.AccountFromRequest(accountId, *req)

	slog.DebugContext(ctx, "create account", slog.Any("account", account), slog.Any("h.storage", h.storage))
	err = h.storage.Add(ctx, accountId, account)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return c.JSON(http.StatusConflict, sldemo.ErrWithMsg("unable to create account"))
	}
	slog.InfoContext(ctx, "created account", slog.Any("account", account))
	return c.JSON(http.StatusCreated, account)
}

func (h *Handler) GetAccount(c echo.Context) error {
	ctx := c.Request().Context()
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		msg := fmt.Sprintf("invalid id: [%v]", c.Param("id"))
		return c.JSON(http.StatusBadRequest, sldemo.ErrWithMsg(msg))
	}

	acc, err := h.storage.GetById(ctx, id)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		msg := fmt.Sprintf("account not found: id: [%v]", c.Param("id"))
		return c.JSON(http.StatusNotFound, sldemo.ErrWithMsg(msg))
	}

	return c.JSON(http.StatusOK, acc)
}

func (h *Handler) DeleteAccount(c echo.Context) error {
	ctx := c.Request().Context()
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		msg := fmt.Sprintf("invalid id: [%v]", c.Param("id"))
		return c.JSON(http.StatusBadRequest, sldemo.ErrWithMsg(msg))
	}

	err = h.storage.Delete(ctx, id)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		msg := fmt.Sprintf("cannot delete id: [%v]", c.Param("id"))
		return c.JSON(http.StatusInternalServerError, sldemo.ErrWithMsg(msg))
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *Handler) GetAllAccounts(c echo.Context) error {
	return c.JSON(http.StatusOK, h.storage.GetAll(c.Request().Context()))
}

package account

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"github.com/ippoippo/slog-lt/zapdemo"
	"github.com/ippoippo/slog-lt/zapdemo/persistence"
	"github.com/ippoippo/slog-lt/zapdemo/zapp"
)

type Handler struct {
	storage persistence.CrudStorage[*zapdemo.Account]
}

func NewHandler() *Handler {
	return &Handler{
		storage: NewStorage(),
	}
}

func (h *Handler) CreateAccount(c echo.Context) error {
	ctx := c.Request().Context()
	req := &zapdemo.AccountRequest{}
	if err := c.Bind(req); err != nil {
		zapp.LoggerFromCtx(ctx).Error("unable to bind", zap.Error(err))
		return c.JSON(http.StatusBadRequest, zapdemo.ErrWithMsg("unable to parse request"))
	}

	accountId, err := uuid.NewV7()
	if err != nil {
		zapp.LoggerFromCtx(ctx).Error("unable to create ID for account", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, zapdemo.ErrWithMsg("unable to create account"))
	}

	account := zapdemo.AccountFromRequest(accountId, *req)

	zapp.LoggerFromCtx(ctx).Debug("create account", zap.Any("account", account), zap.Any("h.storage", h.storage))
	err = h.storage.Add(c.Request().Context(), accountId, account)
	if err != nil {
		zapp.LoggerFromCtx(ctx).Error("unable to create account", zap.Error(err))
		return c.JSON(http.StatusConflict, zapdemo.ErrWithMsg("unable to create account"))
	}
	zapp.LoggerFromCtx(ctx).Info("created account", zap.Any("account", account))
	return c.JSON(http.StatusCreated, account)
}

func (h *Handler) GetAccount(c echo.Context) error {
	ctx := c.Request().Context()
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		zapp.LoggerFromCtx(ctx).Error("unable to parse Id", zap.Error(err))
		msg := fmt.Sprintf("invalid id: [%v]", c.Param("id"))
		return c.JSON(http.StatusBadRequest, zapdemo.ErrWithMsg(msg))
	}

	acc, err := h.storage.GetById(ctx, id)
	if err != nil {
		zapp.LoggerFromCtx(ctx).Error("unable to get account", zap.Error(err))
		msg := fmt.Sprintf("account not found: id: [%v]", c.Param("id"))
		return c.JSON(http.StatusNotFound, zapdemo.ErrWithMsg(msg))
	}

	return c.JSON(http.StatusOK, acc)
}

func (h *Handler) DeleteAccount(c echo.Context) error {
	ctx := c.Request().Context()
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		zapp.LoggerFromCtx(ctx).Error("unable to parse Id", zap.Error(err))
		msg := fmt.Sprintf("invalid id: [%v]", c.Param("id"))
		return c.JSON(http.StatusBadRequest, zapdemo.ErrWithMsg(msg))
	}

	err = h.storage.Delete(ctx, id)
	if err != nil {
		zapp.LoggerFromCtx(ctx).Error("unable to delete account", zap.Error(err))
		msg := fmt.Sprintf("cannot delete id: [%v]", c.Param("id"))
		return c.JSON(http.StatusInternalServerError, zapdemo.ErrWithMsg(msg))
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *Handler) GetAllAccounts(c echo.Context) error {
	return c.JSON(http.StatusOK, h.storage.GetAll(c.Request().Context()))
}

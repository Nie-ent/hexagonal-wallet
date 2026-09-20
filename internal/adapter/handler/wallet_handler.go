package handler

import (
	"context"
	"hexagonal/ccgo01/internal/adapter/handler/dto"
	"hexagonal/ccgo01/internal/core/model"
	"hexagonal/ccgo01/internal/core/service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type WalletHandler struct {
	svc *service.WalletService
}

func NewWalletHandler(svc *service.WalletService) *WalletHandler {
	return &WalletHandler{svc: svc}
}

func (h *WalletHandler) CreateWallet(c *gin.Context) {
	var req dto.CreateWalletRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}
	wallet, err := h.svc.CreateWallet(context.Background(), req.UserID)
	if err != nil {
		c.JSON(http.StatusConflict, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, toWalletReponse(wallet))
}

func (h *WalletHandler) GetWallet(c *gin.Context) {
	id := c.Param("id")

	wallet, err := h.svc.GetWallet(context.Background(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, toWalletReponse(wallet))
}

func (h *WalletHandler) Deposit(c *gin.Context) {
	id := c.Param("id")

	var req dto.TransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}
	wallet, err := h.svc.Deposit(context.Background(), id, req.Amount)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, toWalletReponse(wallet))
}

func (h *WalletHandler) Withdraw(c *gin.Context) {
	id := c.Param("id")

	var req dto.TransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}
	wallet, err := h.svc.Withdraw(context.Background(), id, req.Amount)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, toWalletReponse(wallet))
}

func (h *WalletHandler) DeleteWallet(c *gin.Context) {
	id := c.Param("id")

	if err := h.svc.DeleteWalelt(context.Background(), id); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

func toWalletReponse(w *model.Wallet) dto.WalletResponse {
	return dto.WalletResponse{
		ID:        w.ID,
		UserID:    w.UserID,
		Balance:   w.Balance,
		CreatedAt: w.CreatedAt.Format(time.RFC3339),
		UpdatedAt: w.UpdatedAt.Format(time.RFC3339),
	}
}

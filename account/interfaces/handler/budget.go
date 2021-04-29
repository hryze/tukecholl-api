package handler

import (
	"context"

	"github.com/paypay3/tukecholl-api/account/usecase"
	"github.com/paypay3/tukecholl-api/account/usecase/input"
	"github.com/paypay3/tukecholl-api/proto/accountproto"
)

type budgetHandler struct {
	budgetUsecase usecase.BudgetUsecase
	accountproto.UnimplementedBudgetServiceServer
}

func NewBudgetHandler(budgetUsecase usecase.BudgetUsecase) *budgetHandler {
	return &budgetHandler{
		budgetUsecase: budgetUsecase,
	}
}

func (h *budgetHandler) CreateStandardBudgets(ctx context.Context, r *accountproto.CreateStandardBudgetsRequest) (*accountproto.CreateStandardBudgetsResponse, error) {
	user := &input.User{ID: r.GetUserId()}

	if err := h.budgetUsecase.CreateStandardBudgets(user); err != nil {
		return nil, err
	}

	return &accountproto.CreateStandardBudgetsResponse{}, nil
}

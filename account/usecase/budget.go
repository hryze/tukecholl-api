package usecase

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/paypay3/tukecholl-api/account/domain/budgetdomain"
	"github.com/paypay3/tukecholl-api/account/domain/vo"
	"github.com/paypay3/tukecholl-api/account/usecase/input"
)

type BudgetUsecase interface {
	CreateStandardBudgets(user *input.User) error
}

type budgetUsecase struct {
	budgetRepository budgetdomain.Repository
}

func NewBudgetUsecase(budgetRepository budgetdomain.Repository) *budgetUsecase {
	return &budgetUsecase{
		budgetRepository: budgetRepository,
	}
}

func (u *budgetUsecase) CreateStandardBudgets(user *input.User) error {
	userID, err := vo.NewUserID(user.ID)
	if err != nil {
		return status.Errorf(codes.InvalidArgument, "invalid user id: %v", err)
	}

	if err := u.budgetRepository.CreateStandardBudgets(userID); err != nil {
		return err
	}

	return nil
}

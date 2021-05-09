package usecase

import (
	"github.com/paypay3/tukecholl-api/account/domain/budgetdomain"
	"github.com/paypay3/tukecholl-api/account/domain/vo"
	"github.com/paypay3/tukecholl-api/account/usecase/input"
	"github.com/paypay3/tukecholl-api/pkg/apperrors"
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
		return apperrors.InvalidParameter.AddBadRequestFieldViolation("user id", "ユーザーIDを正しく指定してください").Wrap(err, "invalid user id")
	}

	if err := u.budgetRepository.CreateStandardBudgets(userID); err != nil {
		return apperrors.Wrap(err, "failed to create standard budget initial value")
	}

	return nil
}

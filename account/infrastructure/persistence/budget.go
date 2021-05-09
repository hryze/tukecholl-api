package persistence

import (
	"github.com/paypay3/tukecholl-api/account/domain/vo"
	"github.com/paypay3/tukecholl-api/account/infrastructure/persistence/rdb"
	"github.com/paypay3/tukecholl-api/pkg/apperrors"
)

type budgetRepository struct {
	*rdb.Driver
}

func NewBudgetRepository(rdbDriver *rdb.Driver) *budgetRepository {
	return &budgetRepository{rdbDriver}
}

func (r *budgetRepository) CreateStandardBudgets(userID vo.UserID) error {
	query := `
        INSERT INTO standard_budgets
            (user_id, big_category_id)
        VALUES
            (?,2),
            (?,3),
            (?,4),
            (?,5),
            (?,6),
            (?,7),
            (?,8),
            (?,9),
            (?,10),
            (?,11),
            (?,12),
            (?,13),
            (?,14),
            (?,15),
            (?,16),
            (?,17)`

	if _, err := r.Driver.Conn.Exec(query, userID, userID, userID, userID, userID, userID, userID, userID, userID, userID, userID, userID, userID, userID, userID, userID); err != nil {
		return apperrors.InternalServerError.SetMessage("標準予算の初期値追加に失敗しました").Wrap(err, "rdb unexpected error")
	}

	return nil
}

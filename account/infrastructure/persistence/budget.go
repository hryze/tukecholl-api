package persistence

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/paypay3/tukecholl-api/account/domain/vo"
	"github.com/paypay3/tukecholl-api/account/infrastructure/persistence/rdb"
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
		return status.Errorf(codes.Internal, "rdb unexpected error: %v", err)
	}

	return nil
}

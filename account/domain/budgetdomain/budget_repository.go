package budgetdomain

import "github.com/paypay3/tukecholl-api/account/domain/vo"

type Repository interface {
	CreateStandardBudgets(userID vo.UserID) error
}

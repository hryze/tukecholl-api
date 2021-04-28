package server

import (
	"google.golang.org/grpc"

	"github.com/paypay3/tukecholl-api/account/infrastructure/persistence"
	"github.com/paypay3/tukecholl-api/account/infrastructure/persistence/rdb"
	"github.com/paypay3/tukecholl-api/account/interfaces/handler"
	"github.com/paypay3/tukecholl-api/account/usecase"
	"github.com/paypay3/tukecholl-api/proto/accountproto"
)

func registerBudgetServiceServer(srv *grpc.Server, rdbDriver *rdb.Driver) {
	budgetRepository := persistence.NewBudgetRepository(rdbDriver)
	budgetUsecase := usecase.NewBudgetUsecase(budgetRepository)
	budgetHandler := handler.NewBudgetHandler(budgetUsecase)

	accountproto.RegisterBudgetServiceServer(srv, budgetHandler)
}

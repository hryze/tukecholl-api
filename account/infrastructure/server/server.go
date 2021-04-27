package server

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/paypay3/tukecholl-api/account/config"
	"github.com/paypay3/tukecholl-api/account/infrastructure/persistence"
	"github.com/paypay3/tukecholl-api/account/infrastructure/persistence/rdb"
	"github.com/paypay3/tukecholl-api/account/interfaces/handler"
	"github.com/paypay3/tukecholl-api/account/usecase"
	"github.com/paypay3/tukecholl-api/proto/accountproto"
)

func Run() error {
	rdbDriver, err := rdb.NewDriver()
	if err != nil {
		return err
	}
	defer rdbDriver.Conn.Close()

	budgetRepository := persistence.NewBudgetRepository(rdbDriver)
	budgetUsecase := usecase.NewBudgetUsecase(budgetRepository)
	budgetHandler := handler.NewBudgetHandler(budgetUsecase)

	srv := grpc.NewServer()
	reflection.Register(srv)
	accountproto.RegisterBudgetServiceServer(srv, budgetHandler)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", config.Env.Server.Port))
	if err != nil {
		return err
	}
	defer lis.Close()

	errorCh := make(chan error, 1)
	go func() {
		if err := srv.Serve(lis); err != nil {
			errorCh <- err
		}
	}()

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGTERM, syscall.SIGINT)

	select {
	case err := <-errorCh:
		return err
	case s := <-signalCh:
		log.Printf("SIGNAL %s received", s.String())
		srv.GracefulStop()
	}

	return nil
}

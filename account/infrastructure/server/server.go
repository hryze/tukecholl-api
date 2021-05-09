package server

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/paypay3/tukecholl-api/account/config"
	"github.com/paypay3/tukecholl-api/account/infrastructure/persistence/rdb"
	"github.com/paypay3/tukecholl-api/pkg/Interceptor"
)

func Run() error {
	rdbDriver, err := rdb.NewDriver()
	if err != nil {
		return err
	}
	defer rdbDriver.Conn.Close()

	srv := grpc.NewServer(
		grpc.UnaryInterceptor(grpcMiddleware.ChainUnaryServer(
			Interceptor.TransmitErrorWithStatus(),
			Interceptor.Logging(),
		)),
	)

	// register services to the server.
	reflection.Register(srv)
	registerBudgetServiceServer(srv, rdbDriver)

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

package service

import (
	"errors"
	"time"

	healthcheck "github.com/sigmaott/gest/package/technique/health"

	"go.uber.org/zap"
)

type IHeathCheckService interface {
	HeathCheckPayment() (*healthcheck.Response, error)
}

type heathCheckService struct {
	logger *zap.SugaredLogger
}

func (h *heathCheckService) HeathCheckPayment() (*healthcheck.Response, error) {
	res := healthcheck.HandlerHeathCheck(
		healthcheck.WithTimeout(10 * time.Second),
		// healthcheck.WithChecker("mongodb", healthcheck.CheckerFunc(
		// 	func(ctx context.Context) error {
		// 		err := h.database.Client().Ping(ctx, nil)
		// 		if err != nil {
		// 			h.logger.Errorln(err)
		// 		}
		// 		return err
		// 	},
		// )),
		// healthcheck.WithChecker("lago", healthcheck.CheckerFunc(
		// 	func(ctx context.Context) error {
		// 		_, err := h.lago.Webhook().GetPublicKey(ctx)
		// 		if err != nil {
		// 			h.logger.Errorf("Lago service get error: %v", err)
		// 			return errors.New(fmt.Sprintf("Lago service err: %d", err.HTTPStatusCode))
		// 		}
		// 		return nil
		// 	},
		// )),
		// healthcheck.WithChecker("rabbitmq", healthcheck.CheckerFunc(
		// 	func(ctx context.Context) error {
		// 		if h.rabbitmq.IsClosed() {
		// 			return errors.New("Rabbitmq is not working")
		// 		}
		// 		return nil
		// 	},
		// )),
		// healthcheck.WithChecker("app-grpc", healthcheck.CheckerFunc(
		// 	func(ctx context.Context) error {
		// 		conn, err := grpc.Dial(config.GetConfiguration().Dependencies.AuthService.Uri, grpc.WithInsecure())
		// 		if err != nil {
		// 			h.logger.Errorf("Failed to connect auth service: %v", err)
		// 		}
		// 		defer func(conn *grpc.ClientConn) {
		// 			conn.Close()
		// 			h.logger.Infoln("Closed")
		// 		}(conn)
		// 		state := conn.GetState()
		// 		log.Print(state)
		// 		if string(rune(state)) != "READY" {
		// 			h.logger.Errorf("App-grpc connect err with status: %s", state)
		// 			return errors.New(fmt.Sprintf("App-grpc connect err with status: %s", state))
		// 		}
		// 		return nil
		// 	},
		// )),
		// healthcheck.WithChecker("quota-grpc", healthcheck.CheckerFunc(
		// 	func(ctx context.Context) error {
		// 		conn, err := grpc.Dial(config.GetConfiguration().Dependencies.QuotaService.Uri, grpc.WithInsecure())
		// 		if err != nil {
		// 			return errors.New(fmt.Sprintf("Failed to connect quota service: %v", err))
		// 		}
		// 		defer func(conn *grpc.ClientConn) {
		// 			conn.Close()
		// 			log.Print("Closed")
		// 		}(conn)
		// 		state := conn.GetState()
		// 		log.Print(state)
		// 		if string(rune(state)) != "READY" {
		// 			h.logger.Errorf("Quota-grpc connect err with status: %s", state)
		// 			return errors.New(fmt.Sprintf("Quota-grpc connect err with status: %s", state))
		// 		}
		// 		return nil
		// 	},
		// )),
		// healthcheck.WithChecker("redis", healthcheck.CheckerFunc(
		// 	func(ctx context.Context) error {
		// 		task := asynq.NewTask("health", []byte("health"))
		// 		_, err := h.asyncq.Enqueue(task)
		// 		if err != nil {
		// 			h.logger.Errorf("Redis connection err: %v", err)
		// 			return err
		// 		}
		// 		return nil
		// 	},
		// )),
	)
	if len(res.Errors) > 0 {
		return res, errors.New("Service Unavailable")
	}
	return res, nil
}

func NewHeathCheckService(
	logger *zap.SugaredLogger) IHeathCheckService {
	return &heathCheckService{
		logger: logger,
	}
}

package grpcapi

import (
	"context"
	"fmt"
	"log"

	"github-release-notifier/internal/repository"
)

// GrpcHandler реалізує згенерований інтерфейс NotifierServiceServer
type GrpcHandler struct {
	UnimplementedNotifierServiceServer
	dbRepo *repository.Repository
}

// NewGrpcHandler створює новий екземпляр хендлера
func NewGrpcHandler(repo *repository.Repository) *GrpcHandler {
	return &GrpcHandler{dbRepo: repo}
}

// Subscribe обробляє вхідні gRPC запити на підписку
func (h *GrpcHandler) Subscribe(ctx context.Context, req *SubscribeRequest) (*SubscribeResponse, error) {
	log.Printf("Отримано gRPC запит: email=%s, repo=%s", req.Email, req.Repository)
	err := h.dbRepo.SubscribeUser(req.Email, req.Repository, "")
	if err != nil {
		log.Printf("Помилка збереження через gRPC: %v", err)
		return nil, fmt.Errorf("помилка бази даних: %v", err)
	}

	return &SubscribeResponse{
		Message: "Успішно підписано на " + req.Repository + " через gRPC!",
	}, nil
}

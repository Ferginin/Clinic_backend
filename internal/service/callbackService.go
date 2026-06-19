package service

import (
	"Clinic_backend/internal/entity"
	"Clinic_backend/internal/repository"
	"context"
	"errors"
	"regexp"
)

type CallbackServiceInterface interface {
	CreateCallbackRequest(ctx context.Context, req *entity.CallbackRequestCreate) (*entity.CallbackRequest, error)
	GetAllCallbackRequests(ctx context.Context) ([]entity.CallbackRequest, error)
	GetAllCallbackRequestsPaginated(ctx context.Context, limit, offset int) ([]entity.CallbackRequest, int, error)
	GetCallbackRequestByID(ctx context.Context, id int) (*entity.CallbackRequest, error)
	UpdateCallbackRequest(ctx context.Context, id int, req *entity.CallbackRequestUpdate) (*entity.CallbackRequest, error)
	DeleteCallbackRequest(ctx context.Context, id int) error
}

type CallbackService struct {
	callbackRepo repository.CallbackRepositoryInterface
}

func NewCallbackService(cr repository.CallbackRepositoryInterface) CallbackServiceInterface {
	return &CallbackService{callbackRepo: cr}
}

func (s *CallbackService) CreateCallbackRequest(ctx context.Context, req *entity.CallbackRequestCreate) (*entity.CallbackRequest, error) {
	phoneRegex := regexp.MustCompile(`^\+?[0-9\s\-\(\)]+$`)
	if !phoneRegex.MatchString(req.Phone) {
		return nil, errors.New("invalid phone number format")
	}

	if len(req.Name) < 2 {
		return nil, errors.New("name must be at least 2 characters")
	}

	callbackReq := &entity.CallbackRequest{
		Name:    req.Name,
		Phone:   req.Phone,
		Message: req.Message,
		Status:  "new",
	}

	return s.callbackRepo.Create(ctx, callbackReq)
}

func (s *CallbackService) GetAllCallbackRequests(ctx context.Context) ([]entity.CallbackRequest, error) {
	return s.callbackRepo.GetAll(ctx)
}

func (s *CallbackService) GetAllCallbackRequestsPaginated(ctx context.Context, limit, offset int) ([]entity.CallbackRequest, int, error) {
	return s.callbackRepo.GetAllPaginated(ctx, limit, offset)
}

func (s *CallbackService) GetCallbackRequestByID(ctx context.Context, id int) (*entity.CallbackRequest, error) {
	return s.callbackRepo.GetByID(ctx, id)
}

func (s *CallbackService) UpdateCallbackRequest(ctx context.Context, id int, req *entity.CallbackRequestUpdate) (*entity.CallbackRequest, error) {
	existing, err := s.callbackRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Status != nil {
		existing.Status = *req.Status
	}

	return s.callbackRepo.Update(ctx, id, existing)
}

func (s *CallbackService) DeleteCallbackRequest(ctx context.Context, id int) error {
	return s.callbackRepo.Delete(ctx, id)
}

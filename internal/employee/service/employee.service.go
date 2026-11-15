package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"eci-backend/internal/employee/model"
	"eci-backend/internal/employee/repository"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type ActivityEvent struct {
	Action   string
	ActorID  string
	EntityID string
	Before   any
	After    any
}

type ActivityLogger interface {
	Record(ctx context.Context, event ActivityEvent) error
}

type EmployeeFilter struct {
	Search   string
	Grade    string
	MinAge   int
	MaxAge   int
	Province string
	City     string
	District string
}

type EmployeeService struct {
	repo   *repository.EmployeeRepository
	logger ActivityLogger
}

func NewEmployeeService(
	repo *repository.EmployeeRepository,
	logger ActivityLogger,
) *EmployeeService {
	return &EmployeeService{
		repo:   repo,
		logger: logger,
	}
}

func (s *EmployeeService) Create(
	ctx context.Context,
	employee *model.Employee,
	actorID string,
) error {
	if err := validateEmployee(employee); err != nil {
		return err
	}

	if err := s.repo.Create(ctx, employee); err != nil {
		return err
	}

	s.recordActivity(ctx, ActivityEvent{
		Action:   "CREATE_EMPLOYEE",
		ActorID:  actorID,
		EntityID: employee.ID.Hex(),
		After:    employee,
	})

	return nil
}

func (s *EmployeeService) FindAll(
	ctx context.Context,
	filter EmployeeFilter,
	page int64,
	limit int64,
) ([]model.Employee, int64, error) {
	mongoFilter := bson.M{}

	if filter.Search != "" {
		search := strings.TrimSpace(filter.Search)
		mongoFilter["name"] = bson.M{
			"$regex":   search,
			"$options": "i",
		}
	}

	if filter.Grade != "" {
		mongoFilter["grade"] = filter.Grade
	}

	if filter.Province != "" {
		mongoFilter["province"] = bson.M{
			"$regex":   filter.Province,
			"$options": "i",
		}
	}

	if filter.City != "" {
		mongoFilter["city"] = bson.M{
			"$regex":   filter.City,
			"$options": "i",
		}
	}

	if filter.District != "" {
		mongoFilter["district"] = bson.M{
			"$regex":   filter.District,
			"$options": "i",
		}
	}

	if filter.MinAge > 0 || filter.MaxAge > 0 {
		birthDateFilter := bson.M{}
		now := time.Now()

		if filter.MinAge > 0 {
			birthDateFilter["$lte"] = now.AddDate(-filter.MinAge, 0, 0)
		}

		if filter.MaxAge > 0 {
			birthDateFilter["$gte"] = now.AddDate(-(filter.MaxAge + 1), 0, 1)
		}

		mongoFilter["birthDate"] = birthDateFilter
	}

	return s.repo.FindAll(ctx, mongoFilter, page, limit)
}

func (s *EmployeeService) FindByID(
	ctx context.Context,
	id string,
) (*model.Employee, error) {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid employee ID")
	}

	return s.repo.FindByID(ctx, objectID)
}

func (s *EmployeeService) Update(
	ctx context.Context,
	id string,
	employee *model.Employee,
	actorID string,
) error {
	if err := validateEmployee(employee); err != nil {
		return err
	}

	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid employee ID")
	}

	before, err := s.repo.FindByID(ctx, objectID)
	if err != nil {
		return err
	}

	if err := s.repo.Update(ctx, objectID, employee); err != nil {
		return err
	}

	employee.ID = objectID

	s.recordActivity(ctx, ActivityEvent{
		Action:   "UPDATE_EMPLOYEE",
		ActorID:  actorID,
		EntityID: id,
		Before:   before,
		After:    employee,
	})

	return nil
}

func (s *EmployeeService) Delete(
	ctx context.Context,
	id string,
	actorID string,
) error {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid employee ID")
	}

	before, err := s.repo.FindByID(ctx, objectID)
	if err != nil {
		return err
	}

	if err := s.repo.Delete(ctx, objectID); err != nil {
		return err
	}

	s.recordActivity(ctx, ActivityEvent{
		Action:   "DELETE_EMPLOYEE",
		ActorID:  actorID,
		EntityID: id,
		Before:   before,
	})

	return nil
}

func (s *EmployeeService) recordActivity(
	ctx context.Context,
	event ActivityEvent,
) {
	if s.logger == nil {
		return
	}

	if err := s.logger.Record(ctx, event); err != nil {
		fmt.Printf("failed to record activity: %v\n", err)
	}
}

func validateEmployee(employee *model.Employee) error {
	if employee == nil {
		return errors.New("employee data is required")
	}

	if strings.TrimSpace(employee.Name) == "" {
		return errors.New("employee name is required")
	}

	if strings.TrimSpace(employee.Position) == "" {
		return errors.New("employee position is required")
	}

	if strings.TrimSpace(employee.Grade) == "" {
		return errors.New("employee grade is required")
	}

	if employee.Salary <= 0 {
		return errors.New("salary must be greater than zero")
	}

	if strings.TrimSpace(employee.Address) == "" {
		return errors.New("address is required")
	}

	if strings.TrimSpace(employee.Province) == "" ||
		strings.TrimSpace(employee.City) == "" ||
		strings.TrimSpace(employee.District) == "" {
		return errors.New("complete employee location is required")
	}

	if employee.BirthDate.IsZero() || !employee.BirthDate.Before(time.Now()) {
		return errors.New("valid birth date is required")
	}

	return nil
}

package usecase

import (
	"github.com/ma-miyazaki/go-grpc-neo4j-example/server/domain/model"
	"github.com/ma-miyazaki/go-grpc-neo4j-example/server/domain/repository"
	"github.com/rs/zerolog/log"
)

type EmployeeUseCase interface {
	AddEmployee(email string, lastName string, firstName string) (*model.Employee, error)
	ListEmployees() ([]*model.Employee, error)
}

type employeeUseCase struct {
	repository repository.EmployeeRepository
}

func NewEmployeeUseCase(repository repository.EmployeeRepository) EmployeeUseCase {
	return employeeUseCase{repository}
}

func (uc employeeUseCase) AddEmployee(email string, lastName string, firstName string) (*model.Employee, error) {
	defer uc.repository.Close()
	if err := uc.repository.Begin(); err != nil {
		return nil, err
	}

	employee, err := model.NewEmployee(email, lastName, firstName)
	if err != nil {
		log.Fatal().Stack().Err(err).Msg("failed to create employee model")
		return nil, err
	}

	if err := uc.repository.Create(employee); err != nil {
		log.Fatal().Stack().Err(err).Msg("failed to store employee")
		return nil, err
	}

	log.Info().Msgf("Employee created. [%v]", employee)
	uc.repository.Commit()
	return employee, nil
}

func (uc employeeUseCase) ListEmployees() ([]*model.Employee, error) {
	log.Info().Msgf("repository: %p", &uc.repository)
	return uc.repository.List()
}

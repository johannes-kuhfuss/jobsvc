package repositories

import (
	"github.com/johannes-kuhfuss/jobsvc/domain"
	"github.com/johannes-kuhfuss/jobsvc/dto"
	"github.com/johannes-kuhfuss/services_utils/api_error"
)

type JobRepositoryDb struct {
}

func NewJobRepositoryDb() JobRepositoryDb {
	return JobRepositoryDb{}
}

func (jrd JobRepositoryDb) FindAll(status string) (*[]domain.Job, api_error.ApiErr) {
	panic("not implemented")
}

func (jrd JobRepositoryDb) FindById(id string) (*domain.Job, api_error.ApiErr) {
	panic("not implemented")
}

func (jrd JobRepositoryDb) Store(job domain.Job) api_error.ApiErr {
	panic("not implemented")
}

func (jrd JobRepositoryDb) DeleteById(id string) api_error.ApiErr {
	panic("not implemented")
}

func (jrd JobRepositoryDb) GetNext() (*domain.Job, api_error.ApiErr) {
	panic("not implemented")
}

func (jrd JobRepositoryDb) SetStatus(id string, newStatus dto.UpdateJobStatusRequest) api_error.ApiErr {
	panic("not implemented")
}

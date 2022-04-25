package customer

// NewService is a constructor to get a Service instance
func NewService(repo Repo) Service {
	return &service{
		repo: repo,
	}
}

// service is a struct of service
type service struct {
	repo Repo
}

// Service is used to define the methods in it
type Service interface {
	RetrieveCustomerProfile(userID int) (*Profile, error)
}

// RetrieveCustomerProfile is called to get customer profile through repository
func (s service) RetrieveCustomerProfile(userID int) (*Profile, error) {

	customerProfile, err := s.repo.RetrieveCustomerProfile(userID)

	if err != nil {
		return nil, err
	}

	return customerProfile, nil
}

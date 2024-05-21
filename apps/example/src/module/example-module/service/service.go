package service

type IExampleService interface {
}
type exampleService struct {
}

func NewExampleService() IExampleService {
	return &exampleService{}

}

package repository

import (
	"example/src/module/example-module/model"

	base "github.com/sigmaott/gest/package/technique/database/base"
	mongoRepository "github.com/sigmaott/gest/package/technique/database/mongo"
	"go.mongodb.org/mongo-driver/mongo"
)

type IExampleRepository interface {
	base.IRepository[model.Example]
}
type exampleRepository struct {
	mongoRepository.BaseMongoRepository[model.Example]
}

func NewExampleRepository(db *mongo.Database) IExampleRepository {
	return &exampleRepository{
		BaseMongoRepository: mongoRepository.BaseMongoRepository[model.Example]{
			Collection: db.Collection("examples"),
		},
	}
}

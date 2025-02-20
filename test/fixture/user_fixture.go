package fixture

import (
	"app/src/model"

	"github.com/google/uuid"
)

var UserOne = &model.User{
	BaseModel: model.BaseModel{
		ID: uuid.New(),
	},
	Name:     "Test1",
	Email:    "test1@gmail.com",
	Password: "password1",
	Role:     "user",
}

var UserTwo = &model.User{
	BaseModel: model.BaseModel{
		ID: uuid.New(),
	},
	Name:     "Test2",
	Email:    "test2@gmail.com",
	Password: "password1",
	Role:     "user",
}

var Admin = &model.User{
	BaseModel: model.BaseModel{
		ID: uuid.New(),
	},
	Name:     "Admin",
	Email:    "admin@gmail.com",
	Password: "password1",
	Role:     "admin",
}

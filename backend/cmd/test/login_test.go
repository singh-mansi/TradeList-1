package test

import (
	"testing"
	"tradelist/mocks"
	"tradelist/pkg/api"

	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"
)

func getTestUser() api.User {
	var contact = api.Contact{
		FirstName: "Test",
		LastName:  "Surname",
		Email:     "test@gmail.com",
		Password:  "test",
		Phone:     "",
	}

	var user = api.User{
		ID:       0,
		IsSeller: false,
		Seller:   api.Seller{},
		Contact:  contact,
		Token:    "",
	}
	return user
}

func TestSignup_Email_Exist(test *testing.T) {
	repo := mocks.NewMockRepo(test)
	user := getTestUser()

	loginService := api.CreateLoginService(repo)

	repo.On("IsEmailExisting", user.Contact.Email).Return(true)

	response := loginService.SignUp(user)
	assert.Equal(test, 401, response["status"])
	assert.Equal(test, "Email already exists", response["message"])
}

func TestSignup_Error(test *testing.T) {
	repo := mocks.NewMockRepo(test)
	user := getTestUser()

	loginService := api.CreateLoginService(repo)

	repo.On("IsEmailExisting", user.Contact.Email).Return(false)
	repo.On("CreateUser", user).Return(user, "Error")

	response := loginService.SignUp(user)
	assert.Equal(test, 0, response["status"])
	assert.Equal(test, "Error", response["message"])
}

func TestSignup_Success(test *testing.T) {
	repo := mocks.NewMockRepo(test)
	user := getTestUser()

	loginService := api.CreateLoginService(repo)

	repo.On("IsEmailExisting", user.Contact.Email).Return(false)
	repo.On("CreateUser", user).Return(user, "")

	response := loginService.SignUp(user)
	assert.Equal(test, 201, response["status"])
	assert.Equal(test, "User created successfully", response["message"])

	var resultUser api.User
	mapstructure.Decode(response["data"], &resultUser)

	assert.Equal(test, "", resultUser.Contact.Password)
}

func TestLogin_User_Exist(test *testing.T) {
	repo := mocks.NewMockRepo(test)
	user := getTestUser()

	loginService := api.CreateLoginService(repo)

	repo.On("FetchUserInfo", user.Contact.Email).Return(user, "Error")

	user, response := loginService.FetchUserInfo(user.Contact.Email)
	assert.Equal(test, 404, response["status"])
	assert.Equal(test, "User not found", response["message"])

}

func TestLogin_Success(test *testing.T) {
	repo := mocks.NewMockRepo(test)
	user := getTestUser()

	loginService := api.CreateLoginService(repo)

	repo.On("FetchUserInfo", user.Contact.Email).Return(user, "")

	user, response := loginService.FetchUserInfo(user.Contact.Email)
	assert.Equal(test, 201, response["status"])
	assert.Equal(test, "User found", response["message"])

	var resultUser api.User
	mapstructure.Decode(response["data"], &resultUser)

	assert.Equal(test, "", resultUser.Contact.Password)

}
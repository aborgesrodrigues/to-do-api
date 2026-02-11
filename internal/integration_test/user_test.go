//go:build service
// +build service

package integrationtest

import (
	"net/http"

	"github.com/aborgesrodrigues/to-do-api/internal/common"
	"github.com/google/uuid"
)

func (s *testSuite) listUsers() []common.User {
	users := make([]common.User, 0)
	s.call("GET", "http://localhost:8080/users", nil, &users)

	return users
}

func (s *testSuite) getUser(id string) *common.User {
	user := &common.User{}
	s.call("GET", "http://localhost:8080/users/"+id, nil, user)

	return user
}

func (s *testSuite) getLastUser() *common.User {
	users := s.listUsers()
	if len(users) == 0 {
		user := &common.User{
			Username: "username1",
			Name:     "Name 1",
		}

		s.call("POST", "http://localhost:8080/users/", user, user)
		users = []common.User{*user}
	}

	return &users[len(users)-1]
}

func (s *testSuite) TestAddUser() {
	// check number of users before addint
	oldNumberUsers := len(s.listUsers())

	// add a new user
	user := &common.User{
		Username: "username1" + uuid.New().String(),
		Name:     "User 1" + uuid.New().String(),
	}

	newUser := &common.AuthResponse{}
	res := s.call("POST", "http://localhost:8080/users", user, newUser)

	s.Assert().Equal(http.StatusCreated, res.StatusCode)
	s.Assert().Equal(user.Name, newUser.User.Name)
	s.Assert().Equal(user.Username, newUser.User.Username)

	// check number of users after add user
	newNumberUsers := len(s.listUsers())
	s.Assert().Equal(oldNumberUsers+1, newNumberUsers)
}

func (s *testSuite) TestUpdateUser() {
	lastUser := s.getLastUser()
	newName := "New User Name1" + uuid.New().String()
	newUsername := "newuser1"
	// check user data before update
	s.Assert().NotEqual(lastUser.Name, newName)
	s.Assert().NotEqual(lastUser.Username, newUsername)

	// add a new user
	user := &common.User{
		Username: newName,
		Name:     newUsername,
	}

	res := s.call("PUT", "http://localhost:8080/users/"+lastUser.Id, user, lastUser)

	s.Assert().Equal(http.StatusOK, res.StatusCode)
	s.Assert().Equal(user.Name, lastUser.Name)
	s.Assert().Equal(user.Username, lastUser.Username)

	// check user data after update
	newUser := s.getUser(lastUser.Id)
	s.Assert().Equal(lastUser, newUser)

}

func (s *testSuite) TestGetUser() {
	lastUser := s.getLastUser()
	user := &common.User{}

	res := s.call("GET", "http://localhost:8080/users/"+lastUser.Id, nil, user)

	s.Assert().Equal(http.StatusOK, res.StatusCode)
	s.Assert().Equal(lastUser, user)
}

func (s *testSuite) TestDeleteUser() {
	lastUser := s.getLastUser()
	// check number of users before addint
	oldNumberUsers := len(s.listUsers())

	res := s.call("DELETE", "http://localhost:8080/users/"+lastUser.Id, nil, nil)

	s.Assert().Equal(http.StatusOK, res.StatusCode)

	// check number of users after add user
	newNumberUsers := len(s.listUsers())
	s.Assert().Equal(oldNumberUsers-1, newNumberUsers)

}

package integrationtest

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/aborgesrodrigues/to-do-api/internal/common"
	"github.com/google/uuid"
)

func (s *testSuite) listUsers() []common.User {
	users := make([]common.User, 0)
	req, err := http.NewRequest("GET", "http://localhost:8080/users", nil)
	s.Assert().NoError(err)

	res, err := http.DefaultClient.Do(req)
	s.Assert().NoError(err)

	err = json.NewDecoder(res.Body).Decode(&users)
	s.Assert().NoError(err)

	return users
}

func (s *testSuite) getUser(id string) *common.User {
	user := &common.User{}
	req, err := http.NewRequest("GET", "http://localhost:8080/users/"+id, nil)
	s.Assert().NoError(err)

	res, err := http.DefaultClient.Do(req)
	s.Assert().NoError(err)

	err = json.NewDecoder(res.Body).Decode(&user)
	s.Assert().NoError(err)

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
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(user)
	s.Assert().NoError(err)

	req, err := http.NewRequest("POST", "http://localhost:8080/users", ioutil.NopCloser(&buf))
	s.Assert().NoError(err)

	res, err := http.DefaultClient.Do(req)
	s.Assert().NoError(err)

	newUser := &common.User{}
	// set new user to lastUser variable to use in other tests
	err = json.NewDecoder(res.Body).Decode(&newUser)
	s.Assert().NoError(err)

	s.Assert().Equal(http.StatusCreated, res.StatusCode)
	s.Assert().Equal(user.Name, newUser.Name)
	s.Assert().Equal(user.Username, newUser.Username)

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
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(user)
	s.Assert().NoError(err)

	req, err := http.NewRequest("PUT", "http://localhost:8080/users/"+lastUser.Id, ioutil.NopCloser(&buf))
	s.Assert().NoError(err)

	res, err := http.DefaultClient.Do(req)
	s.Assert().NoError(err)

	err = json.NewDecoder(res.Body).Decode(&lastUser)
	s.Assert().NoError(err)

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

	req, err := http.NewRequest("GET", "http://localhost:8080/users/"+lastUser.Id, nil)
	s.Assert().NoError(err)

	res, err := http.DefaultClient.Do(req)
	s.Assert().NoError(err)

	err = json.NewDecoder(res.Body).Decode(&user)
	s.Assert().NoError(err)

	s.Assert().Equal(http.StatusOK, res.StatusCode)
	s.Assert().Equal(lastUser, user)
}

func (s *testSuite) TestDeleteUser() {
	lastUser := s.getLastUser()
	// check number of users before addint
	oldNumberUsers := len(s.listUsers())

	req, err := http.NewRequest("DELETE", "http://localhost:8080/users/"+lastUser.Id, nil)
	s.Assert().NoError(err)

	res, err := http.DefaultClient.Do(req)
	s.Assert().NoError(err)

	s.Assert().Equal(http.StatusOK, res.StatusCode)

	// check number of users after add user
	newNumberUsers := len(s.listUsers())
	s.Assert().Equal(oldNumberUsers-1, newNumberUsers)

}

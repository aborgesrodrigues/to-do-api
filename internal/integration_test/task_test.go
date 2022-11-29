package integrationtest

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/aborgesrodrigues/to-do-api/internal/common"
	"github.com/google/uuid"
)

func (s *testSuite) listTasks() []common.Task {
	tasks := make([]common.Task, 0)
	req, err := http.NewRequest("GET", "http://localhost:8080/tasks", nil)
	s.Assert().NoError(err)

	res, err := http.DefaultClient.Do(req)
	s.Assert().NoError(err)

	err = json.NewDecoder(res.Body).Decode(&tasks)
	s.Assert().NoError(err)

	return tasks
}

func (s *testSuite) getTask(id string) *common.Task {
	task := &common.Task{}
	req, err := http.NewRequest("GET", "http://localhost:8080/tasks/"+id, nil)
	s.Assert().NoError(err)

	res, err := http.DefaultClient.Do(req)
	s.Assert().NoError(err)

	err = json.NewDecoder(res.Body).Decode(&task)
	s.Assert().NoError(err)

	return task
}

func (s *testSuite) getLastTask() *common.Task {
	tasks := s.listTasks()
	if len(tasks) == 0 {
		task := &common.Task{
			Description: "description 1",
			State:       "to_do",
			UserId:      s.getLastUser().Id,
		}

		s.call("POST", "http://localhost:8080/tasks/", task, task)
		tasks = []common.Task{*task}
	}

	return &tasks[len(tasks)-1]
}

func (s *testSuite) TestAddTask() {
	// check number of tasks before addint
	oldNumberTasks := len(s.listTasks())

	// add a new task
	task := &common.Task{
		Description: "description 1" + uuid.New().String(),
		State:       "to_do",
		UserId:      s.getLastUser().Id,
	}
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(task)
	s.Assert().NoError(err)

	req, err := http.NewRequest("POST", "http://localhost:8080/tasks", ioutil.NopCloser(&buf))
	s.Assert().NoError(err)

	res, err := http.DefaultClient.Do(req)
	s.Assert().NoError(err)

	newTask := &common.Task{}
	// set new task to lastTask variable to use in other tests
	err = json.NewDecoder(res.Body).Decode(&newTask)
	s.Assert().NoError(err)

	s.Assert().Equal(http.StatusCreated, res.StatusCode)
	s.Assert().Equal(task.Description, newTask.Description)
	s.Assert().Equal(task.State, newTask.State)
	s.Assert().Equal(task.UserId, newTask.UserId)

	// check number of tasks after add task
	newNumberTasks := len(s.listTasks())
	s.Assert().Equal(oldNumberTasks+1, newNumberTasks)
}

func (s *testSuite) TestUpdateTask() {
	lastTask := s.getLastTask()
	newDescription := "New Description" + uuid.New().String()
	newState := "done"
	// check task data before update
	s.Assert().NotEqual(lastTask.Description, newDescription)
	s.Assert().NotEqual(lastTask.State, newState)

	// add a new task
	task := &common.Task{
		Description: newDescription,
		State:       common.TaskState(newState),
		UserId:      s.getLastUser().Id,
	}
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(task)
	s.Assert().NoError(err)

	req, err := http.NewRequest("PUT", "http://localhost:8080/tasks/"+lastTask.Id, ioutil.NopCloser(&buf))
	s.Assert().NoError(err)

	res, err := http.DefaultClient.Do(req)
	s.Assert().NoError(err)

	err = json.NewDecoder(res.Body).Decode(&lastTask)
	s.Assert().NoError(err)

	s.Assert().Equal(http.StatusOK, res.StatusCode)
	s.Assert().Equal(task.Description, lastTask.Description)
	s.Assert().Equal(task.State, lastTask.State)

	// check task data after update
	newTask := s.getTask(lastTask.Id)
	s.Assert().Equal(lastTask.Description, newTask.Description)
	s.Assert().Equal(lastTask.State, newTask.State)
	s.Assert().Equal(lastTask.UserId, newTask.UserId)

}

func (s *testSuite) TestGetTask() {
	lastTask := s.getLastTask()
	task := &common.Task{}

	req, err := http.NewRequest("GET", "http://localhost:8080/tasks/"+lastTask.Id, nil)
	s.Assert().NoError(err)

	res, err := http.DefaultClient.Do(req)
	s.Assert().NoError(err)

	err = json.NewDecoder(res.Body).Decode(&task)
	s.Assert().NoError(err)

	s.Assert().Equal(http.StatusOK, res.StatusCode)
	s.Assert().Equal(lastTask.Description, task.Description)
	s.Assert().Equal(lastTask.State, task.State)
	s.Assert().Equal(lastTask.UserId, task.UserId)
}

func (s *testSuite) TestDeleteTask() {
	lastTask := s.getLastTask()
	// check number of tasks before addint
	oldNumberTasks := len(s.listTasks())

	req, err := http.NewRequest("DELETE", "http://localhost:8080/tasks/"+lastTask.Id, nil)
	s.Assert().NoError(err)

	res, err := http.DefaultClient.Do(req)
	s.Assert().NoError(err)

	s.Assert().Equal(http.StatusOK, res.StatusCode)

	// check number of tasks after add task
	newNumberTasks := len(s.listTasks())
	s.Assert().Equal(oldNumberTasks-1, newNumberTasks)

}

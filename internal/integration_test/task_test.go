//go:build service
// +build service

package integrationtest

import (
	"net/http"

	"github.com/aborgesrodrigues/to-do-api/internal/common"
	"github.com/google/uuid"
)

func (s *testSuite) listTasks() []common.Task {
	tasks := make([]common.Task, 0)
	s.call("GET", "http://localhost:8080/tasks", nil, &tasks)

	return tasks
}

func (s *testSuite) getTask(id string) *common.Task {
	task := &common.Task{}
	s.call("GET", "http://localhost:8080/tasks/"+id, nil, task)

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

	newTask := &common.Task{}
	res := s.call("POST", "http://localhost:8080/tasks", task, newTask)

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
	res := s.call("PUT", "http://localhost:8080/tasks/"+lastTask.Id, task, lastTask)

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

	res := s.call("GET", "http://localhost:8080/tasks/"+lastTask.Id, nil, task)

	s.Assert().Equal(http.StatusOK, res.StatusCode)
	s.Assert().Equal(lastTask.Description, task.Description)
	s.Assert().Equal(lastTask.State, task.State)
	s.Assert().Equal(lastTask.UserId, task.UserId)
}

func (s *testSuite) TestDeleteTask() {
	lastTask := s.getLastTask()
	// check number of tasks before addint
	oldNumberTasks := len(s.listTasks())

	res := s.call("DELETE", "http://localhost:8080/tasks/"+lastTask.Id, nil, nil)

	s.Assert().Equal(http.StatusOK, res.StatusCode)

	// check number of tasks after add task
	newNumberTasks := len(s.listTasks())
	s.Assert().Equal(oldNumberTasks-1, newNumberTasks)

}

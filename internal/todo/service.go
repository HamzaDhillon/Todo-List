package todo

import "fmt"

type Service struct {
	tasks []Task
}

func NewService() *Service {
	return &Service{tasks: []Task{}}
}

func (s *Service) nextID() int {
	return len(s.tasks) + 1
}

func (s *Service) Add(description string) error {
	if description == "" {
		return fmt.Errorf("description can not be empty")
	}
	s.tasks = append(s.tasks, NewTask(s.nextID(), description))
	return nil
}

func (s *Service) Complete(id int) error {
	for i := range s.tasks {
		if s.tasks[i].ID == id {
			s.tasks[i].IsComplete = true
			return nil
		}
	}
	return fmt.Errorf("task %d not found", id)
}

func (s *Service) List(all bool) ([]Task, error) {
	if all {
		return s.tasks, nil
	}
	out := make([]Task, 0)
	for _, t := range s.tasks {
		if !t.IsComplete {
			out = append(out, t)
		}
	}
	return out, nil
}

func (s *Service) Delete(id int) error {
	for i := range s.tasks {
		if s.tasks[i].ID == id {
			s.tasks = append(s.tasks[:i], s.tasks[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("task %d not found", id)
}
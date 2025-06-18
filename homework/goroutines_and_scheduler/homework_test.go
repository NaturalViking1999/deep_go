package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type Task struct {
	Identifier int
	Priority   int
}

type tasksHeap struct {
	tasks []*Task
}

func (th *tasksHeap) len() int {
	return len(th.tasks)
}

func (th *tasksHeap) swap(a, b int) {
	th.tasks[a], th.tasks[b] = th.tasks[b], th.tasks[a]
}

func (th *tasksHeap) pop() *Task {
	task := th.tasks[0]
	th.tasks = th.tasks[1:]
	return task
}

func (th *tasksHeap) checkPriorityHigher(a, b int) bool {
	return th.tasks[a].Priority > th.tasks[b].Priority
}

func (th *tasksHeap) sequencingTasks(idx int) {
	if idx == 0 || th.checkPriorityHigher((idx-1)/2, idx) {
		return
	}
	th.swap((idx-1)/2, idx)
	th.sequencingTasks((idx - 1) / 2)
}

func (th *tasksHeap) heapify(idx int) {
	if idx < 0 {
		return
	}
	max, left, right := idx, idx*2+1, idx*2+2
	if left < th.len() && th.checkPriorityHigher(left, max) {
		max = left
	}
	if right < th.len() && th.checkPriorityHigher(right, max) {
		max = right
	}
	if max != idx {
		th.swap(max, idx)
	}
	th.heapify(idx - 1)
}

type Scheduler struct {
	tasksHeap
	taskByIdMap map[int]*Task
}

func NewScheduler() Scheduler {
	return Scheduler{
		taskByIdMap: make(map[int]*Task),
	}
}

func (s *Scheduler) AddTask(task Task) {
	s.tasksHeap.tasks = append(s.tasks, &task)
	s.taskByIdMap[task.Identifier] = &task
	s.sequencingTasks(s.len() - 1)
}

func (s *Scheduler) ChangeTaskPriority(taskID int, newPriority int) {
	task, ok := s.taskByIdMap[taskID]
	if !ok {
		return
	}
	task.Priority = newPriority
	s.heapify(s.len() / 2)
}

func (s *Scheduler) GetTask() Task {
	task := s.tasksHeap.pop()
	delete(s.taskByIdMap, task.Identifier)
	s.heapify(0)
	return *task
}

func TestTrace(t *testing.T) {
	task1 := Task{Identifier: 1, Priority: 10}
	task2 := Task{Identifier: 2, Priority: 20}
	task3 := Task{Identifier: 3, Priority: 30}
	task4 := Task{Identifier: 4, Priority: 40}
	task5 := Task{Identifier: 5, Priority: 50}

	scheduler := NewScheduler()
	scheduler.AddTask(task1)
	scheduler.AddTask(task2)
	scheduler.AddTask(task3)
	scheduler.AddTask(task4)
	scheduler.AddTask(task5)

	task := scheduler.GetTask()
	assert.Equal(t, task5, task)

	task = scheduler.GetTask()
	assert.Equal(t, task4, task)

	scheduler.ChangeTaskPriority(1, 100)
	task1.Priority = 100

	task = scheduler.GetTask()
	assert.Equal(t, task1, task)

	task = scheduler.GetTask()
	assert.Equal(t, task3, task)
}

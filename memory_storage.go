package main

type MemoryStorage struct {
	data   map[int]Task
	nextId int
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		data: make(map[int]Task),
	}
}

func (st *MemoryStorage) CreateTask(t Task) int {
	t.Id = st.nextId
	st.nextId++
	st.data[t.Id] = t
	return t.Id
}

func (st *MemoryStorage) GetTask(id int) Task {
	return st.data[id]
}

func (st *MemoryStorage) UpdateTask(id int, t Task) {
	st.data[id] = t
}

func (st *MemoryStorage) DeleteTask(id int) {
	delete(st.data, id)
}

module usecase

import domain
import repository

pub fn list_tasks(r repository.TodoRepository) ![]domain.Task {
	tasks := r.list_tasks()!
	return tasks
}

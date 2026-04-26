module usecase

import domain
import repository

pub fn add_task(r repository.TodoRepository, task domain.Task) ! {
	r.add_task(task)!
}

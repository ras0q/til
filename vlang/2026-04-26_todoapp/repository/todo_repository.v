module repository

import domain
import db.sqlite

@[table: "tasks"]
struct Task {
	id          int @[primary; serial]
	name        string
	description string
	done        bool
}

pub struct TodoRepository {
}

pub fn (r TodoRepository) list_tasks() ![]domain.Task {
	// TODO: connect to DB only once
	db := sqlite.connect('tasks.db')!
	sql db {
		create table Task
	}!

	tasks := sql db {
		select from Task
	}!

	return tasks.map(fn (t Task) domain.Task {
		return domain.Task{
			name:        t.name
			description: t.description
			done:        t.done
		}
	})
}

pub fn (mut r TodoRepository) add_task(task domain.Task) ! {
	db_task := Task{
		name:        task.name
		description: task.description
		done:        task.done
	}

	db := sqlite.connect('tasks.db')!
	sql db {
		create table Task
		insert db_task into Task
	}!
}

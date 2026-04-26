module main

import os
import cli
import usecase
import domain
import repository

fn main() {
	r := repository.TodoRepository{}

	mut app := cli.Command{
		name:        'todo'
		description: 'TODO App written by Vlang'
		execute:     fn (cmd cli.Command) ! {
			println('Usage: todo <list|add|remove|done>')
			return
		}
		commands:    [
			cli.Command{
				name:    'list'
				execute: fn [r] (cmd cli.Command) ! {
					tasks := usecase.list_tasks(r)!
					dump(tasks)
				}
			},
			cli.Command{
				name:    'add'
				execute: fn [r] (cmd cli.Command) ! {
					assert cmd.args.len > 1

					usecase.add_task(r, domain.Task{
						name:        cmd.args[0]
						description: cmd.args[1]
					})!
				}
			},
			cli.Command{
				name:    'remove'
				alias:   'rm'
				execute: fn (cmd cli.Command) ! {
					println('Usage: todo remove <task>')
				}
			},
			cli.Command{
				name:    'done'
				execute: fn (cmd cli.Command) ! {
					println('Usage: todo done <task>')
				}
			},
		]
	}

	app.setup()
	app.parse(os.args)
}

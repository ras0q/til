module domain

pub struct Task {
pub:
	name        string @[required]
	description string
	done        bool
}

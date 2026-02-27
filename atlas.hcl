env "local" {
	url = "mysql://root:yourpassword@localhost:3306/yourdb"
}

variable "db_name" {
	type = string
	default = "yourdb"
}

variable "db_user" {
	type = string
	default = "root"
}

variable "db_pass" {
	type = string
	default = "yourpassword"
}

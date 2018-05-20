package goforjj

type cmd_data struct {
	command     string   // Command to start
	args        []string // Arrays of args to provide to the command
	socket_path string   // Path to store the socket file
	socket_file string
}

type docker_container struct {
	name        string
	opts        []string
	socket_path string
	volumes     map[string]byte
	envs        map[string]byte
}


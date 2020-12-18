package main

type commandID int

const ( //5 commands
	CMD_USERNAME commandID = iota
	CMD_JOIN
	CMD_ROOMS
	CMD_MSG
	CMD_QUIT
)

type command struct {
	id     commandID //unique command ID
	client *client   //sender of the command
	args   []string  //slice of strings from the client's message
}

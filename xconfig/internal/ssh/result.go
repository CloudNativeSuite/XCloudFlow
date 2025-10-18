package ssh

type CommandResult struct {
	Host       string
	ReturnCode int
	ReturnMsg  string
	Output     string
}

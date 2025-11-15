package commands

// Flags holds global flags shared across all commands
type Flags struct {
	LogLevel string
	LogFile  string
}

// ParseFlags holds flags specific to the parse command
type ParseFlags struct {
	BodyKey string
	Pretty  bool
}

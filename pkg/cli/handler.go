package cli

var forgerArgs *ForgerArgs

// SetForgerArgs sets the forger arguments
func SetForgerArgs(args *ForgerArgs) {
	forgerArgs = args
}

// GetForgerArgs returns the forger arguments
func GetForgerArgs() *ForgerArgs {
	return forgerArgs
}

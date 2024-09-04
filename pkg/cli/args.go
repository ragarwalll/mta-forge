package cli

// ForgerArgs is the struct that holds the arguments for the forger command
type ForgerArgs struct {
	BaseDir      string
	OutputDir    string
	Verbose      bool
	Local        bool
	ExpandSource bool
}

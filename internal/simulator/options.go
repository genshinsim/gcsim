package simulator

import "regexp"

//SimulatorOptions sets out the settings to run the sim by (such as debug mode, etc..)
type SimulatorOptions struct {
	PrintResultSummaryToScreen bool   // print summary output to screen?
	ResultSaveToPath           string // file name (excluding ext) to save the result file; if "" then nothing is saved to file
	GZIPResult                 bool   // should the result file be gzipped; only if ResultSaveToPath is not ""
	NumberOfWorkers            int    // how many workers to run the simulation
	//modes
	CalcMode   bool
	ERCalcMode bool

	//other stuff
	DebugRun   bool // run one extra run and generate debug?
	Iterations int  // how many iterations to run
}

var matchOptions = regexp.MustCompile(`^options([^;]*);`)
var matchDebug = regexp.MustCompile(`debug`)
var matchIter = regexp.MustCompile(`iteration=(\d+)`)
var matchDuration = regexp.MustCompile(`duration=(\d+)`)
var matchWorkers = regexp.MustCompile(`workers=(\d+)`)

//ParseOptionsAndRemove will parse sim options from the config string, stripping it from the string input
//and return the option, the config string (with options stripped), and error if any
func ParseOptionsAndRemove(cfg string) (SimulatorOptions, string, error) {
	//options debug=true iteration=1000 duration=87 workers=24;
	//find the options line; return nil if nothing found
	match := matchOptions
}

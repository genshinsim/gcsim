package simulator

//SimulatorOptions sets out the settings to run the sim by (such as debug mode, etc..)
type SimulatorOptions struct {
	PrintResultSummaryToScreen bool   // print summary output to screen?
	ResultSaveToPath           string // file name (excluding ext) to save the result file; if "" then nothing is saved to file
	GZIPResult                 bool   // should the result file be gzipped; only if ResultSaveToPath is not ""

	// NumberOfWorkers            int    // how many workers to run the simulation
	// //modes
	// CalcMode   bool
	// ERCalcMode bool

	// //other stuff
	// DebugRun   bool // run one extra run and generate debug?
	// Iterations int  // how many iterations to run
}

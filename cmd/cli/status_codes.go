package cli

var COMMON_RETURN_CODES = map[int]string{
	99: "CSync not initialized",
}

var ADD_RETURN_CODES = map[int]string{
	1:  "File removed from staging", // file was staged (ADD), but it got removed
	2:  "Staged file updated",       // file was staged (ADD), but it got modified
	3:  "File already staged",       // the user staged the same file again (ADD)
	4:  "",                          // file was staged (MOD), but it got removed
	5:  "Staged file updated",       // file was staged (MOD), but it got modified
	6:  "File already staged",       // the user staged the same file again (MOD)
	7:  "Staged file update",        // file was staged (REM), but it got added back and modified
	8:  "File already staged",       // the user staged the same file again (REM)
	9:  "Filed added to staging",    // file committed but not staged -> staged (REM)
	10: "Filed added to staging",    // file committed but not staged -> staged (MOD)
	11: "File not modified",         // file committed, not staged, not modified
	12: "File added to staging",     // file not committed, not staged -> staged (ADD)
}

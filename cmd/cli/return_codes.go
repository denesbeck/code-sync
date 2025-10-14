package cli

var COMMON_RETURN_CODES = map[int]string{
	001: "CSync not initialized",
}

var ADD_RETURN_CODES = map[int]string{
	101: "File removed from staging", // file was staged (ADD), but it got removed
	102: "Staged file updated",       // file was staged (ADD), but it got modified
	103: "File already staged",       // the user staged the same file again (ADD)
	104: "",                          // file was staged (MOD), but it got removed
	105: "Staged file updated",       // file was staged (MOD), but it got modified
	106: "File already staged",       // the user staged the same file again (MOD)
	107: "Staged file update",        // file was staged (REM), but it got added back and modified
	108: "File already staged",       // the user staged the same file again (REM)
	109: "Filed added to staging",    // file committed but not staged -> staged (REM)
	110: "Filed added to staging",    // file committed but not staged -> staged (MOD)
	111: "File not modified",         // file committed, not staged, not modified
	112: "File added to staging",     // file not committed, not staged -> staged (ADD)
}

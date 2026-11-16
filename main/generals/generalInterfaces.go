package generals

type Paths struct {
	SystemPath       string `json:"config_files_location"`
	BinariesLocation string `json:"executables_location"`
	LogsLocation     string `json:"logs_location"`
}

type ErrorTypesType struct {
	Error       string `json:"error"`
	RequestLogs string `json:"reqLogs"`
	Action      string `json:"ActionsAppendIPBList"`
}

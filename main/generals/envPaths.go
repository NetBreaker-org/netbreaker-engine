package generals

import "path/filepath"

func SetSysPath(hi string) string {
	a, b := filepath.Abs(hi)
	if b != nil {
		return hi
	}
	return a
}

var EnvPath Paths = Paths{SystemPath: SetSysPath("../FHSsimulation")}

// switch to this before production
// var EnvPath Paths = Paths{SystemPath: ""}

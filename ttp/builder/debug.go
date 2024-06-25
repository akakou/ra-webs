package builder

const DEBUG_UNIQUE = "010203"

func debugBuild(string, string) (string, string, error) {
	return "", DEBUG_UNIQUE, nil
}

func EnableDebug() {
	Build = debugBuild
}

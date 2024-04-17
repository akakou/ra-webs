package builder

const DEBUG_UNIQUE = "debug_unique"

func debugBuild(string, string) (string, string, error) {
	return "", DEBUG_UNIQUE, nil
}

func EnableDebug() {
	Build = debugBuild
}

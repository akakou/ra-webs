package builder

var DEBUG_UNIQUE = []byte("010203")

func debugBuild(_, _ string) ([]byte, error) {
	return DEBUG_UNIQUE, nil
}

func EnableDebug() {
	Build = debugBuild
}

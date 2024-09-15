package builder

const DEBUG_UNIQUE = "010203"

func debugBuild(string) (*BuildOutput, error) {
	return &BuildOutput{
		CommitId: "123456",
		UniqueId: []byte(DEBUG_UNIQUE),
	}, nil
}

func EnableDebug() {
	Build = debugBuild
}

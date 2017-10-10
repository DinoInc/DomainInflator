package Schema

var __resolved = make(map[string]*Structure)
var __baseDir string

func SetBaseDir(baseDir string) {

	if baseDir[len(baseDir)-1] != '/' {
		baseDir = baseDir + "/"
	}

	__baseDir = baseDir
}

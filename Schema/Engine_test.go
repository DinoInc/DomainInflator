package Schema

import "testing"

func TestSetBaseDir(t *testing.T) {

	SetBaseDir("./a")
	if __baseDir != "./a/" {
		t.Error("SetBaseDir not add trailing slash")
	}

}

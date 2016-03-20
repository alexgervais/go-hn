package lib

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
)

func Test_LoadConfig_MissingFile(t *testing.T) {

	config := LoadConfig("/this/file/should/not/exist")

	assert.Equal(t, "https://hacker-news.firebaseio.com/v0", config.BaseURL)
}

func Test_LoadConfig_DefaultFile(t *testing.T) {

	config := LoadConfig("")

	assert.NotEmpty(t, config.BaseURL)
}

func Test_LoadConfig_OverrideFile(t *testing.T) {

	file, err := ioutil.TempFile(os.TempDir(), "testconfig")
	if err != nil {
		t.Errorf("Failed to create temporary file", err)
		return
	}
	defer os.Remove(file.Name())
	ioutil.WriteFile(file.Name(), []byte(`base-url: http://another.url`), 0644)

	config := LoadConfig(file.Name())

	assert.Equal(t, "http://another.url", config.BaseURL)
}

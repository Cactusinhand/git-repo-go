package config

import (
	"github.com/alibaba/git-repo-go/file"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
)

var (
	userMisoperationInput = `
# Example git-repo config file, generated by git-repo.
# DO NOT edit this file! Any modification will be overwritten.
#

# Set console verbosity. 1: show info, 2: show debug, 3: show trace
#verbose: 0

# LogLvel for logging to file
#loglevel: warning

# LogRotate defines max size of the logfile
#logrotate: 20M

# LogFile defines where to save log
#logfile:
this is the user misoperation input string
`
)

func Test_InstallRepoConfig(t *testing.T) {
	// 1. Run InstallRepoConfig for init config example file
	InstallRepoConfig()

	configDir, err := GetConfigDir()
	assert.NoError(t, err)

	filename := filepath.Join(configDir, DefaultGitRepoConfigFile+".yml.example")
	assert.NoError(t, err)

	fileOpen, err := file.New(filename).OpenCreateRewrite()
	assert.NoError(t, err)
	defer fileOpen.Close()

	// 2. Simulate user's misoperation to modify config example file
	_, err = fileOpen.WriteString(userMisoperationInput)
	assert.NoError(t, err)

	InstallRepoConfig()

	fileRead, err := file.New(filename).Open()
	assert.NoError(t, err)
	defer fileRead.Close()

	fileInfo, err := fileRead.Stat()
	assert.NoError(t, err)

	fileSize := fileInfo.Size()
	buffer := make([]byte, fileSize)
	_, err = fileRead.Read(buffer)
	assert.NoError(t, err)

	assert.Equal(t, gitRepoConfigExample, string(buffer))
}

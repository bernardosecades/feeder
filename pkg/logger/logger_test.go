package logger_test

import (
	"github.com/bernardosecades/feeder/pkg/logger"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
)

const fileName = "my_log_file.log"

func TestFileLogger(t *testing.T) {
	firstMessage := "what's up, bro"
	secondMessage := "Hi!"

	l := logger.NewFileLogger(fileName)
	l.Log(firstMessage)
	l.Log(secondMessage)

	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(err)
	}

	text := string(content)

	assert.Contains(t, text, firstMessage)
	assert.Contains(t, text, secondMessage)

	tearDown()
}

func tearDown() {
	err := os.Remove(fileName)
	if err != nil {
		panic(err)
	}
}

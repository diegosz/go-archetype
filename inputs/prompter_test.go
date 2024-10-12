package inputs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPrompt(t *testing.T) {
	assert := assert.New(t)

	assert.IsType(&yesNoPrompter{}, NewPrompt(InputSpec{Type: "yesno"}))
	assert.IsType(&noYesPrompter{}, NewPrompt(InputSpec{Type: "noyes"}))
	assert.IsType(&simpleTextPrompter{}, NewPrompt(InputSpec{Type: "text"}))
	assert.IsType(&selectPrompter{}, NewPrompt(InputSpec{Type: "select"}))
}

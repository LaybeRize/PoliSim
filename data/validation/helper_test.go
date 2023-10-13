package validation

import (
	"PoliSim/data/database"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsRoleValid(t *testing.T) {
	assert.False(t, isRoleValid(int(database.PressAccount)-1))
	assert.True(t, isRoleValid(int(database.PressAccount)))
	assert.False(t, isRoleValid(int(database.NotLoggedIn)))
	assert.True(t, isRoleValid(int(database.User)))
	assert.True(t, isRoleValid(int(database.MediaAdmin)))
	assert.True(t, isRoleValid(int(database.Admin)))
	assert.True(t, isRoleValid(int(database.HeadAdmin)))
	assert.False(t, isRoleValid(int(database.HeadAdmin)+1))
}

func TestIsEmptyOrNotInRange(t *testing.T) {
	assert.False(t, isValidString("", -1))
	assert.False(t, isValidString("", 12))
	assert.True(t, isValidString("aaaa", -1))
	assert.False(t, isValidString("aaaa", 2))
	assert.True(t, isValidString("aaaa", 10))
}

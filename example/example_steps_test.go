package example

import (
	"encoding/json"
	"strconv"
	"testing"

	"github.com/sittercity/gocumber"
	"github.com/stretchr/testify/assert"
)

type userEntity struct {
	username  string `json:"username"`
	first     string `json:"first"`
	last      string `json:"last"`
	activated bool   `json:"activated"`
}

// This is a contrived example using local variables. You can do whatever you want within each step. This is
// just designed as a showcase for the various ways you can pull data and execute steps with gocumber.
func TestExample(t *testing.T) {
	steps := make(gocumber.Definitions)

	var inputData map[string]string
	users := make(map[string]userEntity)

	var responseUser userEntity

	steps.Given(`I have an existing user with the name "(.*)"`, func(matches []string, _ gocumber.StepNode) {
		user := userEntity{
			username:  matches[1],
			activated: true,
		}

		users[user.username] = user
	})

	steps.Given("I have an existing user", func([]string, gocumber.StepNode) {
		steps.Exec(`I have an existing user with name "Cool Dude"`)
	})

	steps.When("I create a user with the following json data:", func(matches []string, step gocumber.StepNode) {
		payload := step.PyString().String()

		userData := make(map[string]interface{})
		err := json.Unmarshal([]byte(payload), &userData)
		assert.NoError(t, err)

		user := userEntity{
			username:  userData["username"].(string),
			first:     userData["first"].(string),
			last:      userData["last"].(string),
			activated: userData["activated"].(bool),
		}

		users[user.username] = user
	})

	steps.When("I create a user with the following table data:", func(matches []string, step gocumber.StepNode) {
		// Please note that the header (key/value) column is ALSO parsed, i.e. inputData["key"] returns "value".
		// But this doesn't affect the rest of the data so it can be ignored. You can leave off the
		// header column if you don't want it for readability.
		inputData = make(map[string]string)
		for _, row := range step.Table().Rows() {
			inputData[row[0]] = row[1]
		}

		user := userEntity{}
		user.username = inputData["username"]
		user.first = inputData["first"]
		user.last = inputData["last"]
		user.activated, _ = strconv.ParseBool(inputData["activated"])

		users[user.username] = user
	})

	steps.When(`I deactivate the user named "(.*)"`, func(matches []string, step gocumber.StepNode) {
		username := matches[1]
		user := users[username]

		user.activated = false

		users[username] = user
	})

	steps.When(`I read the user named "(.*)"`, func(matches []string, step gocumber.StepNode) {
		username := matches[1]
		responseUser = users[username] // Saving this for later, just for illustration purposes
	})

	steps.Then("the user should be created with the expected data", func([]string, gocumber.StepNode) {
		user := users["CoolDude"]

		assert.Equal(t, "CoolDude", user.username)
		assert.Equal(t, "Cool", user.first)
		assert.Equal(t, "Dude", user.last)
		assert.True(t, user.activated)
	})

	steps.Then(`the user "(.*)" should be deactivated`, func(matches []string, step gocumber.StepNode) {
		username := matches[1]
		user := users[username]

		assert.False(t, user.activated)
	})

	steps.Then("I should see the following data:", func(matches []string, step gocumber.StepNode) {
		// Note that the first 'header' row (username, first, last) is also read into the inputData
		// but can be ignored. In this case just start from row 1. You can leave off the header row
		// if this bothers you, there is nothing special about it.
		var expectedUserData []map[string]string
		expectedUserData = make([]map[string]string, 0)

		for _, row := range step.Table().Rows() {
			newExpectedUser := make(map[string]string)
			newExpectedUser["username"] = row[0]
			newExpectedUser["first"] = row[1]
			newExpectedUser["last"] = row[2]

			expectedUserData = append(expectedUserData, newExpectedUser)
		}

		// This is a bit convoluted but we wanted to show how it was just a basic table of data. You can
		// do whatever you want with it.
		assert.Equal(t, expectedUserData[1]["username"], responseUser.username)
		assert.Equal(t, expectedUserData[1]["first"], responseUser.first)
		assert.Equal(t, expectedUserData[1]["last"], responseUser.last)
	})

	steps.Run(t, "example.feature")
}

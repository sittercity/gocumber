package gocumber

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type FuncTestingFramework struct {
	err func(args ...interface{})
	log func(args ...interface{})
}

func (t FuncTestingFramework) Error(args ...interface{}) { t.err(args...) }
func (t FuncTestingFramework) Log(args ...interface{})   { t.log(args...) }

func TestRun_HappyPath(t *testing.T) {
	steps := make(Definitions)

	// Simply defining steps to prove that it parses right, no need to fill them in
	steps.When("I create a user with the following json data:", func([]string, StepNode) {})
	steps.Then("the user should be created with the expected data", func([]string, StepNode) {})

	steps.Run(t, "test/valid.feature")
}

func TestRun_FailsOnMissingFile(t *testing.T) {
	steps := make(Definitions)
	tt := new(testing.T)

	steps.Run(tt, "file_does_not_exist")

	assert.True(t, tt.Failed())
}

func TestRun_FailsOnInvalidGherkin(t *testing.T) {
	steps := make(Definitions)
	tt := new(testing.T)

	steps.Run(tt, "test/invalid.feature")

	assert.True(t, tt.Failed())
}

func TestRun_FailsOnUndefinedSteps(t *testing.T) {
	steps := make(Definitions)
	tt := new(testing.T)

	steps.Run(tt, "test/valid.feature")

	assert.True(t, tt.Failed())
}

func ExampleRun_WithUndefinedSteps() {
	steps := make(Definitions)

	t := FuncTestingFramework{
		err: func(args ...interface{}) { fmt.Println(args...) },
	}
	steps.Run(t, "test/valid_with_url_params.feature")

	// Output:
	// Undefined step:
	// When I get "/something/%{UUID}"
}

func ExampleRun_WithFailingSteps() {
	steps := make(Definitions)

	t := FuncTestingFramework{
		err: func(args ...interface{}) { fmt.Println(args...) },
		log: func(args ...interface{}) { fmt.Println(args...) },
	}
	steps.When("I create a user with the following json data:", func([]string, StepNode) {})
	steps.Then("the user should be created with the expected data", func([]string, StepNode) {
		t.Error("Expectation failed")
	})

	steps.Run(t, "test/valid.feature")

	// Output:
	// Scenario: Create a user with a json payload
	// Expectation failed
}

func TestRun_SuccessWithOutlineSteps(t *testing.T) {
	steps := make(Definitions)

	steps.Given("I have no users", func([]string, StepNode) {})
	steps.When("I create a new user with the following data:", func([]string, StepNode) {})
	steps.Then("no users should be created", func([]string, StepNode) {})

	steps.Run(t, "test/valid_with_outline.feature")
}

func TestRun_SuccessWithPyString(t *testing.T) {
	steps := make(Definitions)

	steps.Given("I do something the following json data:", func([]string, StepNode) {})
	steps.When("I do something", func([]string, StepNode) {})
	steps.Then("something should have happened", func([]string, StepNode) {})

	steps.Run(t, "test/valid_with_pystring.feature")
}

func TestColumnMap_Happy(t *testing.T) {
	steps := make(Definitions)

	var called bool
	steps.When("I create something with the following table data:", func(_ []string, step StepNode) {
		called = true
		assert.Equal(t,
			map[string]string{
				"key":  "value",
				"key1": "value1",
				"key2": "value2",
				"key3": "value3",
			},
			ColumnMap(step.Table()))
	})

	steps.Run(t, "test/valid_with_table_data.feature")

	assert.True(t, called)
}

func TestExec_MatchFound(t *testing.T) {
	steps := make(Definitions)

	steps.When("I create a user with the following json data:", func([]string, StepNode) {})
	steps.Then("the user should be created with the expected data", func([]string, StepNode) {})

	assert.True(t, steps.Exec("the user should be created with the expected data"))
}

func TestExec_NoMatch(t *testing.T) {
	steps := make(Definitions)

	steps.When("I create a user with the following json data:", func([]string, StepNode) {})
	steps.Then("the user should be created with the expected data", func([]string, StepNode) {})

	assert.False(t, steps.Exec("some unknown step"))
}

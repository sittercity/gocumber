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
	tt := new(testing.T)

	// Simply defining steps to prove that it parses right, no need to fill them in
	steps.When("I create a user with the following json data:", func([]string, StepNode) {})
	steps.Then("the user should be created with the expected data", func([]string, StepNode) {})

	steps.Run(tt, "test/valid.feature")

	assert.False(t, tt.Failed())
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

	tt := FuncTestingFramework{
		err: func(args ...interface{}) { fmt.Println(args...) },
	}
	steps.Run(tt, "test/valid_with_url_params.feature")

	// Output:
	// Undefined step:
	// When I get "/something/%{UUID}"
}

func TestRun_FailsOnPendingSteps(t *testing.T) {
	steps := make(Definitions)
	tt := new(testing.T)

	steps.When("I create a user with the following json data:", nil)
	steps.Then("the user should be created with the expected data", nil)

	steps.Run(tt, "test/valid.feature")

	assert.True(t, tt.Failed())
}

func ExampleRun_WithPendingSteps() {
	steps := make(Definitions)

	tt := FuncTestingFramework{
		err: func(args ...interface{}) { fmt.Println(args...) },
		log: func(args ...interface{}) { fmt.Println(args...) },
	}
	steps.When("I create a user with the following json data:", nil)
	steps.Then("the user should be created with the expected data", nil)

	steps.Run(tt, "test/valid.feature")

	// Output:
	// Scenario: Create a user with a json payload
	// Pending step:
	// When I create a user with the following json data:
}

func ExampleRun_WithFailingSteps() {
	steps := make(Definitions)

	tt := FuncTestingFramework{
		err: func(args ...interface{}) { fmt.Println(args...) },
		log: func(args ...interface{}) { fmt.Println(args...) },
	}
	steps.When("I create a user with the following json data:", func([]string, StepNode) {})
	steps.Then("the user should be created with the expected data", func([]string, StepNode) {
		tt.Error("Expectation failed")
	})

	steps.Run(tt, "test/valid.feature")

	// Output:
	// Scenario: Create a user with a json payload
	// Expectation failed
}

func TestRun_SuccessWithBackground(t *testing.T) {
	steps := make(Definitions)
	tt := new(testing.T)

	var called int
	steps.Given("I am able to perform", func([]string, StepNode) {
		called++
	})
	steps.When("I perform admirably", func([]string, StepNode) {})
	steps.When("I act a fool", func([]string, StepNode) {})
	steps.Then("things should go well", func([]string, StepNode) {})

	steps.Run(tt, "test/valid_with_background.feature")

	assert.Equal(t, 2, called, "Expected background to be executed for each scenario")
	assert.False(t, tt.Failed())
}

func TestRun_SuccessWithOutlineSteps(t *testing.T) {
	steps := make(Definitions)
	tt := new(testing.T)

	steps.Given("I have no users", func([]string, StepNode) {})
	steps.Then("no users should be created", func([]string, StepNode) {})

	var called int
	steps.When("I create a new user with the following data:", func(_ []string, step StepNode) {
		called++
		switch called {
		case 1:
			assert.Equal(t,
				map[string]string{
					"key":      "value",
					"username": "",
					"first":    "",
					"last":     "",
				},
				ColumnMap(step.Table()))
		case 2:
			assert.Equal(t,
				map[string]string{
					"key":      "value",
					"username": "",
					"first":    "fname",
					"last":     "lname",
				},
				ColumnMap(step.Table()))
		case 3:
			assert.Equal(t,
				map[string]string{
					"key":      "value",
					"username": "newuser",
					"first":    "",
					"last":     "lname",
				},
				ColumnMap(step.Table()))
		case 4:
			assert.Equal(t,
				map[string]string{
					"key":      "value",
					"username": "newuser",
					"first":    "fname",
					"last":     "",
				},
				ColumnMap(step.Table()))
		case 5:
			assert.Equal(t,
				map[string]string{
					"key":      "value",
					"username": "newuser",
					"first":    "fname",
					"last":     "lname",
				},
				ColumnMap(step.Table()))
		}
	})

	steps.Run(tt, "test/valid_with_outline.feature")

	assert.Equal(t, 5, called, "Expected scenario to be executed for each outline example")
	assert.False(t, tt.Failed())
}

func TestRun_SuccessWithPyString(t *testing.T) {
	steps := make(Definitions)
	tt := new(testing.T)

	steps.Given("I do something the following json data:", func([]string, StepNode) {})
	steps.When("I do something", func([]string, StepNode) {})
	steps.Then("something should have happened", func([]string, StepNode) {})

	steps.Run(tt, "test/valid_with_pystring.feature")

	assert.False(t, tt.Failed())
}

func TestColumnMap_Happy(t *testing.T) {
	steps := make(Definitions)
	tt := new(testing.T)

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

	steps.Run(tt, "test/valid_with_column_table_data.feature")

	assert.True(t, called)
	assert.False(t, tt.Failed())
}

func TestRowMaps_Happy(t *testing.T) {
	steps := make(Definitions)
	tt := new(testing.T)

	var called bool
	steps.When("I create something with the following table data:", func(_ []string, step StepNode) {
		called = true
		assert.Equal(t,
			[]map[string]string{
				map[string]string{
					"key1": "value1",
					"key2": "value2",
					"key3": "value3",
					"key4": "value4",
				},
				map[string]string{
					"key1": "more1",
					"key2": "more2",
					"key3": "more3",
					"key4": "more4",
				},
			},
			RowMaps(step.Table()))
	})

	steps.Run(tt, "test/valid_with_row_table_data.feature")

	assert.True(t, called)
	assert.False(t, tt.Failed())
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

func TestDocstrings_OutlineVariables(t *testing.T) {
	steps := make(Definitions)
	tt := new(testing.T)

	var called int
	steps.Given("something is:", func(_ []string, step StepNode) {
		called++
		assert.Equal(t, []string{"minimally functional"}, step.PyString().Lines())
	})

	steps.Run(t, "test/valid_with_pystring_outline.feature")

	assert.Equal(t, 1, called, "Expected scenario to be executed for each outline example")
	assert.False(t, tt.Failed())
}

func TestDocstrings_OutlineVariablesWithMultipleExamples(t *testing.T) {
	steps := make(Definitions)
	tt := new(testing.T)

	var expected_replaced_pystrings = []string{"minimally functional barely", "incredibly functional overwhelmingly"}

	var called int
	steps.Given("something is:", func(_ []string, step StepNode) {
		assert.Equal(t, []string{expected_replaced_pystrings[called]}, step.PyString().Lines())
		called++
	})

	steps.Run(t, "test/valid_with_pystring_outline_multiple_examples.feature")

	assert.Equal(t, 2, called, "Expected scenario to be executed for each outline example")
	assert.False(t, tt.Failed())
}

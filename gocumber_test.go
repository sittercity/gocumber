package gocumber

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
	steps := make(Definitions)

	// Simply defining steps to prove that it parses right, no need to fill them in
	steps.When("I create a user with the following json data:", func([]string, StepNode) {})
	steps.Then("the user should be created with the expected data", func([]string, StepNode) {})

	steps.Run(t, "test/valid.feature")
}

func TestParseFile_FailsOnEmptyFile(t *testing.T) {
	steps := make(Definitions)
	_, errs := steps.parseFile("file_doesn't_exist")

	assert.NotEmpty(t, errs)
	assert.Error(t, errs[0])
}

func TestParseFile_FailsOnInvalidGherkin(t *testing.T) {
	steps := make(Definitions)
	_, errs := steps.parseFile("test/invalid.feature")

	assert.NotEmpty(t, errs)
	assert.Error(t, errs[0])
}

func TestParseFile_FailsOnUnDefinedSteps(t *testing.T) {
	steps := make(Definitions)
	_, errs := steps.parseFile("test/valid.feature")

	assert.NotEmpty(t, errs)
	assert.Error(t, errs[0])
}

func TestParseFile_SuccessWithOutlineSteps(t *testing.T) {
	steps := make(Definitions)

	steps.Given("I have no users", func([]string, StepNode) {})
	steps.When("I create a new user with the following data:", func([]string, StepNode) {})
	steps.Then("no users should be created", func([]string, StepNode) {})

	_, errs := steps.parseFile("test/valid_with_outline.feature")

	assert.Empty(t, errs)
}

func TestParseFile_SuccessWithPyString(t *testing.T) {
	steps := make(Definitions)

	steps.Given("I do something the following json data:", func([]string, StepNode) {})
	steps.When("I do something", func([]string, StepNode) {})
	steps.Then("something should have happened", func([]string, StepNode) {})

	_, errs := steps.parseFile("test/valid_with_pystring.feature")

	assert.Empty(t, errs)
}

func TestColumnMap_Happy(t *testing.T) {
	steps := make(Definitions)

	steps.When("I create something with the following table data:", func([]string, StepNode) {})

	definitions, errs := steps.parseFile("test/valid_with_table_data.feature")
	assert.Empty(t, errs)

	result := ColumnMap(definitions[0].step.Table())
	assert.Equal(t, "value", result["key"])
	assert.Equal(t, "value1", result["key1"])
	assert.Equal(t, "value2", result["key2"])
	assert.Equal(t, "value3", result["key3"])
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

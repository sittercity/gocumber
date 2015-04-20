package gocumber

import (
	"fmt"
	"io/ioutil"
	"regexp"

	"github.com/muhqu/go-gherkin"
	"github.com/muhqu/go-gherkin/nodes"
)

type Definition func([]string, StepNode)
type Definitions map[*regexp.Regexp]Definition
type StepNode nodes.StepNode
type Table nodes.TableNode

type matchedDefinition struct {
	step       nodes.StepNode
	matches    []string
	definition Definition
}

func (m matchedDefinition) execute()      { m.definition(m.matches, m.step) }
func (m matchedDefinition) pending() bool { return m.definition == nil }

type testingFramework interface {
	Error(args ...interface{})
	Log(args ...interface{})
}

func ColumnMap(table Table) map[string]string {
	result := make(map[string]string)

	for _, row := range table.Rows() {
		result[row[0]] = row[1]
	}

	return result
}

func RowMap(table Table) map[string]string {
	result := make(map[string]string)

	header := table.Rows()[0]
	for _, row := range table.Rows()[1:] {
		for i, key := range header {
			result[key] = row[i]
		}
	}

	return result
}

func (defs Definitions) Step(text string, def Definition) {
	defs[regexp.MustCompile("^"+text+"$")] = def
}

func (defs Definitions) Exec(text string) (found bool) {
	matched := defs.find(nodes.NewMutableStepNode("", text))
	if matched == nil {
		return false
	}

	matched.execute()
	return true
}

func (defs Definitions) find(step nodes.StepNode) *matchedDefinition {
	var matches []string
	text := step.Text()

	for re, definition := range defs {
		matches = re.FindStringSubmatch(text)
		if matches != nil {
			return &matchedDefinition{
				step:       step,
				matches:    matches,
				definition: definition,
			}
		}
	}

	return nil
}

func (defs Definitions) findAll(steps []nodes.StepNode) (matched []matchedDefinition, missing map[string]nodes.StepNode) {
	missing = make(map[string]nodes.StepNode)

	for _, step := range steps {
		found := defs.find(step)

		if found == nil {
			missing[step.Text()] = step
		} else {
			matched = append(matched, *found)
		}
	}

	return matched, missing
}

func (defs *Definitions) Given(text string, def Definition) { defs.Step(text, def) }
func (defs *Definitions) When(text string, def Definition)  { defs.Step(text, def) }
func (defs *Definitions) Then(text string, def Definition)  { defs.Step(text, def) }

func (defs *Definitions) Run(t testingFramework, file string) {
	if buffer, err := ioutil.ReadFile(file); err != nil {
		t.Error(err)
	} else if feature, err := gherkin.ParseGherkinFeature(string(buffer)); err != nil {
		t.Error(err)
	} else {

	ScenarioLoop:
		for _, scenario := range feature.Scenarios() {
			var steps []nodes.StepNode

			if background := feature.Background(); background != nil {
				steps = append(steps, background.Steps()...)
			}

			switch scenario := scenario.(type) {
			case nodes.OutlineNode:
				outlineSteps(scenario, func(step nodes.StepNode) {
					steps = append(steps, step)
				})
			default:
				steps = append(steps, scenario.Steps()...)
			}

			matched, missing := defs.findAll(steps)

			if len(missing) > 0 {
				for _, step := range missing {
					t.Error(fmt.Errorf("Undefined step:\n%s %s", step.StepType(), step.Text()))
				}
				continue ScenarioLoop
			}

			t.Log("Scenario: " + scenario.Title())

			for _, definition := range matched {
				if definition.pending() {
					t.Error(fmt.Errorf("Pending step:\n%s %s", definition.step.StepType(), definition.step.Text()))
					continue ScenarioLoop
				}

				definition.execute()
			}
		}
	}
}

func outlineSteps(outline nodes.OutlineNode, callback func(nodes.StepNode)) {
	matcher := regexp.MustCompile(`<[^>]+>`)
	replacements := make(map[string]string)
	replace := func(text string) string {
		return matcher.ReplaceAllStringFunc(text, func(key string) string {
			return replacements[key]
		})
	}

	header := outline.Examples().Table().Rows()[0]
	for _, row := range outline.Examples().Table().Rows()[1:] {
		for i, key := range header {
			replacements["<"+key+">"] = row[i]
		}

		for _, original := range outline.Steps() {
			step := nodes.NewMutableStepNode(original.StepType(), replace(original.Text()))

			if original.Table() != nil {
				table := nodes.NewMutableTableNode()

				for _, row := range original.Table().Rows() {
					replaced := make([]string, len(row))

					for i, cell := range row {
						replaced[i] = replace(cell)
					}

					table.AddRow(replaced)
				}

				step.SetTable(table)
			} else if original.PyString() != nil {
				pyString := nodes.NewMutablePyStringNode()
				for _, line := range original.PyString().Lines() {
					pyString.AddLine(replace(line))
				}

				step.SetPyString(pyString)
			}

			callback(step)
		}
	}
}

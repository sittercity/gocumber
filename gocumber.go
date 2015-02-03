package gocumber

import (
	"errors"
	"io/ioutil"
	"regexp"
	"testing"

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

func (matched matchedDefinition) execute() {
	matched.definition(matched.matches, matched.step)
}

func ColumnMap(table Table) map[string]string {
	result := make(map[string]string)

	for _, row := range table.Rows() {
		result[row[0]] = row[1]
	}

	return result
}

// Helper method to address weird edge cases. See test case for specific
// behavior.
//TODO Need to revisit this function, the name is misleading
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

func (defs *Definitions) Run(t *testing.T, file string) {
	definitions, errs := defs.parseFile(file)

	if errs != nil && len(errs) != 0 {
		for _, err := range errs {
			t.Errorf(err.Error())
		}
	} else {
		for _, definition := range definitions {
			definition.execute()
		}
	}
}

func (defs *Definitions) parseFile(file string) (definitions []matchedDefinition, errs []error) {
	definitions = make([]matchedDefinition, 0, 0)
	errs = make([]error, 0, 0)

	if buffer, err := ioutil.ReadFile(file); err != nil {
		errs = append(errs, err)
	} else if feature, err := gherkin.ParseGherkinFeature(string(buffer)); err != nil {
		errs = append(errs, err)
	} else {
		for _, scenario := range feature.Scenarios() {
			var steps []nodes.StepNode

			// TODO append Background

			switch scenario := scenario.(type) {
			case nodes.OutlineNode:
				outlineSteps(scenario, func(step nodes.StepNode) {
					steps = append(steps, step)
				})
			default:
				for _, step := range scenario.Steps() {
					steps = append(steps, step)
				}
			}

			matched, missing := defs.findAll(steps)

			if len(missing) == 0 {
				for _, definition := range matched {
					definitions = append(definitions, definition)
				}
			} else {
				for _, step := range missing {
					errs = append(errs, errors.New("Undefined step:\n"+step.StepType()+" "+step.Text()))
				}
			}
		}
	}

	return definitions, errs
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
			// FIXME
			step.SetPyString(nodes.NewMutablePyStringNode())
		}

		callback(step)
	}
}

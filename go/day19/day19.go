package main

import (
	"common"
	"fmt"
	"maps"
	"strings"
)

func main() {
	Solution1()
	Solution2()
}

func Solution2() {
	workflows, _ := parseInterval()
	interval := IntervalPart{"x": {1, 4000}, "m": {1, 4000}, "a": {1, 4000}, "s": {1, 4000}}
	partIntervals := filterAcceptedIntervals(interval, &workflows)
	combinations := sumCombinations(partIntervals)

	fmt.Println(combinations)
}

func Solution1() {
	workflows, parts := parse()
	accepted := filterAccepted(parts, &workflows)
	totalRating := sumPartRatings(accepted)

	fmt.Println(totalRating)
}

func sumCombinations(partIntervals []IntervalPart) uint64 {
	sum := uint64(0)
	for _, match := range partIntervals {
		sum += match.calculateDistinctCombinations()
	}
	return sum
}

func filterAcceptedIntervals(part IntervalPart, workflows *IntervalWorkflows) []IntervalPart {
	matches := []IntervalPart{}
	for current, backlog := pop([]IntervalParts{{"in", part}}); len(backlog) >= 0; current, backlog = backlog[0], backlog[1:] {
		next, matched := nextMatches(current.match, current.part, workflows)
		matches = append(matches, matched...)
		backlog = append(backlog, next...)

		if len(backlog) <= 0 {
			break
		}
	}

	return matches
}

func pop[S ~[]E, E any](s S) (E, S) {
	lastIndex := len(s) - 1
	last := s[lastIndex]
	return last, s[:lastIndex]
}

func nextMatches(workflow string, part IntervalPart, workflows *IntervalWorkflows) ([]IntervalParts, []IntervalPart) {
	currentWorkflow := (*workflows)[workflow]
	prevMatches := []func(part IntervalPart) IntervalPart{}
	intervals := []IntervalParts{}
	accepted := []IntervalPart{}

	for _, rule := range currentWorkflow.rules {
		latestPart := maps.Clone(part)

		for _, mm := range prevMatches {
			latestPart = mm(latestPart)
		}

		interval, next := rule(latestPart)
		prevMatches = append(prevMatches, next)

		if !latestPart.isValid() || !interval.part.isValid() {
			continue
		}

		if interval.match == ACCEPT {
			accepted = append(accepted, interval.part)
		} else if interval.match != REJECT {
			intervals = append(intervals, interval)
		}
	}

	return intervals, accepted
}

func sumPartRatings(parts []Part) int {
	sum := 0

	for _, part := range parts {
		sum += calculatePartRating(part)
	}

	return sum
}

func calculatePartRating(part Part) int {
	sum := 0

	for _, rating := range part {
		sum += rating
	}

	return sum
}

func filterAccepted(parts []Part, workflows *Workflows) []Part {
	evaluate := createWorkflowValidator(workflows)
	accepted := []Part{}

	for _, part := range parts {
		if evaluate(part) {
			accepted = append(accepted, part)
		}
	}

	return accepted
}

func createWorkflowValidator(workflows *Workflows) func(part Part) bool {
	return func(part Part) bool {
		nextWorkflow := createWorkflowIterator("in", workflows)

		for value := nextWorkflow(part); true; value = nextWorkflow(part) {
			if value == ACCEPT {
				return true
			}

			if value == REJECT {
				return false
			}
		}

		return false
	}
}

func createWorkflowIterator(start string, workflows *Workflows) func(part Part) string {
	currentWorkflow := (*workflows)[start]

	return func(part Part) string {
		for _, rule := range currentWorkflow.rules {
			if nextWorkflow, fulfilled := rule(part); fulfilled {
				currentWorkflow = (*workflows)[nextWorkflow]
				return nextWorkflow
			}
		}
		return ""
	}
}

func parse() (Workflows, []Part) {
	workflows := Workflows{}
	parts := []Part{}

	parsePart := createPartParser()

	common.ReadFile("day19/day19.txt", func(line string) {
		if line == "" {
			return
		} else if strings.HasPrefix(line, "{") {
			parts = append(parts, parsePart(line))
		} else {
			workflow := parseWorkflow(line)
			workflows[workflow.value] = workflow
		}
	})

	return workflows, parts
}

func parseInterval() (IntervalWorkflows, []Part) {
	workflows := IntervalWorkflows{}
	parts := []Part{}

	parsePart := createPartParser()

	common.ReadFile("day19/day19.txt", func(line string) {
		if line == "" {
			return
		} else if strings.HasPrefix(line, "{") {
			parts = append(parts, parsePart(line))
		} else {
			workflow := parseIntervalWorkflow(line)
			workflows[workflow.value] = workflow
		}
	})

	return workflows, parts
}

func createPartParser() func(line string) Part {
	replacer := strings.NewReplacer("{", "", "}", "")

	return func(line string) Part {
		values := replacer.Replace(line)
		ratings := Part{}

		for _, value := range strings.Split(values, ",") {
			category, rating := common.Split2(value, "=")
			ratings[category] = common.ParseInt(rating)
		}

		return ratings
	}
}

func parseWorkflow(line string) Workflow {
	value, rules := common.Split2(line, "{")
	rules = strings.TrimSuffix(rules, "}")

	workflowRules := []Rule{}

	for _, rule := range strings.Split(rules, ",") {
		if !strings.Contains(rule, ":") {
			workflowRules = append(workflowRules, func(part Part) (string, bool) {
				return rule, true
			})
		} else {
			unparsedCondition, mappedValue := common.Split2(rule, ":")

			condition := parseCondition(unparsedCondition)
			workflowRules = append(workflowRules, func(part Part) (string, bool) {
				return mappedValue, condition(part)
			})
		}
	}

	return Workflow{value, workflowRules}
}

func parseIntervalWorkflow(line string) IntervalWorkflow {
	value, rules := common.Split2(line, "{")
	rules = strings.TrimSuffix(rules, "}")

	workflowRules := []IntervalRule{}

	for _, rule := range strings.Split(rules, ",") {
		if !strings.Contains(rule, ":") {
			workflowRules = append(workflowRules, func(part IntervalPart) (IntervalParts, func(part IntervalPart) IntervalPart) {
				return IntervalParts{rule, part}, nil
			})
		} else {
			unparsedCondition, mappedValue := common.Split2(rule, ":")
			nextValue := parseIntervalCondition(unparsedCondition)
			workflowRules = append(workflowRules, func(part IntervalPart) (IntervalParts, func(part IntervalPart) IntervalPart) {
				mappedPart, nextPart := nextValue(part)
				return IntervalParts{mappedValue, mappedPart}, nextPart
			})
		}
	}

	return IntervalWorkflow{value, workflowRules}
}

func parseCondition(condition string) func(part Part) bool {
	if strings.Contains(condition, "<") {
		category, conditionValueString := common.Split2(condition, "<")
		conditionValue := common.ParseInt(conditionValueString)

		return func(part Part) bool {
			return part[category] < conditionValue
		}
	} else if strings.Contains(condition, ">") {
		category, conditionValueString := common.Split2(condition, ">")
		conditionValue := common.ParseInt(conditionValueString)

		return func(part Part) bool {
			return part[category] > conditionValue
		}
	}

	return func(part Part) bool {
		return true
	}
}

func parseIntervalCondition(condition string) func(part IntervalPart) (IntervalPart, func(part IntervalPart) IntervalPart) {
	if strings.Contains(condition, "<") {
		category, conditionValueString := common.Split2(condition, "<")
		conditionValue := common.ParseInt(conditionValueString)

		return func(part IntervalPart) (IntervalPart, func(part IntervalPart) IntervalPart) {
			matchedInterval := part
			currentInterval := part[category]
			matchedInterval[category] = Interval{currentInterval.from, conditionValue - 1}

			return matchedInterval, func(part IntervalPart) IntervalPart {
				nextInterval := part
				nextInterval[category] = Interval{conditionValue, currentInterval.to}
				return nextInterval
			}
		}
	} else if strings.Contains(condition, ">") {
		category, conditionValueString := common.Split2(condition, ">")
		conditionValue := common.ParseInt(conditionValueString)

		return func(part IntervalPart) (IntervalPart, func(part IntervalPart) IntervalPart) {
			matchedInterval := part
			currentInterval := part[category]
			matchedInterval[category] = Interval{conditionValue + 1, currentInterval.to}

			return matchedInterval, func(part IntervalPart) IntervalPart {
				nextInterval := part
				nextInterval[category] = Interval{currentInterval.from, conditionValue}
				return nextInterval
			}
		}
	}

	return func(part IntervalPart) (IntervalPart, func(part IntervalPart) IntervalPart) {
		return part, nil
	}
}

type Workflow struct {
	value string
	rules []Rule
}

type Workflows map[string]Workflow
type Rule func(part Part) (string, bool)
type Part map[string]int

type IntervalWorkflow struct {
	value string
	rules []IntervalRule
}

type Interval struct {
	from, to int
}

type IntervalWorkflows map[string]IntervalWorkflow
type IntervalRule func(part IntervalPart) (IntervalParts, func(part IntervalPart) IntervalPart)
type IntervalPart map[string]Interval

func (this IntervalPart) calculateDistinctCombinations() uint64 {
	multiple := uint64(1)
	for _, part := range this {
		multiple *= uint64((part.to - part.from) + 1)
	}
	return multiple
}

func (this IntervalPart) isValid() bool {
	for _, part := range this {
		if part.to < part.from {
			return false
		}
	}
	return true
}

type IntervalParts struct {
	match string
	part  IntervalPart
}

type Condition struct {
	lessThan, greaterThan int
}

// 18446523642893126191
// 167409079868000
// 196513904884482
const (
	ACCEPT = "A"
	REJECT = "R"
)

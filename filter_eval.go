package main

import (
	"regexp"
	"strconv"
	"time"

	"github.com/sachaos/todoist/lib"
)

var priorityRegex = regexp.MustCompile("^p([1-4])$")

// Eval ...
func Eval(e Expression, item todoist.Item, c *todoist.Client) (result bool, err error) {
	result = false
	switch e.(type) {
	case BoolInfixOpExpr:
		e := e.(BoolInfixOpExpr)
		lr, err := Eval(e.left, item, c)
		rr, err := Eval(e.right, item, c)
		if err != nil {
			return false, nil
		}
		switch e.operator {
		case '&':
			return lr && rr, nil
		case '|':
			return lr || rr, nil
		}
	case StringExpr:
		e := e.(StringExpr)
		return EvalAsPriority(e, item), err
	case DueDateExpr:
		e := e.(DueDateExpr)
		return EvalDueDate(e, item), err
	case LabelExpr:
		e := e.(LabelExpr)
		return EvalLabel(e, item, c), err
	case ProjectExpr:
		e := e.(ProjectExpr)
		return EvalProject(e, item, c), err
	case NotOpExpr:
		e := e.(NotOpExpr)
		r, err := Eval(e.expr, item, c)
		if err != nil {
			return false, nil
		}
		return !r, nil
	default:
		return true, err
	}
	return
}

func EvalDueDate(e DueDateExpr, item todoist.Item) (result bool) {
	itemDueDate := item.DueDateTime()
	if (itemDueDate == time.Time{}) {
		if e.operation == NO_DUE_DATE {
			return true
		}
		return false
	}
	allDay := e.allDay
	dueDate := e.datetime
	switch e.operation {
	case DUE_ON:
		var startDate, endDate time.Time
		if allDay {
			startDate = dueDate
			endDate = dueDate.AddDate(0, 0, 1)
			if itemDueDate.Equal(startDate) || (itemDueDate.After(startDate) && itemDueDate.Before(endDate)) {
				return true
			}
		}
		return false
	case DUE_BEFORE:
		if itemDueDate.Before(dueDate) {
			return true
		}
		return false
	case DUE_AFTER:
		endDateTime := dueDate
		if allDay {
			endDateTime = dueDate.AddDate(0, 0, 1).Add(-time.Duration(time.Microsecond))
		}
		if itemDueDate.After(endDateTime) {
			return true
		}
		return false
	default:
		return false
	}
}

func EvalLabel(e LabelExpr, item todoist.Item, c *todoist.Client) (result bool) {
	for _, l := range c.Store.Labels {
		for _, id := range item.LabelIDs {
			if id == l.ID && e.label == l.Name {
				return true
			}
		}
	}
	return false
}

func EvalProject(e ProjectExpr, item todoist.Item, c *todoist.Client) (result bool) {
	for _, p := range c.Store.Projects {
		if item.ProjectID == p.GetID() && e.project == p.Name {
			return true
		}
	}
	return false
}

func EvalAsPriority(e StringExpr, item todoist.Item) (result bool) {
	matched := priorityRegex.FindStringSubmatch(e.literal)
	if len(matched) == 0 {
		return false
	} else {
		p, _ := strconv.Atoi(matched[1])
		if p == item.Priority {
			return true
		}
	}
	return false
}

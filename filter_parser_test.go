package main

import (
	"testing"

	"time"

	"github.com/stretchr/testify/assert"
)

// Test ...
func TestFilter(t *testing.T) {
	assert.Equal(t, nil, Filter(""), "they should be equal")
}

func TestPriorityFilter(t *testing.T) {
	assert.Equal(t, StringExpr{literal: "p1"}, Filter("p1"), "they should be equal")
}

func TestLabelFilter(t *testing.T) {
	assert.Equal(t, LabelExpr{label: "label"}, Filter("@label"), "they should be equal")
}

func TestProjectFilter(t *testing.T) {
	assert.Equal(t, ProjectExpr{project: "project"}, Filter("#project"), "they should be equal")
}

func TestBoolInfixFilter(t *testing.T) {
	assert.Equal(t,
		BoolInfixOpExpr{
			left:     StringExpr{literal: "p1"},
			operator: '|',
			right:    StringExpr{literal: "p2"},
		},
		Filter("p1 | p2"), "they should be equal")

	assert.Equal(t,
		BoolInfixOpExpr{
			left:     StringExpr{literal: "p1"},
			operator: '&',
			right:    StringExpr{literal: "p2"},
		},
		Filter("p1 & p2"), "they should be equal")

	assert.Equal(t,
		BoolInfixOpExpr{
			left:     StringExpr{literal: "p1"},
			operator: '&',
			right: BoolInfixOpExpr{
				left:     StringExpr{literal: "p2"},
				operator: '|',
				right:    StringExpr{literal: "p3"},
			},
		},
		Filter("p1 & (p2 | p3 )"), "they should be equal")
}

func setNow(t time.Time) {
	now = func() time.Time { return t }
}

func TestDateTimeFilter(t *testing.T) {
	timeNow := time.Date(2017, time.January, 2, 1, 0, 0, 0, time.Local)
	setNow(timeNow)

	assert.Equal(t,
		DueDateExpr{operation: DUE_ON, datetime: time.Date(2017, time.October, 5, 0, 0, 0, 0, time.Local), allDay: true},
		Filter("10/5/2017"), "they should be equal")

	assert.Equal(t,
		DueDateExpr{operation: DUE_ON, datetime: time.Date(timeNow.Year(), time.January, 3, 0, 0, 0, 0, time.Local), allDay: true},
		Filter("Jan 3"), "they should be equal")

	assert.Equal(t,
		DueDateExpr{operation: DUE_ON, datetime: time.Date(timeNow.Year(), time.August, 8, 0, 0, 0, 0, time.Local), allDay: true},
		Filter("8 August"), "they should be equal")

	assert.Equal(t,
		DueDateExpr{operation: DUE_ON, datetime: time.Date(2020, time.February, 10, 0, 0, 0, 0, time.Local), allDay: true},
		Filter("10 Feb 2020"), "they should be equal")

	assert.Equal(t,
		DueDateExpr{operation: DUE_ON, datetime: time.Date(timeNow.Year(), time.May, 16, 0, 0, 0, 0, time.Local), allDay: true},
		Filter("16/05"), "they should be equal")

	assert.Equal(t,
		DueDateExpr{operation: DUE_ON, datetime: time.Date(timeNow.Year(), timeNow.Month(), timeNow.Day(), 16, 0, 0, 0, time.Local), allDay: false},
		Filter("16:00"), "they should be equal")

	assert.Equal(t,
		DueDateExpr{operation: DUE_ON, datetime: time.Date(timeNow.Year(), timeNow.Month(), timeNow.Day(), 16, 10, 3, 0, time.Local), allDay: false},
		Filter("16:10:03"), "they should be equal")

	assert.Equal(t,
		DueDateExpr{operation: DUE_ON, datetime: time.Date(timeNow.Year(), timeNow.Month(), timeNow.Day(), 15, 0, 0, 0, time.Local), allDay: false},
		Filter("3pm"), "they should be equal")

	assert.Equal(t,
		DueDateExpr{operation: DUE_ON, datetime: time.Date(timeNow.Year(), timeNow.Month(), timeNow.Day(), 7, 0, 0, 0, time.Local), allDay: false},
		Filter("7am"), "they should be equal")

	assert.Equal(t,
		DueDateExpr{operation: DUE_ON, datetime: time.Date(2020, time.February, 10, 15, 0, 0, 0, time.Local), allDay: false},
		Filter("10 Feb 2020 3pm"), "they should be equal")

	assert.Equal(t,
		DueDateExpr{operation: DUE_ON, datetime: time.Date(timeNow.Year(), timeNow.Month(), timeNow.Day(), 7, 0, 0, 0, time.Local), allDay: false},
		Filter("7am"), "they should be equal")
}

func TestSpecialDateTimeFilter(t *testing.T) {
	timeNow := time.Date(2017, time.January, 1, 1, 0, 0, 0, time.Local)
	setNow(timeNow)
	assert.Equal(t,
		DueDateExpr{operation: DUE_ON, datetime: time.Date(2017, time.January, 1, 0, 0, 0, 0, time.Local), allDay: true},
		Filter("today"), "they should be equal")

	assert.Equal(t,
		DueDateExpr{operation: DUE_ON, datetime: time.Date(2017, time.January, 1, 0, 0, 0, 0, time.Local), allDay: true},
		Filter("tod"), "they should be equal")

	assert.Equal(t,
		DueDateExpr{operation: DUE_ON, datetime: time.Date(2017, time.January, 1, 0, 0, 0, 0, time.Local), allDay: true},
		Filter("Today"), "they should be equal")

	assert.Equal(t,
		DueDateExpr{operation: DUE_ON, datetime: time.Date(2017, time.January, 2, 0, 0, 0, 0, time.Local), allDay: true},
		Filter("tomorrow"), "they should be equal")

	assert.Equal(t,
		DueDateExpr{operation: DUE_ON, datetime: time.Date(2016, time.December, 31, 0, 0, 0, 0, time.Local), allDay: true},
		Filter("yesterday"), "they should be equal")
}

func TestDateTimeElapsedFilter(t *testing.T) {
	timeNow := time.Date(2017, time.January, 2, 18, 0, 0, 0, time.Local)
	setNow(timeNow)
	assert.Equal(t,
		DueDateExpr{operation: DUE_ON, datetime: time.Date(timeNow.Year(), timeNow.Month(), timeNow.Day()+1, 16, 0, 0, 0, time.Local), allDay: false},
		Filter("16:00"), "they should be equal")

	timeNow = time.Date(2017, time.May, 16, 23, 59, 59, 0, time.Local)
	setNow(timeNow)
	assert.Equal(t,
		DueDateExpr{operation: DUE_ON, datetime: time.Date(timeNow.Year(), time.May, 16, 0, 0, 0, 0, time.Local), allDay: true},
		Filter("16/05"), "they should be equal")

	timeNow = time.Date(2017, time.May, 17, 0, 0, 0, 0, time.Local)
	setNow(timeNow)
	assert.Equal(t,
		DueDateExpr{operation: DUE_ON, datetime: time.Date(timeNow.Year()+1, time.May, 16, 0, 0, 0, 0, time.Local), allDay: true},
		Filter("16/05"), "they should be equal")
}

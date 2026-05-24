package models

import (
	"fmt"
	"time"
)

type Keymap string

type Keymaps map[Keymap]struct{}

func newKeymap(value string, keymaps Keymaps) (Keymap, error) {
	if value == "" {
		return "", ValidationError{
			Field:   "keymap",
			Message: "keymap cannot be blank",
		}
	}

	keymap := Keymap(value)
	_, ok := keymaps[keymap]
	if !ok {
		return "", ValidationError{
			Field:   "keymap",
			Message: fmt.Sprintf("keymap '%s' is not supported", value),
		}
	}

	return keymap, nil
}

type LessonStepResult struct {
	ID        ID
	UserID    ID
	Keymap    Keymap
	LessonID  uint32
	StepID    uint32
	WPM       WPM
	CPM       CPM
	Accuracy  Accuracy
	Duration  time.Duration
	CreatedAt time.Time
}

type LessonStepResultOptions struct {
	UserID    ID
	Keymap    string
	LessonID  uint32
	StepID    uint32
	WPM       float64
	CPM       float64
	Accuracy  float64
	Duration  time.Duration
	CreatedAt time.Time
}

func NewLessonStepResult(opts LessonStepResultOptions, keymaps Keymaps) (*LessonStepResult, error) {
	keymap, err := newKeymap(opts.Keymap, keymaps)
	if err != nil {
		return nil, err
	}

	wpm, err := newWPM(opts.WPM)
	if err != nil {
		return nil, err
	}

	cpm, err := newCPM(opts.CPM)
	if err != nil {
		return nil, err
	}

	accuracy, err := newAccuracy(opts.Accuracy)
	if err != nil {
		return nil, err
	}

	duration, err := newDuration(opts.Duration)
	if err != nil {
		return nil, err
	}

	return &LessonStepResult{
		UserID:    opts.UserID,
		Keymap:    keymap,
		LessonID:  opts.LessonID,
		StepID:    opts.StepID,
		WPM:       wpm,
		CPM:       cpm,
		Accuracy:  accuracy,
		Duration:  duration,
		CreatedAt: opts.CreatedAt,
	}, nil
}

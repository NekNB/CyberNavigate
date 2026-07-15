package models

import "time"

type SessionData struct {
	UUID          string
	CreatedAt     time.Time
	CurrentStepId string
	CurrentTrust  int
	FinishedAt    *time.Time
}

type ScenarioData struct {
	UUID        string
	Title       string
	Description string
	Difficulty  string
	FirstStep   string
	ArticleIds  *[]string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type StepData struct {
	UUID           string
	PreviousStep   *string
	PreviousAnswer *string
	MinTrust       *int
	MaxTrust       *int
	ScenarioId     string
	ActionIds      *[]string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type ActionData struct {
	UUID string
	// StepId    *string
	Type      string
	MessageId string
	Delay     int
}

type MessageData struct {
	UUID       string
	SenderId   string
	SenderName string
	Text       *string
}

type AnswerData struct {
	UUID     string
	Text     string
	AddTrust int
	Error    *string
}

type FileData struct {
	UUID     string
	Filename string
	IsSafe   bool
	Size     int
	Error    *string
}

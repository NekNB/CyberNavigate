package http

import (
	"fmt"
)

type Scenario struct {
	Id          string `json:"id"`
	Difficulty  string `json:"difficulty"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

func (a *APIClient) CreateScenario(
	title, difficulty, description string,
) (*Response[Scenario], error) {
	var resp Scenario
	r, err := a.apiClient.R().SetSuccessResult(&resp).SetBody(map[string]string{
		"difficulty":  difficulty,
		"title":       title,
		"description": description,
	}).Post("/simulator/scenarios")
	if err != nil {
		a.log.Warn(err)
		return nil, err
	}
	if !r.IsSuccessState() {
		a.log.Warn(r.StatusCode)
		return nil, fmt.Errorf("Неизвестная ошибка %d", r.StatusCode)
	}

	return &Response[Scenario]{
		Code: r.StatusCode,
		Body: &resp,
	}, nil
}

type Message struct {
	Id         string    `json:"id"`
	SenderId   *string   `json:"senderId"`
	SenderName string    `json:"senderName"`
	Text       *string   `json:"text"`
	Files      *[]File   `json:"files"`
	Answers    *[]Answer `json:"answers"`
}

func (a *APIClient) CreateMessage(
	senderId, text *string,
	senderName string,
	files *[]string,
	answers *[]string,
) (*Response[Message], error) {
	var resp Message
	req := struct {
		SenderId   *string   `json:"senderId,omitempty"`
		SenderName string    `json:"senderName"`
		Text       *string   `json:"text,omitempty"`
		Files      *[]string `json:"files,omitempty"`
		Answers    *[]string `json:"answers,omitempty"`
	}{
		SenderId:   senderId,
		SenderName: senderName,
		Text:       text,
		Files:      files,
		Answers:    answers,
	}

	r, err := a.apiClient.R().SetSuccessResult(&resp).SetBody(req).Post("/simulator/messages")
	if err != nil {
		a.log.Warn(err)
		return nil, err
	}
	if !r.IsSuccessState() {
		a.log.Warn(r.StatusCode)
		return nil, fmt.Errorf("Неизвестная ошибка %d", r.StatusCode)
	}

	return &Response[Message]{
		Code: r.StatusCode,
		Body: &resp,
	}, nil
}

type File struct {
	Id       string  `json:"id"`
	Filename string  `json:"filename"`
	IsSafe   bool    `json:"isSafe"`
	Error    *string `json:"error"`
	Size     int     `json:"size"`
}

func (a *APIClient) CreateFile(
	Filename string,
	IsSafe bool,
	Error *string,
	Size int,
) (*Response[Message], error) {
	var resp Message
	req := struct {
		Filename string  `json:"filename"`
		IsSafe   bool    `json:"isSafe"`
		Error    *string `json:"error,omitempty"`
		Size     int     `json:"size"`
	}{
		Filename: Filename,
		IsSafe:   IsSafe,
		Error:    Error,
		Size:     Size,
	}

	r, err := a.apiClient.R().SetSuccessResult(&resp).SetBody(req).Post("/simulator/files")
	if err != nil {
		a.log.Warn(err)
		return nil, err
	}
	if !r.IsSuccessState() {
		a.log.Warn(r.StatusCode)
		return nil, fmt.Errorf("Неизвестная ошибка %d", r.StatusCode)
	}

	return &Response[Message]{
		Code: r.StatusCode,
		Body: &resp,
	}, nil
}

type Answer struct {
	Id       string  `json:"id"`
	Text     string  `json:"text"`
	Error    *string `json:"error"`
	AddTrust int     `json:"addTrust"`
}

func (a *APIClient) CreateAnswer(
	Text string,
	Error *string,
	AddTrust int,
) (*Response[Message], error) {
	var resp Message
	req := struct {
		Text     string  `json:"text"`
		Error    *string `json:"error,omitempty"`
		AddTrust int     `json:"addTrust"`
	}{
		Text:     Text,
		Error:    Error,
		AddTrust: AddTrust,
	}

	r, err := a.apiClient.R().SetSuccessResult(&resp).SetBody(req).Post("/simulator/answers")
	if err != nil {
		a.log.Warn(err)
		return nil, err
	}
	if !r.IsSuccessState() {
		a.log.Warn(r.StatusCode)
		return nil, fmt.Errorf("Неизвестная ошибка %d", r.StatusCode)
	}

	return &Response[Message]{
		Code: r.StatusCode,
		Body: &resp,
	}, nil
}

type Step struct {
	Id             string   `json:"id"`
	ScenarioId     string   `json:"scenarioId"`
	PreviousStep   *string  `json:"previousStep"`
	PreviousAnswer *string  `json:"previousAnswer"`
	MinTrust       *int     `json:"minTrust"`
	MaxTrust       *int     `json:"maxTrust"`
	Actions        []Action `json:"actions"`
}

func (a *APIClient) CreateStep(
	ScenarioId string,
	PreviousStep *[]string,
	PreviousAnswer *[]string,
	MinTrust *int,
	MaxTrust *int,
	Actions []string,
) (*Response[Message], error) {
	var resp Message
	req := struct {
		ScenarioId      string    `json:"scenarioId"`
		PreviousSteps   *[]string `json:"previousSteps,omitempty"`
		PreviousAnswers *[]string `json:"previousAnswers,omitempty"`
		MinTrust        *int      `json:"minTrust,omitempty"`
		MaxTrust        *int      `json:"maxTrust,omitempty"`
		Actions         []string  `json:"actions,omitempty"`
	}{
		ScenarioId:      ScenarioId,
		PreviousSteps:   PreviousStep,
		PreviousAnswers: PreviousAnswer,
		MinTrust:        MinTrust,
		MaxTrust:        MaxTrust,
		Actions:         Actions,
	}

	r, err := a.apiClient.R().SetSuccessResult(&resp).SetBody(req).Post("/simulator/step")
	if err != nil {
		a.log.Warn(err)
		return nil, err
	}
	if !r.IsSuccessState() {
		a.log.Warn(r.StatusCode)
		return nil, fmt.Errorf("Неизвестная ошибка %d", r.StatusCode)
	}

	return &Response[Message]{
		Code: r.StatusCode,
		Body: &resp,
	}, nil
}

type Action struct {
	Id      string  `json:"id"`
	Type    string  `json:"type"`
	Message Message `json:"message"`
	Delay   int     `json:"delay"`
}

func (a *APIClient) CreateAction(
	Type string,
	MessageId string,
	Delay int,
) (*Response[Action], error) {
	var resp Action
	req := struct {
		Type      string `json:"type"`
		MessageId string `json:"message"`
		Delay     int    `json:"delay"`
	}{
		Type:      Type,
		MessageId: MessageId,
		Delay:     Delay,
	}

	r, err := a.apiClient.R().SetSuccessResult(&resp).SetBody(req).Post("/simulator/action")
	if err != nil {
		a.log.Warn(err)
		return nil, err
	}
	if !r.IsSuccessState() {
		a.log.Warn(r.StatusCode)
		return nil, fmt.Errorf("Неизвестная ошибка %d", r.StatusCode)
	}

	return &Response[Action]{
		Code: r.StatusCode,
		Body: &resp,
	}, nil
}

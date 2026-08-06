package scenarios

import (
	"fmt"
	"strings"

	"github.com/NekNB/CyberNavigate/init/internal/http"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

// ==========================================
// 1. СТРУКТУРЫ ДЛЯ ПАРСИНГА YAML
// ==========================================

type Config struct {
	Metadata   Metadata          `yaml:"metadata"`
	StartsStep string            `yaml:"starts_step"`
	Senders    map[string]Sender `yaml:"senders"`
	Steps      map[string]Step   `yaml:"steps"`
	Ends       []string          `yaml:"ends"`
}

type Metadata struct {
	Title       string `yaml:"title"`
	Difficulty  string `yaml:"difficulty"`
	Description string `yaml:"description"`
}

type Sender struct {
	SenderName string `yaml:"senderName"`
}

type Step struct {
	NextStep string   `yaml:"next_step,omitempty"`
	Actions  []Action `yaml:"actions"`
	MinTrust *int     `yaml:"min_trust,omitempty"`
	MaxTrust *int     `yaml:"max_trust,omitempty"` // Используем указатель для парсинга YAML, чтобы отличать 0 от отсутствующего поля
}

type Action struct {
	Type          string   `yaml:"type"`
	SenderName    string   `yaml:"senderName,omitempty"`
	SenderID      string   `yaml:"sender_id,omitempty"`
	SenderNameRef string   `yaml:"sender_name,omitempty"`
	Text          string   `yaml:"text"`
	Delay         int      `yaml:"delay"`
	Answers       []Answer `yaml:"answers,omitempty"`
	Files         []File   `yaml:"files,omitempty"`
}

type Answer struct {
	Text     string `yaml:"text"`
	AddTrust int    `yaml:"add_trust"`
	NextStep string `yaml:"next_step,omitempty"`
}

type File struct {
	Filename string `yaml:"filename"`
	Error    string `yaml:"error"`
	Size     int    `yaml:"size"`
	IsSafe   bool   `yaml:"is_safe"`
}

// ==========================================
// 2. ВСПОМОГАТЕЛЬНЫЕ ФУНКЦИИ
// ==========================================

// ptrString возвращает указатель на строку, если она не пустая, иначе nil
func ptrString(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// isEnd проверяет, является ли шаг финальным
func isEnd(step string, ends []string) bool {
	for _, e := range ends {
		if e == step {
			return true
		}
	}
	return false
}

// ==========================================
// 3. ОСНОВНОЙ СЦЕНАРИЙ ГЕНЕРАЦИИ
// ==========================================

func GenerateScenario(yamlData []byte, apiClient *http.APIClient, log *logrus.Logger) error {
	var config Config
	if err := yaml.Unmarshal(yamlData, &config); err != nil {
		return fmt.Errorf("ошибка парсинга YAML: %w", err)
	}

	// 1. Создаем сам сценарий
	scenarioResp, err := apiClient.CreateScenario(config.Metadata.Title, config.Metadata.Difficulty, config.Metadata.Description)
	if err != nil || scenarioResp.Body == nil {
		return fmt.Errorf("ошибка создания сценария: %w", err)
	}
	scenarioID := scenarioResp.Body.Id

	// Кэш для хранения senderId по внутреннему имени (например, "alex")
	senderCache := make(map[string]string)

	// Хранилища для связей
	actionIDs := make(map[string][]string)   // stepKey -> []actionID
	answerLinks := make(map[string][]string) // nextStepKey -> []answerID
	stepLinks := make(map[string]string)     // nextStepKey -> previousStepKey

	// Предварительно собираем связи шагов через поле next_step на уровне Step
	for stepKey, stepData := range config.Steps {
		if stepData.NextStep != "" && !isEnd(stepData.NextStep, config.Ends) {
			stepLinks[stepData.NextStep] = stepKey
		}
	}

	// 2. ИТЕРАТИВНОЕ СОЗДАНИЕ: Files -> Answers -> Messages -> Actions
	for stepKey, stepData := range config.Steps {
		for _, action := range stepData.Actions {

			// 2.1 Создаем Files
			var fileIDs []string
			for _, f := range action.Files {
				// В новом API CreateFile возвращает *Response[Message], но ID находится там же
				fileResp, err := apiClient.CreateFile(f.Filename, f.IsSafe, ptrString(f.Error), f.Size)
				if err != nil || fileResp.Body == nil {
					return fmt.Errorf("ошибка создания файла %s: %w", f.Filename, err)
				}
				fileIDs = append(fileIDs, fileResp.Body.Id)
			}

			// 2.2 Создаем Answers
			var answerIDs []string
			for _, ans := range action.Answers {
				// CreateAnswer тоже возвращает *Response[Message]
				ansResp, err := apiClient.CreateAnswer(ans.Text, nil, ans.AddTrust)
				if err != nil || ansResp.Body == nil {
					return fmt.Errorf("ошибка создания ответа '%s': %w", ans.Text, err)
				}
				ansID := ansResp.Body.Id
				answerIDs = append(answerIDs, ansID)

				// Сохраняем связь ответа со следующим шагом
				if ans.NextStep != "" && !isEnd(ans.NextStep, config.Ends) {
					answerLinks[ans.NextStep] = append(answerLinks[ans.NextStep], ansID)
				}
			}

			// 2.3 Подготовка данных для Message
			var senderIDPtr *string
			senderName := ""

			switch action.Type {
			case "sms":
				senderName = action.SenderName
			case "message":
				// Получаем senderName из конфига senders по ссылке (например "alex.senderName")
				parts := strings.Split(action.SenderNameRef, ".")
				if len(parts) == 2 {
					if sender, ok := config.Senders[parts[0]]; ok {
						senderName = sender.SenderName
					}
				}

				// Проверяем кэш для senderId
				if cachedID, ok := senderCache[action.SenderID]; ok {
					senderIDPtr = &cachedID
				}
			}

			// Подготавливаем указатели для слайсов
			var filesPtr *[]string
			if len(fileIDs) > 0 {
				filesPtr = &fileIDs
			}
			var answersPtr *[]string
			if len(answerIDs) > 0 {
				answersPtr = &answerIDs
			}

			textPtr := &action.Text

			// Создаем Message
			msgResp, err := apiClient.CreateMessage(senderIDPtr, textPtr, senderName, filesPtr, answersPtr)
			if err != nil || msgResp.Body == nil {
				return fmt.Errorf("ошибка создания сообщения: %w", err)
			}
			msgID := msgResp.Body.Id

			// Обновляем кэш отправителей для message
			if action.Type == "message" && senderIDPtr == nil && msgResp.Body.SenderId != nil {
				senderCache[action.SenderID] = *msgResp.Body.SenderId
			}

			// 2.4 Создаем Action
			actResp, err := apiClient.CreateAction(action.Type, msgID, action.Delay)
			if err != nil || actResp.Body == nil {
				return fmt.Errorf("ошибка создания action: %w", err)
			}
			actionIDs[stepKey] = append(actionIDs[stepKey], actResp.Body.Id)
		}
	}

	// 3. ИТЕРАТИВНОЕ СОЗДАНИЕ STEPS
	var stepAPIMap = make(map[string][]string) // stepKey -> API steps ID
	var createdSteps = make(map[string]bool)

	var createStepFunc func(stepKey string) error
	createStepFunc = func(stepKey string) error {
		if _, ok := createdSteps[stepKey]; ok {
			return nil // Уже создан
		}
		if isEnd(stepKey, config.Ends) {
			return nil // Финальные шаги не создаем
		}

		stepData := config.Steps[stepKey]
		var prevSteps, prevAnsPtr *[]string

		// Если это не первый шаг, ищем связи
		if stepKey != config.StartsStep {
			// Приоритет у связи через Answer
			if ansIDs, ok := answerLinks[stepKey]; ok && len(ansIDs) > 0 {
				prevAnsPtr = &ansIDs
			} else if prevYAMLStep, ok := stepLinks[stepKey]; ok {
				// Иначе ищем связь через Step.next_step
				if _, ok := stepAPIMap[prevYAMLStep]; !ok {
					if err := createStepFunc(prevYAMLStep); err != nil {
						return err
					}
				}
				prevStepsFromMap := stepAPIMap[prevYAMLStep]
				prevSteps = &prevStepsFromMap
			}
		}

		// MinTrust не указываем (nil), MaxTrust берем из распарсенного YAML (он уже *int)
		stepResp, err := apiClient.CreateStep(scenarioID, prevSteps, prevAnsPtr, stepData.MinTrust, stepData.MaxTrust, actionIDs[stepKey])
		if err != nil || stepResp.Body == nil {
			return fmt.Errorf("ошибка создания шага %s: %w", stepKey, err)
		}

		stepAPIMap[stepKey] = append(stepAPIMap[stepKey], stepResp.Body.Id)
		createdSteps[stepKey] = true

		return nil
	}

	// Запускаем создание шагов
	for stepKey := range config.Steps {
		if err := createStepFunc(stepKey); err != nil {
			return err
		}
	}

	fmt.Println("Сценарий успешно сгенерирован!")
	return nil
}

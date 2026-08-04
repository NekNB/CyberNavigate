// package main

// import (
// 	"flag"
// 	"fmt"
// 	"os"

// 	"github.com/imroc/req/v3"
// )

// func getGatewayPath() string {
// 	var res string
// 	flag.StringVar(&res, "url", "http://127.0.0.1:9080/api/v1", "")
// 	flag.Parse()

// 	return res
// }

// func main() {
// 	var baseUrl string
// 	var startStep int
// 	flag.StringVar(&baseUrl, "url", "http://127.0.0.1:9080/api/v1", "")
// 	flag.IntVar(&startStep, "step", 1, "")
// 	flag.Parse()
// 	args := os.Args[1:]

// 	client := req.C()
// 	client.BaseURL = baseUrl

// 	_ = args

// 	for i := startStep; i <= 3; i++ {
// 		switch i {
// 		case 1:
// 			step1(client)
// 		case 2:
// 			step2(client)
// 		}
// 	}

// }

// func step1(client *req.Client) {
// 	fmt.Println("STEP 1")

// 	resp, err := client.R().
// 		SetBody(map[string]string{"username": "admin", "password": "admin"}).Post("/auth/login")
// 	if err != nil {
// 		fmt.Println(err)
// 		os.Exit(1)
// 	}
// 	if resp.StatusCode != 200 {
// 		fmt.Println(resp.StatusCode)
// 		os.Exit(1)
// 	}
// }

// func step2(client *req.Client) {
// 	fmt.Println("STEP 2")

//		resp, err := client.R().
//			SetBody(map[string]string{
//				"difficulty":  "middle",
//				"title":       "Служба безопасности госуслуг",
//				"description": "Учимся распозновать обман и не устанавливать вредоносные файлы",
//			}).Post("/simulator/scenarios")
//		if err != nil {
//			fmt.Println(err)
//			os.Exit(1)
//		}
//		if resp.StatusCode != 201 {
//			fmt.Println(resp.StatusCode, resp)
//			os.Exit(1)
//		}
//	}
package main

import (
	"fmt"
	"log"

	"github.com/imroc/req/v3"
)

// Структуры для создания сущностей (согласно спецификации)
type ScenarioRequest struct {
	Difficulty  string   `json:"difficulty"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	ArticleIDs  []string `json:"articleIds,omitempty"`
}

type FileRequest struct {
	Filename string `json:"filename"`
	IsSafe   bool   `json:"isSafe"`
	Error    string `json:"error,omitempty"`
	Size     int    `json:"size"`
}

type AnswerRequest struct {
	Text     string `json:"text"`
	Error    string `json:"error,omitempty"`
	AddTrust int    `json:"addTrust"`
}

type MessageRequest struct {
	SenderName string   `json:"senderName"`
	Text       string   `json:"text"`
	Files      []string `json:"files,omitempty"`
	Answers    []string `json:"answers,omitempty"`
}

type ActionRequest struct {
	Type    string `json:"type"`
	Message string `json:"message"`
	Delay   int    `json:"delay"`
}

type StepRequest struct {
	ScenarioID string   `json:"scenarioId"`
	MinTrust   int      `json:"minTrust"`
	MaxTrust   int      `json:"maxTrust"`
	Actions    []string `json:"actions"`
}

// Универсальная структура ответа для получения ID
type IDResponse struct {
	ID string `json:"id"`
}

const baseURL = "http://127.0.0.1:9080/api/v1"

func main() {
	client := req.C().SetBaseURL(baseURL)
	client.R().
		SetBody(map[string]string{"username": "admin", "password": "admin"}).Post("/auth/login")
	// Если для админских эндпоинтов нужен Bearer токен, раскомментируйте строку ниже и вставьте токен
	// client.SetCommonBearerTokenAuth("YOUR_ADMIN_JWT_TOKEN")

	log.Println("Начинаем создание сценария...")

	// 1. Создание сценария
	scenarioID, err := createEntity(client, "/simulator/scenarios", ScenarioRequest{
		Difficulty:  "middle",
		Title:       "Звонок из службы безопасности",
		Description: "Злоумышленник представляется сотрудником банка. Задача игрока — не довериться и не скачать вредоносный файл.",
	})
	if err != nil {
		log.Fatalf("Ошибка создания сценария: %v", err)
	}
	log.Printf("Сценарий создан: %s\n", scenarioID)

	// 2. Создание файла (Ловушка)
	fileID, err := createEntity(client, "/simulator/files", FileRequest{
		Filename: "bank_security_update.apk",
		IsSafe:   false,
		Error:    "Вы скачали вредоносное ПО. Мошенники получили доступ к вашему онлайн-банку.",
		Size:     15,
	})
	if err != nil {
		log.Fatalf("Ошибка создания файла: %v", err)
	}
	log.Printf("Файл создан: %s\n", fileID)

	// 3. Создание ответов
	// Ответы для первого сообщения
	ansTrust1, _ := createEntity(client, "/simulator/answers", AnswerRequest{
		Text:     "Ой, помогите! Что мне делать?",
		AddTrust: 40, // Игрок доверился
	})
	ansSusp1, _ := createEntity(client, "/simulator/answers", AnswerRequest{
		Text:     "Я сам перезвоню в банк по официальному номеру.",
		AddTrust: -50, // Игрок заподозрил неладное
	})

	// Ответы для второго сообщения (где присылается файл)
	ansDownload, _ := createEntity(client, "/simulator/answers", AnswerRequest{
		Text:     "Хорошо, я скачиваю приложение.",
		AddTrust: 50, // Игрок полностью повелся
		Error:    "Никогда не скачивайте приложения по ссылкам из сообщений и звонков.",
	})
	ansRefuse, _ := createEntity(client, "/simulator/answers", AnswerRequest{
		Text:     "Я не буду ничего скачивать. Вы мошенник!",
		AddTrust: -60, // Игрок отказался и раскрыл мошенника
	})

	log.Println("Ответы созданы.")

	// 4. Создание сообщений
	// Сообщение 1: Начало
	msg1ID, _ := createEntity(client, "/simulator/messages", MessageRequest{
		SenderName: "Служба безопасности Банка",
		Text:       "Здравствуйте! Мы зафиксировали попытку списания 50 000 рублей с вашей карты. Это вы совершаете операцию?",
		Answers:    []string{ansTrust1, ansSusp1},
	})

	// Сообщение 2: Присылка файла (появляется, если доверие > 0)
	msg2ID, _ := createEntity(client, "/simulator/messages", MessageRequest{
		SenderName: "Служба безопасности Банка",
		Text:       "Чтобы заблокировать перевод, вам нужно срочно установить наше защитное приложение. Скачайте его по ссылке ниже.",
		Files:      []string{fileID}, // Прикрепляем файл-ловушку
		Answers:    []string{ansDownload, ansRefuse},
	})

	// Сообщение 3: Good End (Если игрок был подозрителен)
	msgGoodID, _ := createEntity(client, "/simulator/messages", MessageRequest{
		SenderName: "Служба безопасности Банка",
		Text:       "Алло? Алло?... *мошенник бросил трубку*",
		// Нет ответов, конец игры
	})

	// Сообщение 4: Game Over (Если игрок скачал файл)
	msgBadID, _ := createEntity(client, "/simulator/messages", MessageRequest{
		SenderName: "Служба безопасности Банка",
		Text:       "Отлично! Теперь ваши данные у нас. Спасибо за сотрудничество!",
		// Нет ответов, конец игры
	})

	log.Println("Сообщения созданы.")

	// 5. Создание действий (Actions)
	act1ID, _ := createEntity(client, "/simulator/action", ActionRequest{
		Type:    "sms",
		Message: msg1ID,
		Delay:   5,
	})
	act2ID, _ := createEntity(client, "/simulator/action", ActionRequest{
		Type:    "message",
		Message: msg2ID,
		Delay:   10,
	})
	actGoodID, _ := createEntity(client, "/simulator/action", ActionRequest{
		Type:    "message",
		Message: msgGoodID,
		Delay:   5,
	})
	actBadID, _ := createEntity(client, "/simulator/action", ActionRequest{
		Type:    "message",
		Message: msgBadID,
		Delay:   5,
	})

	log.Println("Действия созданы.")

	// 6. Создание шагов (Управление доверием и ветвлением)
	// Шаг 1: Стартовый (Доступен при любом доверии, с которого начинается игра)
	createStep(client, StepRequest{
		ScenarioID: scenarioID,
		MinTrust:   -100,
		MaxTrust:   100,
		Actions:    []string{act1ID},
	})

	// Шаг 2: Мошенник присылает файл. Доступен только если игрок доверился (Доверие от 1 до 100)
	createStep(client, StepRequest{
		ScenarioID: scenarioID,
		MinTrust:   1,
		MaxTrust:   100,
		Actions:    []string{act2ID},
	})

	// Шаг 3 (Good End): Игрок раскрыл мошенника. Доступен, если доверие упало (от -100 до 0)
	createStep(client, StepRequest{
		ScenarioID: scenarioID,
		MinTrust:   -100,
		MaxTrust:   0,
		Actions:    []string{actGoodID},
	})

	// Шаг 4 (Game Over): Игрок скачал файл (Доверие сильно выросло, от 50 до 100)
	createStep(client, StepRequest{
		ScenarioID: scenarioID,
		MinTrust:   50,
		MaxTrust:   100,
		Actions:    []string{actBadID},
	})

	log.Println("Шаги созданы.")
	log.Println("✅ Сценарий успешно развернут на сервере!")
}

// Вспомогательная функция для отправки POST запросов и получения ID созданной сущности
func createEntity(client *req.Client, path string, payload interface{}) (string, error) {
	var resp IDResponse

	r, err := client.R().SetBody(payload).SetSuccessResult(&resp).Post(path)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}

	if r.StatusCode != 200 && r.StatusCode != 201 {
		return "", fmt.Errorf("bad status: %d, body: %s", r.StatusCode, r.String())
	}

	return resp.ID, nil
}

// Вспомогательная функция для создания шагов
func createStep(client *req.Client, payload StepRequest) {
	var resp IDResponse
	r, err := client.R().SetBody(payload).SetSuccessResult(&resp).Post("/simulator/step")
	if err != nil {
		log.Fatalf("Ошибка создания шага: %v", err)
	}
	if r.StatusCode != 201 {
		log.Fatalf("Ошибка создания шага. Статус: %d, Тело: %s", r.StatusCode, r.String())
	}
	log.Printf("Создан шаг: %s (MinTrust: %d, MaxTrust: %d)\n", resp.ID, payload.MinTrust, payload.MaxTrust)
}

package http

import (
	"fmt"
	"strings"

	"github.com/imroc/req/v3"
	"github.com/sirupsen/logrus"
)

type Response[T any] struct {
	Code int
	Body *T
}

type APIClient struct {
	log       *logrus.Logger
	apiClient *req.Client
}

// doRequest - вспомогательный метод для устранения дублирования кода
func (a *APIClient) doRequest(method, url string, body any, result any) (*req.Response, error) {
	var r *req.Response
	var err error

	// Явно создаем запрос, чтобы задать result
	req := a.apiClient.R().SetSuccessResult(result)

	switch strings.ToUpper(method) {
	case "GET":
		r, err = req.Get(url)
	case "POST":
		r, err = req.SetBodyJsonMarshal(body).Post(url)
	case "PUT":
		r, err = req.SetBodyJsonMarshal(body).Put(url)
	case "PATCH":
		r, err = req.SetBodyJsonMarshal(body).Patch(url)
	case "DELETE":
		r, err = req.SetBodyJsonMarshal(body).Delete(url)
	default:
		return nil, fmt.Errorf("неизвестный HTTP метод: %s", method)
	}

	if err != nil {
		a.log.Warnf("Ошибка запроса %s %s: %v", method, url, err)
		return nil, err
	}

	if !r.IsSuccessState() {
		err = fmt.Errorf("неизвестная ошибка %d от %s %s", r.StatusCode, method, url)
		a.log.Warn(err)
		return nil, err
	}

	return r, nil
}

func New(log *logrus.Logger, baseUrl string) *APIClient {

	apiClient := req.NewClient()

	apiClient.
		SetBaseURL(baseUrl).
		SetCommonRetryCount(1).
		SetCommonRetryCondition(func(resp *req.Response, err error) bool {
			if err != nil {
				return false
			}
			return resp.StatusCode == 401 && !strings.Contains(resp.Request.URL.Path, "/auth/refresh")
		}).
		SetCommonRetryHook(func(resp *req.Response, err error) {
			log.Info("Получили 401 пытаемся обновить токены")
			refreshResp, err := apiClient.R().Put("/auth/refresh")
			if err != nil {
				log.Errorf("Ошибка сети при обновлении куки: %v", err)
				return
			}
			if !refreshResp.IsSuccessState() {
				log.Errorf("Сервер вернул ошибку при обновлении куки: %d", refreshResp.StatusCode)
				return
			}
			log.Info("Куки успешно обновлены, повторяем исходный запрос...")
		})

	return &APIClient{
		log:       log,
		apiClient: apiClient,
	}
}

func (a *APIClient) Login(username, pasword string) error {
	resp, err := a.apiClient.R().SetBody(map[string]string{
		"username": username,
		"password": pasword,
	}).Post("/auth/login")
	if err != nil {
		a.log.Warn(err)
		return err
	}
	if !resp.IsSuccessState() {
		a.log.Warn(resp.StatusCode)
		return fmt.Errorf("Неизвестная ошибка %d", resp.StatusCode)
	}
	return nil
}

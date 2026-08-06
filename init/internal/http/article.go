package http

import (
	"fmt"

	"github.com/imroc/req/v3"
)

type ArticleResponse struct {
	Id     string `json:"id"`
	Title  string `json:"title"`
	Status string `json:"status"`
	Resp   *req.Response
}

func (a *APIClient) CreateArticle(title string) (*Response[ArticleResponse], error) {

	var resp ArticleResponse
	r, err := a.apiClient.R().SetBody(map[string]string{
		"articleTitle": title,
	}).SetSuccessResult(&resp).Post("/articles")
	if err != nil {
		a.log.Warn(err)
		return nil, err
	}
	if !r.IsSuccessState() {
		if r.StatusCode == 400 {
			return &Response[ArticleResponse]{Code: r.StatusCode, Body: nil}, nil
		}
		a.log.Warn(r.StatusCode)
		return nil, fmt.Errorf("Неизвестная ошибка %d", r.StatusCode)
	}

	return &Response[ArticleResponse]{Code: r.StatusCode, Body: &resp}, nil
}

func (a *APIClient) UploadAcrticleText(articleId, text string) (*Response[ArticleResponse], error) {
	var resp ArticleResponse
	r, err := a.apiClient.R().SetBody(text).SetSuccessResult(&resp).Post(fmt.Sprintf("/articles/%s/text", articleId))
	if err != nil {
		a.log.Warn(err)
		return nil, err
	}
	if !r.IsSuccessState() {
		a.log.Warn(r.StatusCode)
		return nil, fmt.Errorf("Неизвестная ошибка %d", r.StatusCode)
	}

	return &Response[ArticleResponse]{Code: r.StatusCode, Body: &resp}, nil
}

func (a *APIClient) UploadAcrticleData(articleId string, status, title *string) (*Response[ArticleResponse], error) {
	var resp ArticleResponse
	req := make(map[string]string)

	if status != nil {
		req["articleStatus"] = *status

	}
	if title != nil {
		req["articleTitle"] = *title
	}

	r, err := a.apiClient.R().SetBody(req).SetSuccessResult(&resp).Patch(fmt.Sprintf("/articles/%s", articleId))
	if err != nil {
		a.log.Warn(err)
		return nil, err
	}
	if !r.IsSuccessState() {
		a.log.Warn(r.StatusCode)
		return nil, fmt.Errorf("Неизвестная ошибка %d", r.StatusCode)
	}

	return &Response[ArticleResponse]{Code: r.StatusCode, Body: &resp}, nil
}

func (a *APIClient) GetArticleData(articleId string) (*Response[ArticleResponse], error) {
	var resp ArticleResponse
	r, err := a.apiClient.R().SetSuccessResult(&resp).Get(fmt.Sprintf("/articles/%s", articleId))
	if err != nil {
		a.log.Warn(err)
		return nil, err
	}
	if !r.IsSuccessState() {
		a.log.Warn(r.StatusCode)
		return nil, fmt.Errorf("Неизвестная ошибка %d", r.StatusCode)
	}

	return &Response[ArticleResponse]{Code: r.StatusCode, Body: &resp}, nil
}

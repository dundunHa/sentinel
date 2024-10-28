package api

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"net/http"
	"sentinel/server/service"
)

// API 结构体
type API struct {
	svc *service.Service
}

// NewAPI 创建新的 API 实例
func NewAPI(svc *service.Service) *API {
	return &API{svc: svc}
}

func (a *API) Routes() http.Handler {
	r := chi.NewRouter()

	// 定义 API 路由
	r.Get("/messages", a.GetMessages)

	return r
}

func (a *API) GetMessages(w http.ResponseWriter, r *http.Request) {
	messages, err := a.svc.GetMessages(0)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 返回消息列表
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(messages); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

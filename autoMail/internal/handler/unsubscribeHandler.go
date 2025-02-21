package handler

import (
	"net/http"

	"automail/autoMail/internal/logic"
	"automail/autoMail/internal/svc"
	"automail/autoMail/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func UnsubscribeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.Request
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewUnsubscribeLogic(r.Context(), svcCtx)
		resp, err := l.Unsubscribe(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

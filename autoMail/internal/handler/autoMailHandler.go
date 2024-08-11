package handler

import (
	"net/http"

	"automail/autoMail/internal/logic"
	"automail/autoMail/internal/svc"
	"automail/autoMail/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func AutoMailHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.Request
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewAutoMailLogic(r.Context(), svcCtx)
		resp, err := l.AutoMail(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

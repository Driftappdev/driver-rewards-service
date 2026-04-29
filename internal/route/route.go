package route

import (
	"net/http"

	"dift_backend_driver/driver-rewards-service/internal/handler"
)

func Register(mux *http.ServeMux, h *handler.RewardsHandler) {
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	mux.HandleFunc("/api/v1/driver/rewards", h.GetRewardProgress)
	mux.HandleFunc("/api/v1/driver/campaigns", h.GetCampaigns)
	mux.HandleFunc("/api/v1/driver/rank", h.GetRank)
	mux.HandleFunc("/api/v1/driver/rewards/redeem", h.RedeemReward)
	mux.HandleFunc("/api/v1/driver/campaigns/claim", h.ClaimCampaign)
}

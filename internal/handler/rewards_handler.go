package handler

import (
	"encoding/json"
	"net/http"

	"dift_backend_driver/driver-rewards-service/internal/dto"
	response "github.com/PlatformCore/engine-core/core/contracts/response"

	"dift_backend_driver/driver-rewards-service/internal/service"
)

type RewardsHandler struct{ svc *service.RewardsService }

func NewRewardsHandler(svc *service.RewardsService) *RewardsHandler {
	return &RewardsHandler{svc: svc}
}

func (h *RewardsHandler) GetRewardProgress(w http.ResponseWriter, r *http.Request) {
	driverID := r.URL.Query().Get("driver_id")
	if driverID == "" {
		writeJSON(w, http.StatusBadRequest, response.Envelope[any]{
			Error: &response.AppError{Code: "bad_request", Message: "driver_id required"},
		})
		return
	}
	writeJSON(w, http.StatusOK, response.Envelope[any]{Data: h.svc.GetRewardProgress(driverID)})
}

func (h *RewardsHandler) GetCampaigns(w http.ResponseWriter, r *http.Request) {
	driverID := r.URL.Query().Get("driver_id")
	if driverID == "" {
		writeJSON(w, http.StatusBadRequest, response.Envelope[any]{
			Error: &response.AppError{Code: "bad_request", Message: "driver_id required"},
		})
		return
	}
	writeJSON(w, http.StatusOK, response.Envelope[any]{Data: h.svc.GetCampaigns(driverID)})
}

func (h *RewardsHandler) GetRank(w http.ResponseWriter, r *http.Request) {
	driverID := r.URL.Query().Get("driver_id")
	if driverID == "" {
		writeJSON(w, http.StatusBadRequest, response.Envelope[any]{
			Error: &response.AppError{Code: "bad_request", Message: "driver_id required"},
		})
		return
	}
	writeJSON(w, http.StatusOK, response.Envelope[any]{Data: h.svc.GetRank(driverID)})
}

func (h *RewardsHandler) RedeemReward(w http.ResponseWriter, r *http.Request) {
	var req dto.RewardRedemptionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, response.Envelope[any]{Error: &response.AppError{Code: "bad_request", Message: err.Error()}})
		return
	}
	if req.DriverID == "" || req.MilestoneID == "" {
		writeJSON(w, http.StatusBadRequest, response.Envelope[any]{Error: &response.AppError{Code: "bad_request", Message: "driver_id and milestone_id required"}})
		return
	}
	out, err := h.svc.RedeemReward(req.DriverID, req.MilestoneID)
	if err != nil {
		writeJSON(w, http.StatusConflict, response.Envelope[any]{Error: &response.AppError{Code: "reward_redeem_failed", Message: err.Error()}})
		return
	}
	writeJSON(w, http.StatusOK, response.Envelope[any]{Data: out})
}

func (h *RewardsHandler) ClaimCampaign(w http.ResponseWriter, r *http.Request) {
	var req dto.CampaignClaimRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, response.Envelope[any]{Error: &response.AppError{Code: "bad_request", Message: err.Error()}})
		return
	}
	if req.DriverID == "" || req.CampaignID == "" {
		writeJSON(w, http.StatusBadRequest, response.Envelope[any]{Error: &response.AppError{Code: "bad_request", Message: "driver_id and campaign_id required"}})
		return
	}
	out, err := h.svc.ClaimCampaign(req.DriverID, req.CampaignID)
	if err != nil {
		writeJSON(w, http.StatusConflict, response.Envelope[any]{Error: &response.AppError{Code: "campaign_claim_failed", Message: err.Error()}})
		return
	}
	writeJSON(w, http.StatusOK, response.Envelope[any]{Data: out})
}

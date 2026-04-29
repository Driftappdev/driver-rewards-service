package service

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"dift_backend_driver/driver-rewards-service/internal/dto"
)

type RewardsService struct {
	mu                sync.Mutex
	redeemedMilestone map[string]map[string]bool
	claimedCampaign   map[string]map[string]bool
}

func NewRewardsService() *RewardsService {
	return &RewardsService{
		redeemedMilestone: map[string]map[string]bool{},
		claimedCampaign:   map[string]map[string]bool{},
	}
}

func (s *RewardsService) GetRewardProgress(driverID string) dto.RewardProgressResponse {
	s.mu.Lock()
	defer s.mu.Unlock()
	return dto.RewardProgressResponse{
		DriverID:         driverID,
		PointsBalance:    2840,
		Tier:             "Gold",
		NextTier:         "Platinum",
		PointsToNextTier: 360,
		Milestones: []dto.RewardMilestone{
			{ID: "reward-1", Title: "Peak-hour completions", Progress: 15, Target: 15, RewardValue: 450, Redeemable: !s.isRedeemed(driverID, "reward-1"), Redeemed: s.isRedeemed(driverID, "reward-1")},
			{ID: "reward-2", Title: "Airport trips", Progress: 4, Target: 6, RewardValue: 600, Redeemable: false, Redeemed: s.isRedeemed(driverID, "reward-2")},
		},
	}
}

func (s *RewardsService) GetCampaigns(driverID string) dto.CampaignsResponse {
	s.mu.Lock()
	defer s.mu.Unlock()
	return dto.CampaignsResponse{
		DriverID: driverID,
		Items: []dto.Campaign{
			{ID: "campaign-1", Title: "Weekend boost", Description: "Complete 12 trips between Friday and Sunday to unlock a THB 800 payout.", Status: campaignStatus(s.isClaimed(driverID, "campaign-1")), ExpiresAt: time.Now().Add(48 * time.Hour).Format(time.RFC3339), RewardValue: 800, Claimable: !s.isClaimed(driverID, "campaign-1"), Claimed: s.isClaimed(driverID, "campaign-1")},
			{ID: "campaign-2", Title: "Morning airport run", Description: "Accept airport pickups before 09:00 to earn an extra THB 120 each.", Status: campaignStatus(s.isClaimed(driverID, "campaign-2")), ExpiresAt: time.Now().Add(7 * 24 * time.Hour).Format(time.RFC3339), RewardValue: 120, Claimable: !s.isClaimed(driverID, "campaign-2"), Claimed: s.isClaimed(driverID, "campaign-2")},
		},
	}
}

func (s *RewardsService) GetRank(driverID string) dto.RankResponse {
	seed := float64(len(strings.TrimSpace(driverID)))
	return dto.RankResponse{
		DriverID:        driverID,
		CurrentRank:     "Gold",
		Score:           91.6 + seed/10,
		Percentile:      87.2,
		TripsToNextRank: 18,
		NextRank:        "Platinum",
	}
}

func (s *RewardsService) RedeemReward(driverID, milestoneID string) (dto.RewardRedemptionResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	progress := s.GetRewardProgressUnlocked(driverID)
	for _, milestone := range progress.Milestones {
		if milestone.ID != milestoneID {
			continue
		}
		if milestone.Redeemed {
			return dto.RewardRedemptionResponse{}, fmt.Errorf("milestone already redeemed")
		}
		if !milestone.Redeemable {
			return dto.RewardRedemptionResponse{}, fmt.Errorf("milestone not redeemable yet")
		}
		s.markRedeemed(driverID, milestoneID)
		return dto.RewardRedemptionResponse{
			DriverID:        driverID,
			RedemptionID:    fmt.Sprintf("redeem-%s-%d", milestoneID, time.Now().Unix()),
			MilestoneID:     milestoneID,
			Status:          "redeemed",
			RewardValue:     milestone.RewardValue,
			RemainingPoints: progress.PointsBalance + int(milestone.RewardValue),
			RedeemedAt:      time.Now().Format(time.RFC3339),
		}, nil
	}
	return dto.RewardRedemptionResponse{}, fmt.Errorf("milestone not found")
}

func (s *RewardsService) ClaimCampaign(driverID, campaignID string) (dto.CampaignClaimResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	campaigns := s.GetCampaignsUnlocked(driverID)
	for _, campaign := range campaigns.Items {
		if campaign.ID != campaignID {
			continue
		}
		if campaign.Claimed {
			return dto.CampaignClaimResponse{}, fmt.Errorf("campaign already claimed")
		}
		if !campaign.Claimable {
			return dto.CampaignClaimResponse{}, fmt.Errorf("campaign not claimable")
		}
		s.markClaimed(driverID, campaignID)
		return dto.CampaignClaimResponse{
			DriverID:    driverID,
			ClaimID:     fmt.Sprintf("claim-%s-%d", campaignID, time.Now().Unix()),
			CampaignID:  campaignID,
			Status:      "claimed",
			RewardValue: campaign.RewardValue,
			ClaimedAt:   time.Now().Format(time.RFC3339),
		}, nil
	}
	return dto.CampaignClaimResponse{}, fmt.Errorf("campaign not found")
}

func (s *RewardsService) GetRewardProgressUnlocked(driverID string) dto.RewardProgressResponse {
	return dto.RewardProgressResponse{
		DriverID:         driverID,
		PointsBalance:    2840,
		Tier:             "Gold",
		NextTier:         "Platinum",
		PointsToNextTier: 360,
		Milestones: []dto.RewardMilestone{
			{ID: "reward-1", Title: "Peak-hour completions", Progress: 15, Target: 15, RewardValue: 450, Redeemable: !s.isRedeemed(driverID, "reward-1"), Redeemed: s.isRedeemed(driverID, "reward-1")},
			{ID: "reward-2", Title: "Airport trips", Progress: 4, Target: 6, RewardValue: 600, Redeemable: false, Redeemed: s.isRedeemed(driverID, "reward-2")},
		},
	}
}

func (s *RewardsService) GetCampaignsUnlocked(driverID string) dto.CampaignsResponse {
	return dto.CampaignsResponse{
		DriverID: driverID,
		Items: []dto.Campaign{
			{ID: "campaign-1", Title: "Weekend boost", Description: "Complete 12 trips between Friday and Sunday to unlock a THB 800 payout.", Status: campaignStatus(s.isClaimed(driverID, "campaign-1")), ExpiresAt: time.Now().Add(48 * time.Hour).Format(time.RFC3339), RewardValue: 800, Claimable: !s.isClaimed(driverID, "campaign-1"), Claimed: s.isClaimed(driverID, "campaign-1")},
			{ID: "campaign-2", Title: "Morning airport run", Description: "Accept airport pickups before 09:00 to earn an extra THB 120 each.", Status: campaignStatus(s.isClaimed(driverID, "campaign-2")), ExpiresAt: time.Now().Add(7 * 24 * time.Hour).Format(time.RFC3339), RewardValue: 120, Claimable: !s.isClaimed(driverID, "campaign-2"), Claimed: s.isClaimed(driverID, "campaign-2")},
		},
	}
}

func (s *RewardsService) isRedeemed(driverID, milestoneID string) bool {
	return s.redeemedMilestone[driverID] != nil && s.redeemedMilestone[driverID][milestoneID]
}

func (s *RewardsService) markRedeemed(driverID, milestoneID string) {
	if s.redeemedMilestone[driverID] == nil {
		s.redeemedMilestone[driverID] = map[string]bool{}
	}
	s.redeemedMilestone[driverID][milestoneID] = true
}

func (s *RewardsService) isClaimed(driverID, campaignID string) bool {
	return s.claimedCampaign[driverID] != nil && s.claimedCampaign[driverID][campaignID]
}

func (s *RewardsService) markClaimed(driverID, campaignID string) {
	if s.claimedCampaign[driverID] == nil {
		s.claimedCampaign[driverID] = map[string]bool{}
	}
	s.claimedCampaign[driverID][campaignID] = true
}

func campaignStatus(claimed bool) string {
	if claimed {
		return "claimed"
	}
	return "active"
}

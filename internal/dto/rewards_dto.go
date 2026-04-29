package dto

type RewardMilestone struct {
	ID          string  `json:"id"`
	Title       string  `json:"title"`
	Progress    float64 `json:"progress"`
	Target      float64 `json:"target"`
	RewardValue float64 `json:"reward_value"`
	Redeemable  bool    `json:"redeemable"`
	Redeemed    bool    `json:"redeemed"`
}

type RewardProgressResponse struct {
	DriverID         string            `json:"driver_id"`
	PointsBalance    int               `json:"points_balance"`
	Tier             string            `json:"tier"`
	NextTier         string            `json:"next_tier"`
	PointsToNextTier int               `json:"points_to_next_tier"`
	Milestones       []RewardMilestone `json:"milestones"`
}

type Campaign struct {
	ID          string  `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Status      string  `json:"status"`
	ExpiresAt   string  `json:"expires_at"`
	RewardValue float64 `json:"reward_value"`
	Claimable   bool    `json:"claimable"`
	Claimed     bool    `json:"claimed"`
}

type CampaignsResponse struct {
	DriverID string     `json:"driver_id"`
	Items    []Campaign `json:"items"`
}

type RankResponse struct {
	DriverID        string  `json:"driver_id"`
	CurrentRank     string  `json:"current_rank"`
	Score           float64 `json:"score"`
	Percentile      float64 `json:"percentile"`
	TripsToNextRank int     `json:"trips_to_next_rank"`
	NextRank        string  `json:"next_rank"`
}

type RewardRedemptionRequest struct {
	DriverID    string `json:"driver_id"`
	MilestoneID string `json:"milestone_id"`
}

type RewardRedemptionResponse struct {
	DriverID        string  `json:"driver_id"`
	RedemptionID    string  `json:"redemption_id"`
	MilestoneID     string  `json:"milestone_id"`
	Status          string  `json:"status"`
	RewardValue     float64 `json:"reward_value"`
	RemainingPoints int     `json:"remaining_points"`
	RedeemedAt      string  `json:"redeemed_at"`
}

type CampaignClaimRequest struct {
	DriverID   string `json:"driver_id"`
	CampaignID string `json:"campaign_id"`
}

type CampaignClaimResponse struct {
	DriverID    string  `json:"driver_id"`
	ClaimID     string  `json:"claim_id"`
	CampaignID  string  `json:"campaign_id"`
	Status      string  `json:"status"`
	RewardValue float64 `json:"reward_value"`
	ClaimedAt   string  `json:"claimed_at"`
}

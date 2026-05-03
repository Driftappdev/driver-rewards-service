package natsadmin

import (
	"context"
	"encoding/json"
	"strings"

	"dift_backend_driver/driver-rewards-service/internal/service"

	"github.com/nats-io/nats.go"
)

type AdminConsumer struct {
	svc *service.RewardsService
}

func NewAdminConsumer(svc *service.RewardsService) *AdminConsumer { return &AdminConsumer{svc: svc} }

func (c *AdminConsumer) Subscribe(ctx context.Context, nc *nats.Conn, subject string) error {
	_, err := nc.Subscribe(subject, func(msg *nats.Msg) {
		_ = c.handle(msg.Data)
	})
	if err != nil {
		return err
	}
	go func() {
		<-ctx.Done()
		_ = nc.Drain()
	}()
	return nil
}

func (c *AdminConsumer) handle(raw []byte) error {
	var cmd struct {
		Action  string         `json:"action"`
		Payload map[string]any `json:"payload"`
	}
	if err := json.Unmarshal(raw, &cmd); err != nil {
		return err
	}
	driverID, _ := cmd.Payload["driver_id"].(string)
	switch strings.TrimSpace(cmd.Action) {
	case "driver.rewards.redeem":
		milestoneID, _ := cmd.Payload["milestone_id"].(string)
		if driverID != "" && milestoneID != "" {
			_, _ = c.svc.RedeemReward(driverID, milestoneID)
		}
	case "driver.campaign.claim":
		campaignID, _ := cmd.Payload["campaign_id"].(string)
		if driverID != "" && campaignID != "" {
			_, _ = c.svc.ClaimCampaign(driverID, campaignID)
		}
	}
	return nil
}

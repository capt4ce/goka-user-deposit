package aboveThresholdProcessor

import (
	"context"
	"fmt"
	"time"

	"github.com/capt4ce/goka-user-deposit/model"
	"github.com/capt4ce/goka-user-deposit/topics"
	"github.com/lovoo/goka"
)

var (
	thresholdAmount   float32 = 10000
	thresholdDuration int64   = 120 // seconds
)

func Start(ctx context.Context, topicDeposits *topics.TopicDeposits) func() error {
	fmt.Println("starting aboveThresholdProcessor...")
	return topicDeposits.GenerateListener(ctx, topics.DepositFlagGroup, []goka.Edge{
		goka.Input(topics.DepositStream, new(topics.DepositCodec), func(ctx goka.Context, msg interface{}) {
			fmt.Println("aboveThresholdProcessor: receiving topic", fmt.Sprintf("%+v", msg))
			var depositFlag *model.DepositFlag

			deposit := msg.(*model.Deposit)
			if v := ctx.Value(); v == nil {
				depositFlag = new(model.DepositFlag)
			} else {
				depositFlag = v.(*model.DepositFlag)
			}

			if depositFlag.Flagged {
				return
			} else if depositFlag.DepositStartTimestamp == 0 || time.Now().Unix()-depositFlag.DepositStartTimestamp > thresholdDuration {
				depositFlag.DepositStartTimestamp = time.Now().Unix()
				depositFlag.BalanceAccumulation = deposit.Amount
			} else if time.Now().Unix()-depositFlag.DepositStartTimestamp <= thresholdDuration {
				depositFlag.BalanceAccumulation += deposit.Amount
				if (depositFlag.BalanceAccumulation) > thresholdAmount {
					depositFlag.Flagged = true
				}
			}

			ctx.SetValue(depositFlag)
		}),
		goka.Persist(new(topics.DepositFlagCodec)),
	})
}

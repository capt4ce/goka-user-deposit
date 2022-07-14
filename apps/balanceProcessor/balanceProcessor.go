package balanceProcessor

import (
	"context"
	"fmt"

	"github.com/capt4ce/goka-user-deposit/model"
	"github.com/capt4ce/goka-user-deposit/topics"
	"github.com/lovoo/goka"
)

func Start(ctx context.Context, brokers []string) func() error {
	return func() error {
		fmt.Println("starting balanceProcessor...")
		topicDeposit := topics.NewTopicDeposits(brokers)
		return topicDeposit.GenerateListener(ctx, topics.BalanceGroup, []goka.Edge{
			goka.Input(topics.DepositStream, new(topics.DepositCodec), func(ctx goka.Context, msg interface{}) {
				deposit := msg.(*model.Deposit)
				depositArray := &model.DepositArray{}
				if v := ctx.Value(); v != nil {
					depositArray = v.(*model.DepositArray)
				}

				depositArray.WalletId = deposit.WalletId
				depositArray.Deposits = append(depositArray.Deposits, deposit)

				ctx.SetValue(depositArray)
			}),
			goka.Persist(new(topics.DepositArrayCodec)),
		})
	}
}

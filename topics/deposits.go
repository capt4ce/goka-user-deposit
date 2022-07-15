package topics

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/capt4ce/goka-user-deposit/model"
	"github.com/lovoo/goka"
	"google.golang.org/protobuf/proto"
)

var (
	DepositStream    goka.Stream = "deposits"
	BalanceGroup     goka.Group  = "balance"
	BalanceTable     goka.Table  = goka.GroupTable(BalanceGroup)
	DepositFlagGroup goka.Group  = "aboveThreshold"
	DepositFlagTable goka.Table  = goka.GroupTable(DepositFlagGroup)
)

type TopicDeposits struct {
	brokers         []string
	depositView     *goka.View
	depositFlagView *goka.View
}

func NewTopicDeposits(brokers []string) *TopicDeposits {
	depositView, err := goka.NewView(brokers, BalanceTable, new(DepositArrayCodec))
	if err != nil {
		panic(err)
	}
	go func() {
		err := depositView.Run(context.Background())
		if err != nil {
			fmt.Println("depositsView err:", err)
		}
	}()

	depositFlagView, err := goka.NewView(brokers, DepositFlagTable, new(DepositFlagCodec))
	if err != nil {
		panic(err)
	}
	go func() {
		err := depositFlagView.Run(context.Background())
		if err != nil {
			fmt.Println("depositFlagView err:", err)
		}
	}()

	return &TopicDeposits{
		brokers:         brokers,
		depositView:     depositView,
		depositFlagView: depositFlagView,
	}
}

func (td *TopicDeposits) Emit(walletId string, deposit *model.Deposit) error {
	emitter, err := goka.NewEmitter(td.brokers, DepositStream, new(DepositCodec))
	if err != nil {
		log.Fatalf("error creating emitter: %v", err)
	}
	defer emitter.Finish()
	err = emitter.EmitSync(walletId, deposit)
	if err != nil {
		fmt.Println("error emitting")
	}
	return err
}

func (td *TopicDeposits) GenerateListener(ctx context.Context, groupName goka.Group, groupFunctions []goka.Edge) func() error {
	return func() error {
		gokaGroup := goka.DefineGroup(groupName,
			groupFunctions...,
		)
		p, err := goka.NewProcessor(td.brokers, gokaGroup)
		if err != nil {
			fmt.Println("Error new processor", err)
			return err
		}

		err = p.Run(ctx)
		if err != nil {
			fmt.Println("Error run processor", err)
			return err
		}

		return err
	}
}

func (td *TopicDeposits) GetDeposits(walletId string) (*model.DepositArray, error) {
	val, err := td.depositView.Get(walletId)
	if err != nil {
		return nil, err
	}
	if val == nil {
		return nil, errors.New("GetDeposits: key not found")
	}

	return val.(*model.DepositArray), nil
}

func (td *TopicDeposits) GetDepositFlag(walletId string) (*model.DepositFlag, error) {
	val, err := td.depositFlagView.Get(walletId)
	if err != nil {
		return nil, err
	}
	if val == nil {
		return nil, errors.New("GetDepositFlag: key not found")
	}

	return val.(*model.DepositFlag), nil
}

type DepositCodec struct{}

func (c *DepositCodec) Encode(value interface{}) ([]byte, error) {
	return proto.Marshal(value.(*model.Deposit))
}

func (c *DepositCodec) Decode(data []byte) (interface{}, error) {
	var m model.Deposit
	return &m, proto.Unmarshal(data, &m)
}

type DepositArrayCodec struct{}

func (c *DepositArrayCodec) Encode(value interface{}) ([]byte, error) {
	return proto.Marshal(value.(*model.DepositArray))
}

func (c *DepositArrayCodec) Decode(data []byte) (interface{}, error) {
	var m model.DepositArray
	return &m, proto.Unmarshal(data, &m)
}

type DepositFlagCodec struct{}

func (c *DepositFlagCodec) Encode(value interface{}) ([]byte, error) {
	return proto.Marshal(value.(*model.DepositFlag))
}

func (c *DepositFlagCodec) Decode(data []byte) (interface{}, error) {
	var m model.DepositFlag
	return &m, proto.Unmarshal(data, &m)
}

package service

import (
	"github.com/capt4ce/goka-user-deposit/model"
	"github.com/capt4ce/goka-user-deposit/topics"
)

type BalanceService struct {
	topicDeposits *topics.TopicDeposits
}

func NewBalanceService(topicDeposits *topics.TopicDeposits) *BalanceService {
	return &BalanceService{
		topicDeposits: topicDeposits,
	}
}

func (bs *BalanceService) ProcessDeposit(walletId string, amount float32) error {
	deposit := &model.Deposit{
		WalletId: walletId,
		Amount:   amount,
	}
	return bs.topicDeposits.Emit(walletId, deposit)
}

func (bs *BalanceService) GetBalance(walletId string) (float32, error) {
	var totalBalance float32
	depositArray, err := bs.topicDeposits.GetDeposits(walletId)
	if err != nil {
		return 0, err
	}
	for _, d := range depositArray.Deposits {
		totalBalance += d.Amount
	}
	return totalBalance, nil
}

func (bs *BalanceService) GetDepositFlag(walletId string) (bool, error) {
	depositFlag, err := bs.topicDeposits.GetDepositFlag(walletId)
	if err != nil {
		return false, err
	}
	return depositFlag.Flagged, nil
}

package service

type BalanceService struct{}

func (*BalanceService) ProcessDeposit(walletId string, amount int32) error {
	return nil
}

func (*BalanceService) GetBalance(walletId string) (int32, error) {
	return 0, nil
}

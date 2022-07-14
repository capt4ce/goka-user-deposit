package controller

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/capt4ce/goka-user-deposit/apps/restApi/service"
	"github.com/capt4ce/goka-user-deposit/apps/restApi/utils"
	"github.com/gorilla/mux"
)

type BalanceController struct {
	balanceService *service.BalanceService
}

func NewBalanceController(balanceService *service.BalanceService) *BalanceController {
	return &BalanceController{
		balanceService: balanceService,
	}
}

type depositRequest struct {
	WalletId string  `json:"wallet_id"`
	Amount   float32 `json:"amount"`
}

func (bc *BalanceController) Deposit(w http.ResponseWriter, r *http.Request) {
	var req depositRequest

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	err = json.Unmarshal(b, &req)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	if !(req.Amount > 0) {
		utils.RespondWithError(w, http.StatusUnprocessableEntity, "amount must be more than 0")
		return
	}

	err = bc.balanceService.ProcessDeposit(req.WalletId, req.Amount)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, nil)
}

type getBalanceResponse struct {
	WalletId       string  `json:"wallet_id"`
	Balance        float32 `json:"balance"`
	AboveThreshold bool    `json:"above_threshold"`
}

func (bc *BalanceController) GetBalance(w http.ResponseWriter, r *http.Request) {
	walletId := mux.Vars(r)["wallet_id"]

	response := getBalanceResponse{
		WalletId:       walletId,
		Balance:        0,
		AboveThreshold: false,
	}
	balance, err := bc.balanceService.GetBalance(walletId)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusOK, response)
		return
	}
	response.Balance = balance

	flag, err := bc.balanceService.GetDepositFlag(walletId)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusOK, response)
		return
	}
	response.AboveThreshold = flag

	utils.RespondWithJSON(w, http.StatusOK, response)
}

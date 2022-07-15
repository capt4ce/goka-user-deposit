package restApi

import (
	"fmt"
	"log"
	"net/http"

	"github.com/capt4ce/goka-user-deposit/apps/restApi/controller"
	"github.com/capt4ce/goka-user-deposit/apps/restApi/service"
	"github.com/capt4ce/goka-user-deposit/topics"

	"github.com/gorilla/mux"
)

func Start(port string, topicDeposits *topics.TopicDeposits) {
	balanceService := service.NewBalanceService(topicDeposits)
	balanceController := controller.NewBalanceController(balanceService)

	router := mux.NewRouter()
	router.HandleFunc("/deposit", balanceController.Deposit).Methods("POST")
	router.HandleFunc("/balance/{wallet_id}", balanceController.GetBalance).Methods("GET")

	log.Printf("Listen port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), router))

}

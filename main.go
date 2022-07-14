package main

import (
	"github.com/capt4ce/goka-user-deposit/topics"
)

var topicDepositsInstance topics.TopicDeposits

func init() {
	// initializing topics
	topicDepositsInstance := topics.TopicDeposits{}
	go topicDepositsInstance.RunEmitter() // emits one message and stops
	topicDepositsInstance.RunProcessor()  // press ctrl-c to stop

	// initialize Api
}

func main() {

}

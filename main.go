package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/capt4ce/goka-user-deposit/apps/aboveThresholdProcessor"
	"github.com/capt4ce/goka-user-deposit/apps/balanceProcessor"
	"github.com/capt4ce/goka-user-deposit/apps/restApi"
	"github.com/capt4ce/goka-user-deposit/topics"
	"golang.org/x/sync/errgroup"
)

func main() {
	brokers := []string{"localhost:9092"}
	topicDeposits := topics.NewTopicDeposits(brokers)

	ctx, cancel := context.WithCancel(context.Background())

	grp, ctx := errgroup.WithContext(ctx)
	grp.Go(balanceProcessor.Start(ctx, topicDeposits))
	grp.Go(aboveThresholdProcessor.Start(ctx, topicDeposits))
	go restApi.Start("8000", topicDeposits)

	waiter := make(chan os.Signal, 1)
	signal.Notify(waiter, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-waiter:
	case <-ctx.Done():
	}
	cancel()
	if err := grp.Wait(); err != nil {
		log.Println(err)
	}
}

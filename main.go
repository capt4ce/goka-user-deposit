package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/capt4ce/goka-user-deposit/apps/aboveTresholdProcessor"
	"github.com/capt4ce/goka-user-deposit/apps/balanceProcessor"
	"github.com/capt4ce/goka-user-deposit/apps/restApi"
	"golang.org/x/sync/errgroup"
)

func main() {
	brokers := []string{"localhost:9092"}
	ctx, cancel := context.WithCancel(context.Background())

	grp, ctx := errgroup.WithContext(ctx)
	grp.Go(aboveTresholdProcessor.Start(ctx, brokers))
	grp.Go(balanceProcessor.Start(ctx, brokers))
	grp.Go(restApi.Start("8080", brokers))

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

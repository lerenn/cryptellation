package main

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/lerenn/cryptellation/v1/api"
	"github.com/spf13/cobra"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

var (
	ticksListenExchangeFlag string
	ticksListenPairFlag     string
)

var ticksCmd = &cobra.Command{
	Use:     "ticks",
	Aliases: []string{"t"},
	Short:   "Manage ticks",
}

var ticksListenCmd = &cobra.Command{
	Use:     "listen",
	Aliases: []string{"l"},
	Short:   "Listen to ticks",
	RunE: func(cmd *cobra.Command, _ []string) error {
		// Create temporary worker
		tq := fmt.Sprintf("cryptellation-cli-ticks-listen-%s", uuid.New().String())
		w := worker.New(cryptellationClient.Temporal(), tq, worker.Options{})
		w.RegisterWorkflowWithOptions(ticksListenCallbackWorkflow, workflow.RegisterOptions{
			Name: tq,
		})

		// Start worker
		irq := worker.InterruptCh()
		go w.Run(irq)

		// Listen to ticks
		_, err := cryptellationClient.ListenToTicks(cmd.Context(),
			api.RegisterForTicksListeningWorkflowParams{
				Exchange: ticksListenExchangeFlag,
				Pair:     ticksListenPairFlag,
				CallbackWorkflow: api.ListenToTicksCallbackWorkflow{
					Name:          tq,
					TaskQueueName: tq,
				},
			})
		if err != nil {
			return err
		}

		// Wait for interrupt
		<-irq

		// Stop worker
		w.Stop()

		return nil
	},
}

func ticksListenCallbackWorkflow(ctx workflow.Context, params api.ListenToTicksCallbackWorkflowParams) error {
	fmt.Println(params.Tick.String())
	return nil
}

func addTicksCommands() {
	ticksListenCmd.Flags().StringVarP(&ticksListenExchangeFlag, "exchange", "e", "binance", "Exchange")
	ticksListenCmd.Flags().StringVarP(&ticksListenPairFlag, "pair", "p", "BTC-USDT", "Pair")
	ticksCmd.AddCommand(ticksListenCmd)
	rootCmd.AddCommand(ticksCmd)
}

package main

import (
	"context"
	"fmt"
	"github.com/ledgerwatch/turbo-geth/cmd/rpctest/rpctest"
	"github.com/ledgerwatch/turbo-geth/log"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	log.SetupDefaultTerminalLogger(log.Lvl(3), "", "")

	var (
		needCompare bool
		fullTest    bool
		gethURL     string
		tgURL       string
		blockNum    uint64
		chaindata   string
	)
	withTGUrl := func(cmd *cobra.Command) {
		cmd.Flags().StringVar(&tgURL, "tgUrl", "http://localhost:8545", "turbogeth rpcdaemon url")
	}
	withGethUrl := func(cmd *cobra.Command) {
		cmd.Flags().StringVar(&gethURL, "gethUrl", "http://localhost:8546", "geth rpc url")
	}
	withBlockNum := func(cmd *cobra.Command) {
		cmd.Flags().Uint64Var(&blockNum, "block", 2000000, "Block number")
	}
	withNeedCompare := func(cmd *cobra.Command) {
		cmd.Flags().BoolVar(&needCompare, "needCompare", false, "need compare with geth")
	}
	with := func(cmd *cobra.Command, opts ...func(*cobra.Command)) {
		for i := range opts {
			opts[i](cmd)
		}
	}

	var bench1Cmd = &cobra.Command{
		Use:   "bench1",
		Short: "",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			rpctest.Bench1(tgURL, gethURL, needCompare, fullTest)
		},
	}
	with(bench1Cmd, withTGUrl, withGethUrl, withNeedCompare)
	bench1Cmd.Flags().BoolVar(&fullTest, "fullTest", false, "some text")

	var bench2Cmd = &cobra.Command{
		Use:   "bench2",
		Short: "",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			rpctest.Bench2(tgURL)
		},
	}
	var bench3Cmd = &cobra.Command{
		Use:   "bench3",
		Short: "",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			rpctest.Bench3(tgURL, gethURL)
		},
	}
	with(bench3Cmd, withTGUrl, withGethUrl)

	var bench4Cmd = &cobra.Command{
		Use:   "bench4",
		Short: "",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			rpctest.Bench4(tgURL)
		},
	}
	with(bench4Cmd, withTGUrl)

	var bench5Cmd = &cobra.Command{
		Use:   "bench5",
		Short: "",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			rpctest.Bench5(tgURL)
		},
	}
	with(bench5Cmd, withTGUrl)
	var bench6Cmd = &cobra.Command{
		Use:   "bench6",
		Short: "",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			rpctest.Bench6(tgURL)
		},
	}
	with(bench6Cmd, withTGUrl)

	var bench7Cmd = &cobra.Command{
		Use:   "bench7",
		Short: "",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			rpctest.Bench7(tgURL, gethURL)
		},
	}
	with(bench7Cmd, withTGUrl, withGethUrl)

	var bench8Cmd = &cobra.Command{
		Use:   "bench6",
		Short: "",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			rpctest.Bench8(tgURL)
		},
	}
	with(bench8Cmd, withTGUrl)

	var bench9Cmd = &cobra.Command{
		Use:   "bench9",
		Short: "",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			rpctest.Bench9(tgURL, gethURL, needCompare)
		},
	}
	with(bench9Cmd, withTGUrl, withGethUrl, withNeedCompare)

	var bench10Cmd = &cobra.Command{
		Use:   "bench6",
		Short: "",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			err := rpctest.Bench10(tgURL, gethURL, blockNum)
			if err != nil {
				log.Error("bench 10 err", "err", err)
			}
		},
	}
	with(bench10Cmd, withGethUrl, withTGUrl, withBlockNum)

	var proofsCmd = &cobra.Command{
		Use:   "proofs",
		Short: "",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			rpctest.Proofs(chaindata, gethURL, blockNum)
		},
	}
	proofsCmd.Flags().StringVar(&chaindata, "chaindata", "", "")

	var fixStateCmd = &cobra.Command{
		Use:   "fixstate",
		Short: "",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			rpctest.FixState(chaindata, gethURL)
		},
	}
	fixStateCmd.Flags().StringVar(&chaindata, "chaindata", "", "")

	var tmpDataDir, tmpDataDirOrig string
	var notRegenerateGethData bool
	var compareAccountRange = &cobra.Command{
		Use:   "compareAccountRange",
		Short: "",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			rpctest.CompareAccountRange(tgURL, gethURL, tmpDataDir, tmpDataDirOrig, blockNum, notRegenerateGethData)
		},
	}
	with(compareAccountRange, withTGUrl, withGethUrl, withBlockNum)
	compareAccountRange.Flags().BoolVar(&notRegenerateGethData, "regenGethData", false, "")
	compareAccountRange.Flags().StringVar(&tmpDataDir, "tmpdir", "/media/b00ris/nvme/accrange1", "dir for tmp db")
	compareAccountRange.Flags().StringVar(&tmpDataDirOrig, "gethtmpdir", "/media/b00ris/nvme/accrangeorig1", "dir for tmp db")

	var rootCmd = &cobra.Command{Use: "test"}
	rootCmd.Flags().StringVar(&tgURL, "tgUrl", "http://localhost:8545", "turbogeth rpcdaemon url")
	rootCmd.Flags().StringVar(&gethURL, "gethUrl", "http://localhost:8546", "geth rpc url")
	rootCmd.Flags().Uint64Var(&blockNum, "block", 2000000, "Block number")

	rootCmd.AddCommand(
		bench1Cmd,
		bench2Cmd,
		bench3Cmd,
		bench4Cmd,
		bench5Cmd,
		bench6Cmd,
		bench7Cmd,
		bench8Cmd,
		bench9Cmd,
		bench10Cmd,
		proofsCmd,
		fixStateCmd,
		compareAccountRange,
	)
	if err := rootCmd.ExecuteContext(rootContext()); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func rootContext() context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
		defer signal.Stop(ch)

		select {
		case <-ch:
			log.Info("Got interrupt, shutting down...")
		case <-ctx.Done():
		}

		cancel()
	}()
	return ctx
}

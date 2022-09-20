package runner

import "github.com/urfave/cli"

func (runner *EthereumTransactionScannerRunner) getStartFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:        "ethereum-host",
			Value:       "http://localhost:8545",
			Usage:       "host endpoint of the ethereum node",
			Destination: &runner.endpoint,
		},
		cli.StringFlag{
			Name:        "send-ethereum-host",
			Value:       "http://localhost:8545",
			Usage:       "host endpoint of the ethereum node",
			Destination: &runner.sendEndpoint,
		},
		cli.IntFlag{
			Name:        "block-workers",
			Value:       1,
			Usage:       "number of routines that will pull block information",
			Destination: &runner.blockWorkerNum,
		},
		cli.Int64Flag{
			Name:        "start-block",
			Value:       1,
			Usage:       "block number of the first block to begin scanning",
			Destination: &runner.startBlock,
		},
		cli.Int64Flag{
			Name:        "end-block",
			Value:       6000000,
			Usage:       "block number of the first block to begin scanning",
			Destination: &runner.endBlock,
		},
		cli.StringFlag{
			Name:        "filter-address",
			Usage:       "filters all transactions to those only containing the specified address",
			Destination: &runner.filterAddress,
		},
	}
}

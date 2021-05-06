package runner

import (
	"fmt"
	"os/signal"
	"syscall"
)

func (runner *EthereumTransactionScannerRunner) handleShutdown() {
	defer runner.waitGroup.Done()

	signal.Notify(runner.sigKillChan, syscall.SIGINT, syscall.SIGTERM)

	<-runner.sigKillChan
	signal.Reset()

	runner.blockWorkerManager.Stop()

	runner.transactionWorker.Stop()

	runner.transactionReporter.Stop()

	close(runner.rawTransactions)
	close(runner.filteredTransactions)

	if err := runner.outputfile.SaveAs("Transaction.xlsx"); err != nil {
		fmt.Println(err)
	}

	runner.done = true
}

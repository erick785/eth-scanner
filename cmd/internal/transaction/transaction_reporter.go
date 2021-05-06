package transaction

import (
	"encoding/json"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
)

var categories = map[string]string{
	"hash": "A", "nonce": "B", "blockHash": "C",
	"blockNumber": "D", "transactionIndex": "E", "from": "F",
	"to": "G", "value": "H", "gas": "I",
	"gasPrice": "J", "input": "K", "raw": "L"}

type TransactionReporter struct {
	outputFile           *excelize.File
	filteredTransactions chan *TransactionResult
	waitGroup            *sync.WaitGroup
	txCount              int
	done                 bool
}

func NewTransactionReporter(outputFile *excelize.File, transactions chan *TransactionResult, wg *sync.WaitGroup) *TransactionReporter {

	return &TransactionReporter{
		txCount:              2,
		outputFile:           outputFile,
		filteredTransactions: transactions,
		waitGroup:            wg,
	}
}

func (reporter *TransactionReporter) Start() {
	defer reporter.waitGroup.Done()

	for !reporter.done {
		select {
		case transaction := <-reporter.filteredTransactions:

			t, _ := json.Marshal(transaction)

			result := make(map[string]interface{})
			json.Unmarshal(t, &result)

			for k, v := range result {
				if err := reporter.outputFile.SetCellValue("Sheet1", categories[k]+strconv.Itoa(reporter.txCount), v); err != nil {
					log.Println("Error writing transaction to file:", err.Error())
					reporter.done = true
				}
			}

			reporter.txCount++
		default:
			time.Sleep(time.Millisecond * 100)
		}
	}
}

func (reporter *TransactionReporter) Stop() {
	reporter.done = true
}

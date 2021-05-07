package transaction

import (
	"encoding/json"
	"log"
	"math/big"
	"strconv"
	"sync"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

var categories = map[string]string{
	"hash": "A", "nonce": "B", "blockHash": "C",
	"blockNumber": "D", "transactionIndex": "E", "from": "F",
	"to": "G", "value": "H", "gas": "I",
	"gasPrice": "J", "input": "K", "raw": "L", "isContract": "M"}

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

				if k == "value" {
					vBig := hexutil.MustDecodeBig(v.(string))
					v = FromWei(vBig, 18).String()
				}

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

// FromWei transforms wei to a decimal value with `decimals` orders of magnitude (divide by 10^`decimals`)
func FromWei(wei *big.Int, decimals int) *big.Float {
	floatValue := new(big.Float).SetInt(wei)
	mul := big.NewInt(10)
	mul = mul.Exp(mul, big.NewInt(int64(decimals)), nil)
	floatMul := new(big.Float).SetInt(mul)
	result := new(big.Float).Quo(floatValue, floatMul)
	return result
}

// ToWei transforms decimal value to wei with `decimals` orders of magnitude (multiply by 10^`decimals`)
func ToWei(value *big.Float, decimals int) *big.Int {
	floatValue := new(big.Float).Copy(value)
	mul := big.NewInt(10)
	mul = mul.Exp(mul, big.NewInt(int64(decimals)), nil)
	floatMul := new(big.Float).SetInt(mul)
	floatValue = floatValue.Mul(floatValue, floatMul)
	wei, _ := floatValue.Int(new(big.Int))
	return wei
}

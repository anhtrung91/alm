package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
	// "github.com/example_cc/common"

	"github.com/example_cc/common"
	"github.com/example_cc/models"
	"github.com/example_cc/util"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type Transaction models.Transaction

func row_keys_of_transaction(transaction *Transaction) []string {
	return []string{strconv.Itoa(transaction.TransactionId)}
}

func (transaction *Transaction) CreateTransaction(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	TransactionId, err := strconv.Atoi(args[0])
	UserMerchantId, err2 := strconv.Atoi(args[1])
	Amount := args[2]
	ReceiveValue := args[3]
	ReceiveMerchantId := args[4]
	ReceiveProductId := args[5]
	ReceiveRedeemId := args[6]
	Blockid := args[7]
	TransactionDatetime := time.Now().Format("2006.01.02 15:04:05")
	MerchantId, err := strconv.Atoi(args[8])

	err = create_transaction_(stub, &Transaction{TransactionId: TransactionId, UserMerchantId: UserMerchantId,
		Amount: Amount, ReceiveValue: ReceiveValue, ReceiveMerchantId: ReceiveMerchantId,
		ReceiveProductId: ReceiveProductId, ReceiveRedeemId: ReceiveRedeemId, Blockid: Blockid,
		TransactionDatetime: TransactionDatetime, MerchantId: MerchantId})

	if err != nil || err2 != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func create_transaction_(stub shim.ChaincodeStubInterface, transaction *Transaction) error {
	var old_transaction Transaction
	row_was_found, err := util.InsertTableRow(stub, models.TRANSACTION_TABLE, row_keys_of_transaction(transaction), transaction, util.FAIL_BEFORE_OVERWRITE, &old_transaction)
	if err != nil {
		return err
	}
	if row_was_found {
		return fmt.Errorf("Could not create transaction %v because an transaction already exists", *transaction)
	}
	return nil // success
}

func (transaction *Transaction) GetAllTransaction(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	transactionList := []*Transaction{}
	transactionList, err := get_all_transaction_(stub, args)
	transactionJson, err2 := json.Marshal(transactionList)
	if err != nil || err2 != nil {
		return shim.Error("convert json error")
	}

	return shim.Success(transactionJson)
}

func get_all_transaction_(stub shim.ChaincodeStubInterface, Name []string) ([]*Transaction, error) {
	row_json_bytes_channel, err := util.GetTableRows(stub, models.TRANSACTION_TABLE, []string{}) // empty row_keys to get all entries
	if err != nil {
		return nil, fmt.Errorf("Could not get %v", err.Error())
	}

	transaction := new(Transaction)
	transactionList := []*Transaction{}

	for row_json_bytes := range row_json_bytes_channel {
		transaction = new(Transaction)
		err = json.Unmarshal(row_json_bytes, transaction)
		if err != nil {
			return nil, fmt.Errorf("Could not get; json.Unmarshal of \"%s\" failed with error %v", string(row_json_bytes), err)
		}
		transactionList = append(transactionList, transaction)
	}

	return transactionList, nil
}

func (transaction *Transaction) GetExchanged(stub shim.ChaincodeStubInterface, arg []string) pb.Response {
	transactionList := []*Transaction{}

	for i := range arg {
		TransactionId := arg[i]
		transaction, err := get_transaction_(stub, TransactionId)
		if err != nil {
			return shim.Error("get Client Asset error")
		}

		//fmt.Printf("query_amount Response: %s\n", string(transaction.Amount));
		transactionList = append(transactionList, transaction)
	}

	amountTotal := 0
	for j := range transactionList {
		amount, err := strconv.Atoi(transactionList[j].Amount)
		if err != nil {
			return shim.Error("casting error")
		}
		amountTotal += amount
	}

	transactionJson, err2 := json.Marshal(transactionList)
	if err2 != nil {
		return shim.Error("convert json error")
	}

	fmt.Println("amountTotal: ", amountTotal)
	//resSuc := common.ResponseSuccess{common.SUCCESS, common.ResCodeDict[common.SUCCESS], strconv.Itoa(amountTotal)}
	//return common.RespondSuccess(&resSuc)
	return shim.Success(transactionJson)
}

func get_transaction_(stub shim.ChaincodeStubInterface, TransactionId string) (*Transaction, error) {
	var transaction Transaction
	row_was_found, err := util.GetTableRow(stub, models.TRANSACTION_TABLE, []string{TransactionId}, &transaction, util.FAIL_IF_MISSING)

	if err != nil || !row_was_found {
		return nil, fmt.Errorf("Could not retrieve Client Asset Client ID")
	}
	// if !row_was_found {
	// 	return nil, fmt.Errorf("Client Asset of Client ID does not exist")
	// }

	return &transaction, nil
}

func (transaction *Transaction) GetTransactionByMerchant(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//MerchantId := args[0]
	MerchantId, err := strconv.Atoi(args[0])

	queryString := fmt.Sprintf("{\"selector\":{\"MerchantId\":%d}}", MerchantId)
	result, err := get_detail_transaction_(stub, queryString)
	if err != nil {
		return shim.Error("casting error")
	}

	if string(result[:]) != "[]" {
		resSuc := common.ResponseSuccess{common.SUCCESS, common.ResCodeDict[common.SUCCESS], string(result[:])}
		return common.RespondSuccess(&resSuc)
	} else {
		resErr := common.ResponseError{common.ERR1, common.ResCodeDict[common.ERR1]}
		return common.RespondError(&resErr)
	}
}

func get_detail_transaction_(stub shim.ChaincodeStubInterface, queryString string) ([]byte, error) {
	fmt.Printf("- getQueryResultForQueryString queryString:\n%s\n", queryString)
	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryRecords
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- getQueryResultForQueryString queryResult:\n%s\n", buffer.String())
	return buffer.Bytes(), nil
}

func (transaction *Transaction) GetTopExchanged(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	sortingType := args[0]
	transactionList := []*Transaction{}
	sortingList := []*Transaction{}
	transactionList, err := get_all_transaction_(stub, args)
	//nowDate := time.Now().Format("2006.01.02 15:04:05")
	nowYear, nowMonth, nowDay := time.Now().Date()

	for i := range transactionList {
		txtDatetime, err := time.Parse("2006.01.02 15:04:05", transactionList[i].TransactionDatetime)
		if err != nil {
			fmt.Println(err)
		}
		txtYear, txtMonth, txtDay := txtDatetime.Date()

		switch sortingType {
		case "day":
			if txtDay == nowDay {
				sortingList = append(sortingList, transactionList[i])
				break
			}
		case "month":
			if txtMonth == nowMonth {
				sortingList = append(sortingList, transactionList[i])
				break
			}
		case "year":
			if txtYear == nowYear {
				sortingList = append(sortingList, transactionList[i])
				break
			}
		default:
			sortingList = append(sortingList, transactionList[i])
		}
	}

	duplicate_frequency := make(map[int]int)
	for _, item := range sortingList {
		_, exist := duplicate_frequency[item.MerchantId]
		txtAmount, err := strconv.Atoi(item.Amount)
		if err != nil {
			return shim.Error("casting error")
		}
		if exist {

			duplicate_frequency[item.MerchantId] += txtAmount
		} else {
			duplicate_frequency[item.MerchantId] = txtAmount
		}
	}
	for k, v := range duplicate_frequency {
		fmt.Printf("Item : %d , Count : %d\n", k, v)
	}

	jsonString, err := json.Marshal(duplicate_frequency)
	if err != nil {
		return shim.Error("casting error")
	}
	fmt.Println("jsonString: ", string(jsonString))
	resSuc := common.ResponseSuccess{common.SUCCESS, common.ResCodeDict[common.SUCCESS], string(jsonString)}
	return common.RespondSuccess(&resSuc)
	// transactionJson, err2 := json.Marshal(sortingList)
	// if err != nil || err2 != nil {
	// 	return shim.Error("convert json error")
	// }

	// return shim.Success(transactionJson)
}

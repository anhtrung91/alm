package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/example_cc/common"
	"github.com/example_cc/models"
	"github.com/example_cc/util"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type Account models.Account
type ResponseSuccess common.ResponseSuccess
type ResponseError common.ResponseError

// func GetTotalUsers (totalUser, totalToken string) (string, string) {
// 	fmt.Println("GetTotalUsers: ", totalUser)
// 	fmt.Println("GetTotalToken: ", totalToken)
// 	return totalUser, totalToken
// }

func (account *Account)GetTotalToken(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	accountList := []*Account{}
	accountList, err := get_User_(stub, args)

	//accountJson, err2 := json.Marshal(accountList)
	if err != nil {
		return shim.Error("convert json error")
	}

	totalToken := 0
	for i := range accountList {
		tokenAmount, err := strconv.Atoi(accountList[i].TokenAmount)
		if err != nil {
			return shim.Error("casting error")
		}
		totalToken += tokenAmount
	}

	//return shim.Success(accountJson)
	resSuc := common.ResponseSuccess{common.SUCCESS, common.ResCodeDict[common.SUCCESS], ("[{\"totalToken\": \"" + strconv.Itoa(totalToken) + "\"}]")}
	return common.RespondSuccess(&resSuc)
}

//"data":"[{"MerchantId":1,"MerchantName":"CGV","SumCoin":"10","RateCash":"10"}]
//"data":"totalToken: "100""

func (account *Account) GetAllUsers(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	accountList := []*Account{}
	accountList, err := get_User_(stub, args)

	accountJson, err2 := json.Marshal(accountList)
	if err != nil || err2 != nil {
		return shim.Error("convert json error")
	}

	// totalUser := 0
	// totalToken := 0
	// for i := range accountList {
	// 	totalUser += i
	// 	fmt.Println("TokenAmount: ", accountList[i].TokenAmount)
	// 	tokenAmount, err := strconv.Atoi(accountList[i].TokenAmount)
	// 	if err != nil {
	// 		return shim.Error("casting error")
	// 	}
	// 	totalToken += tokenAmount
	// }

	//GetTotalUsers(strconv.Itoa(totalUser + 1), strconv.Itoa(totalToken))

	return shim.Success(accountJson)
}

func get_User_(stub shim.ChaincodeStubInterface, Name []string) ([]*Account, error) {
	row_json_bytes_channel, err := util.GetTableRows(stub, models.ACCOUNT_TABLE, []string{}) // empty row_keys to get all entries

	if err != nil {
		return nil, fmt.Errorf("Could not get account names; %v", err.Error())
	}

	account := new(Account)
	accountList := []*Account{}

	for row_json_bytes := range row_json_bytes_channel {
		account = new(Account)
		err = json.Unmarshal(row_json_bytes, account)
		if err != nil {
			return nil, fmt.Errorf("Could not get account names; json.Unmarshal of \"%s\" failed with error %v", string(row_json_bytes), err)
		}
		accountList = append(accountList, account)
	}

	return accountList, nil
}

func (account *Account) CreateAccount(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments.  Expecting 5")
	}
	ID, err := strconv.Atoi(args[0])
	UserName := args[1]
	Password := args[2]
	TokenAmount := "50"
	Name := args[3]
	Email := args[4]
	if err != nil {
		return shim.Error(fmt.Sprintf("Malformed initial_balance string \"%s\"; expecting nonnegative integer", args[1]))
	}

	err = create_account_(stub, &Account{ID: ID, UserName: UserName, Password: Password, TokenAmount: TokenAmount, Name: Name, Email: Email})
	if err != nil {
		return shim.Error(err.Error())
	}

	resSuc := common.ResponseSuccess{common.SUCCESS, common.ResCodeDict[common.SUCCESS], ""}
	return common.RespondSuccess(&resSuc)
}

func (account *Account) Login(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments.  Expecting 2")
	}
	UserName := args[0]
	password := args[1]
	queryString := fmt.Sprintf("{\"selector\":{\"UserName\":\"%s\",\"Password\":\"%s\"}}", UserName, password)

	queryResults, err := getQueryResultForQueryString(stub, queryString)

	if err != nil {
		return shim.Error(err.Error())
	}

	if string(queryResults[:]) != "[]" {
		resSuc := common.ResponseSuccess{common.SUCCESS, common.ResCodeDict[common.SUCCESS], string(queryResults[:])}
		return common.RespondSuccess(&resSuc)
	} else {
		resErr := common.ResponseError{common.ERR1, common.ResCodeDict[common.ERR1]}
		return common.RespondError(&resErr)
	}
}

func row_keys_of_Account(account *Account) []string {
	return []string{account.UserName}
}

func getQueryResultForQueryString(stub shim.ChaincodeStubInterface, queryString string) ([]byte, error) {

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
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- getQueryResultForQueryString queryResult:\n%s\n", buffer.String())
	return buffer.Bytes(), nil
}

func create_account_(stub shim.ChaincodeStubInterface, account *Account) error {
	var old_account Account
	row_was_found, err := util.InsertTableRow(stub, models.ACCOUNT_TABLE, row_keys_of_Account(account), account, util.FAIL_BEFORE_OVERWRITE, &old_account)
	if err != nil {
		return err
	}
	if row_was_found {
		return fmt.Errorf("Could not create account %v because an account with that Name already exists", *account)
	}
	return nil // success
}

package controllers

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/example_cc/common"
	"github.com/example_cc/models"
	"github.com/example_cc/util"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type Merchant models.Merchant

func row_keys_of_merchant(merchant *Merchant) []string {
	return []string{strconv.Itoa(merchant.MerchantId)}
}

func (merchant *Merchant) CreateMerchant(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//logger.Info("########### create_merchant Init ###########")

	// if len(args) != 3 {
	// 	return shim.Error("Incorrect number of arguments.  Expecting 3")
	// }

	MerchantId, err := strconv.Atoi(args[0])
	MerchantName := args[1]
	TokenAmount := args[2]
	RemainCoin := args[3]
	SumCoin := args[4]
	RateCash := args[5]

	// if err != nil {
	// 	return shim.Error(fmt.Sprintf("Malformed initial_balance string \"%s\"; expecting nonnegative integer", args[1]))
	// }

	err = create_merchant_(stub, &Merchant{MerchantId: MerchantId, MerchantName: MerchantName,
		TokenAmount: TokenAmount, RemainCoin: RemainCoin, SumCoin: SumCoin, RateCash: RateCash})

	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func create_merchant_(stub shim.ChaincodeStubInterface, merchant *Merchant) error {
	var old_merchant Merchant
	row_was_found, err := util.InsertTableRow(stub, models.MERCHANT_TABLE, row_keys_of_merchant(merchant), merchant, util.FAIL_BEFORE_OVERWRITE, &old_merchant)
	if err != nil {
		return err
	}
	if row_was_found {
		return fmt.Errorf("Could not create merchant %v because an merchant with that Name already exists", *merchant)
	}

	return nil // success
}

func (merchant *Merchant) GetMerchant(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	MerchantName := args[0]
	merchant, err := get_Merchant_(stub, MerchantName)
	if err != nil {
		return shim.Error("get Client Asset error")
	}

	fmt.Printf("query_merchant Response: %s\n", string(merchant.MerchantName))
	// Serialize ClientAsset struct as JSON
	bytes, err := json.Marshal(merchant)
	if err != nil {
		return shim.Error("convert json error")
	}

	fmt.Printf("query_merchant Response: %s\n", string(bytes))
	return shim.Success(bytes)
}

func get_Merchant_(stub shim.ChaincodeStubInterface, MerchantName string) (*Merchant, error) {
	var merchant Merchant
	row_was_found, err := util.GetTableRow(stub, models.MERCHANT_TABLE, []string{MerchantName}, &merchant, util.FAIL_IF_MISSING)

	if err != nil {
		return nil, fmt.Errorf("Could not retrieve Client Asset Client ID \"%s\"; error was %v", MerchantName, err.Error())
	}
	if !row_was_found {
		return nil, fmt.Errorf("Client Asset of Client ID  \"%s\" does not exist", MerchantName)
	}

	fmt.Printf("query_merchant Response: %s\n", string(merchant.MerchantName))
	return &merchant, nil
}

// func GetTotalMerchants (totalMerchant string) string {
// 	fmt.Println("GetTotalMerchants: ", totalMerchant)
// 	return totalMerchant
// }


func (merchant *Merchant) GetAllMerchant(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	merchantList := []*Merchant{}
	merchantList, err := get_all_merchant_(stub, args)
	merchantJson, err2 := json.Marshal(merchantList)
	if err != nil || err2 != nil {
		return shim.Error("convert json error")
	}

	// totalMerchant := 0
	// for i := range merchantList {
	// 	totalMerchant += i
	// }
	// GetTotalMerchants(strconv.Itoa(totalMerchant + 1))

	//return shim.Success(merchantJson)
	resSuc := common.ResponseSuccess{common.SUCCESS, common.ResCodeDict[common.SUCCESS], string(merchantJson)}
	return common.RespondSuccess(&resSuc)
}

func get_all_merchant_(stub shim.ChaincodeStubInterface, Name []string) ([]*Merchant, error) {
	row_json_bytes_channel, err := util.GetTableRows(stub, models.MERCHANT_TABLE, []string{}) // empty row_keys to get all entries
	if err != nil {
		return nil, fmt.Errorf("Could not get %v", err.Error())
	}

	merchant := new(Merchant)
	merchantList := []*Merchant{}

	for row_json_bytes := range row_json_bytes_channel {
		merchant = new(Merchant)
		err = json.Unmarshal(row_json_bytes, merchant)
		if err != nil {
			return nil, fmt.Errorf("Could not get; json.Unmarshal of \"%s\" failed with error %v", string(row_json_bytes), err)
		}

		merchantList = append(merchantList, merchant)
	}

	// totalMerchant := 0
	// for i := range merchantList {
	// 	totalMerchant += i
	// }
	// fmt.Println("totalMerchant: ", totalMerchant)

	return merchantList, nil
}

// func getQueryResultForQueryString2(stub shim.ChaincodeStubInterface, queryString string) ([]byte, error) {

// 	fmt.Printf("- getQueryResultForQueryString queryString:\n%s\n", queryString)

// 	resultsIterator, err := stub.GetQueryResult(queryString)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resultsIterator.Close()

// 	// buffer is a JSON array containing QueryRecords
// 	var buffer bytes.Buffer
// 	buffer.WriteString("[")

// 	bArrayMemberAlreadyWritten := false
// 	for resultsIterator.HasNext() {
// 		queryResponse, err := resultsIterator.Next()
// 		if err != nil {
// 			return nil, err
// 		}
// 		// Add a comma before array members, suppress it for the first array member
// 		if bArrayMemberAlreadyWritten == true {
// 			buffer.WriteString(",")
// 		}
// 		buffer.WriteString("{\"Key\":")
// 		buffer.WriteString("\"")
// 		buffer.WriteString(queryResponse.Key)
// 		buffer.WriteString("\"")

// 		buffer.WriteString(", \"Record\":")
// 		// Record is a JSON object, so we write as-is
// 		buffer.WriteString(string(queryResponse.Value))
// 		buffer.WriteString("}")
// 		bArrayMemberAlreadyWritten = true
// 	}
// 	buffer.WriteString("]")

// 	fmt.Printf("- getQueryResultForQueryString queryResult:\n%s\n", buffer.String())

// 	return buffer.Bytes(), nil
// }

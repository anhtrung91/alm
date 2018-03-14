/*
Copyright IBM Corp. 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"fmt"
	"strconv"

	"github.com/example_cc/controllers"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

var logger = shim.NewLogger("example_cc0")

// Chaincode example simple Chaincode implementation
type Chaincode struct {
}

func (t *Chaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	logger.Info("########### example_cc0 Init ###########")

	_, args := stub.GetFunctionAndParameters()
	var A, B string    // Entities
	var Aval, Bval int // Asset holdings
	var err error

	// Initialize the chaincode
	A = args[0]
	Aval, err = strconv.Atoi(args[1])
	if err != nil {
		return shim.Error("Expecting integer value for asset holding")
	}
	B = args[2]
	Bval, err = strconv.Atoi(args[3])
	if err != nil {
		return shim.Error("Expecting integer value for asset holding")
	}
	logger.Info("Aval = %d, Bval = %d\n", Aval, Bval)

	// Write the state to the ledger
	err = stub.PutState(A, []byte(strconv.Itoa(Aval)))
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(B, []byte(strconv.Itoa(Bval)))
	if err != nil {
		return shim.Error(err.Error())
	}

	// // Accounts
	// fmt.Println("Toodles Is Starting Up")
	// var accounts []Account
	// bytes, err := json.Marshal( accounts )

	// if err != nil {
	//   return shim.Error("Error initializing accounts.")
	// }

	// err = stub.PutState( "toodles_accounts", bytes )

	// // Users
	// var users []Users

	// bytes, err = json.Marshal( users )

	// if err != nil {
	//   return shim.Error("Error initializing users.")
	// }

	// err = stub.PutState( "toodles_users", bytes )

	return shim.Success(nil)
}

// Transaction makes payment of X units from A to B
func (t *Chaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	logger.Info("########### example_cc0 Invoke ###########")

	function, args := stub.GetFunctionAndParameters()

	if function == "delete" {
		// Deletes an entity from its state
		return t.delete(stub, args)
	}

	if function == "query" {
		// queries an entity state
		return t.query(stub, args)
	}
	if function == "move" {
		// Deletes an entity from its state
		return t.move(stub, args)
	}

	//goto function using account.go

	var account controllers.Account
	if function == "login" {
		return account.Login(stub, args)
	}

	if function == "get_al_user" {
		return account.GetAllUsers(stub, args)
	}

	if function == "create_account" {
		return account.CreateAccount(stub, args)
	}

	if function == "get_total_token" {
		return account.GetTotalToken(stub, args)
	}
	// end function account using

	// Merchant
	var merchant controllers.Merchant
	if function == "create_merchant" {
		return merchant.CreateMerchant(stub, args)
	}

	if function == "get_merchant" {
		return merchant.GetMerchant(stub, args)
	}

	if function == "get_all_merchant" {
		return merchant.GetAllMerchant(stub, args)
	}

	// Transaction
	var transaction controllers.Transaction
	if function == "create_transaction" {
		return transaction.CreateTransaction(stub, args)
	}

	if function == "get_all_transaction" {
		return transaction.GetAllTransaction(stub, args)
	}

	if function == "get_exchanged" {
		return transaction.GetExchanged(stub, args)
	}

	if function == "get_top_exchanged" {
		return transaction.GetTopExchanged(stub, args)
	}

	if function == "get_transaction_by_merchant" {
		return transaction.GetTransactionByMerchant(stub, args)
	}

	logger.Errorf("Unknown action, check the first argument, must be one of 'delete', 'query', or 'move'. But got: %v", args[0])
	return shim.Error(fmt.Sprintf("Unknown action, check the first argument, must be one of 'delete', 'query', or 'move'. But got: %v", args[0]))
}

func (t *Chaincode) move(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// must be an invoke
	var A, B string    // Entities
	var Aval, Bval int // Asset holdings
	var X int          // Transaction value
	var err error

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 4, function followed by 2 names and 1 value")
	}

	A = args[0]
	B = args[1]

	// Get the state from the ledger
	// TODO: will be nice to have a GetAllState call to ledger
	Avalbytes, err := stub.GetState(A)
	if err != nil {
		return shim.Error("Failed to get state")
	}

curl -s -X POST \
  http://localhost:4000/channels/mychannel/chaincodes/mycc \
  -H "content-type: application/json" \
  -d '{
	"peers": ["peer0.org1.example.com","peer1.org1.example.com"],
	"fcn":"get_top_exchanged",
	"args":["day"]
}'

curl -s -X POST \
  http://localhost:4000/channels/mychannel/chaincodes/mycc \
  -H "content-type: application/json" \
  -d '{
	"peers": ["peer0.org1.example.com","peer1.org1.example.com"],
	"fcn":"create_transaction",
	"args":["1","1","50","10","10","10","10","10","1"]
}'
}

func main() {
	err := shim.Start(new(Chaincode))
	if err != nil {
		logger.Errorf("Error starting Simple chaincode: %s", err)
	}
}

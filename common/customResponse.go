package common

import (
	pb "github.com/hyperledger/fabric/protos/peer"
)

const (
	SUCCESS = "200"
	ERR1    = "ALM0001"
)

var ResCodeDict = map[string]string{
	"200":     "OK",
	"ALM0001": "Error",
}

type ResponseSuccess struct {
	ResCode string
	Msg     string
	Payload string
}

type ResponseError struct {
	ResCode string
	Msg     string
}

func RespondSuccess(res *ResponseSuccess) pb.Response {
	if res.Payload == "" {
		res.Payload = "[]"
	}

	return pb.Response{
		Payload: []byte("{\"status\":\"" + res.ResCode + "\", \"msg\":\"" + res.Msg + "\", \"data\":" + res.Payload + "}"),
	}
}

func RespondError(err *ResponseError) pb.Response {
	return pb.Response{
		Payload: []byte("{\"status\":\"" + err.ResCode + "\", \"msg\":\"" + err.Msg + "\"}"),
	}
}

// func GetTotalUsers (total int) int {
// 	fmt.Println("GetTotalUsers: ", total)
// 	return i
// }

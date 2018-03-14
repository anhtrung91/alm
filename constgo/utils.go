package constgo

import (
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/example_cc/models"
)

type AppError models.AppError
type AppSuccess  models.AppSuccess

// type AppError struct {
//     Msg   		string	`json:"Msg"`
//     Status 		string	`json:"Status"`
//     Payload    	string	`json:"Payload"`
// }

func ReponseSuccess(appSuccess AppSuccess) pb.Response {
	return pb.Response{
		Payload: []byte("{status: " + appSuccess.Status + ",msg: " + appSuccess.Msg + ",data: " + appSuccess.Payload + "}"),
	}
}

func ReponseError(appError AppError) pb.Response {
	return pb.Response{
		Payload: []byte("{status: " + appError.Status + ",msg: " + appError.Msg + "}"),
	}
}

// func Reponse(msg string, status string, payload string) pb.Response {
// 	return pb.Response{
// 		Payload: []byte("{status: " + status + ",msg: " + msg + ",data: " + payload + "}"),
// 	}
// }

// func ReponseNotData(msg string, status string) pb.Response {
// 	return pb.Response{
// 		Payload: []byte("{status: " + status + ",msg: " + msg + "}"),
// 	}
// }
package client

import (
	"encoding/json"
)

type ArgsContent struct {
	Args []string `json:"args"`
}

type SignedData struct {
	Content       string `json:"content"`
	MessageNumber int    `json:"messageNumber"`
	UserName      string `json:"userName"`
}

type UserMessage struct {
	SignedData `json:"signedData"`

	Signature interface{} `json:"signature"`
}

func UnmarashalInnerArgs(objectString string) ([]string, error) {
	var argsContent ArgsContent

	json.Unmarshal([]byte(objectString), &argsContent)

	return argsContent.Args, nil
}

func UnmarashalUserMessage(objectString string) (UserMessage, error) {
	var userMessage UserMessage

	json.Unmarshal([]byte(objectString), &userMessage)

	return userMessage, nil
}

func (userMessage *UserMessage) VerifySignature(publicKey interface{}) bool {
	//
	return true
}

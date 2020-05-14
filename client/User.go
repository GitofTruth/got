package client

import (
	"encoding/json"
)

// This enum identifies the differnt user update request types.
type UserUpdateType int

const (
	CreateNewUser       UserUpdateType = 1
	ChangeUserUserName  UserUpdateType = 2
	ChangeUserPublicKey UserUpdateType = 3
	DeleteUser          UserUpdateType = 4
)

type UserUpdate struct {
	UserUpdateType `json:"userUpdateType"`
	UserInfo       `json:"userInfo"`

	OldUserName string `json:"oldUserName"`
}

func UnmarashalUserUpdate(objectString string) (UserUpdate, error) {
	var userUpdate UserUpdate

	json.Unmarshal([]byte(objectString), &userUpdate)

	return userUpdate, nil
}

// TODO: userName raceCase if userchange pupkey

// contains all the public avaiable information about a user
type UserInfo struct {
	UserName          string      `json:"userName"`
	PublicKey         interface{} `json:"publicKey"`
	LastMessageNumber int         `json:"lastMessageNumber"`
}

// contains all the avaiable information about a user
type User struct {
	UserInfo   `json:"userInfo"`
	PrivateKey interface{} `json:"privateKey"`
}

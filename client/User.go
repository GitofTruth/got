package client

import (
	"encoding/json"
)

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

type UserInfo struct {
	UserName          string      `json:"userName"`
	PublicKey         interface{} `json:"publicKey"`
	LastMessageNumber int         `json:"lastMessageNumber"`
}

type User struct {
	UserInfo   `json:"userInfo"`
	PrivateKey interface{} `json:"privateKey"`
}

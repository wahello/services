package messaging

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"google.golang.org/grpc/codes"
	"gorm.io/gorm"

	"github.com/gidyon/micro/v2/utils/errs"
	"github.com/gidyon/services/pkg/api/messaging"
)

const (
	messages = "messages"
)

// Message model
type Message struct {
	UserID      uint   `gorm:"index;not null"`
	Title       string `gorm:"type:varchar(256);not null"`
	Message     string `gorm:"type:varchar(2048);not null"`
	Link        string `gorm:"type:varchar(512);not null"`
	Seen        bool   `gorm:"type:tinyint(1);not null;default:0"`
	Type        int8   `gorm:"type:tinyint(1);not null;default:0"`
	SendMethods []byte `gorm:"type:json;not null"`
	Details     []byte `gorm:"type:json;null"`
	gorm.Model
}

// TableName returns the name of the table
func (*Message) TableName() string {
	tableName := os.Getenv("MESSAGING_TABLE")
	if tableName != "" {
		return tableName
	}
	return messages
}

// GetMessageDB creates message model from proto message
func GetMessageDB(messagePB *messaging.Message) (*Message, error) {
	id, err := strconv.Atoi(messagePB.UserId)
	if err != nil {
		return nil, errs.WrapErrorWithCodeAndMsg(codes.InvalidArgument, err, "failed to get message")
	}

	messageDB := &Message{
		UserID:  uint(id),
		Title:   messagePB.Title,
		Message: messagePB.Data,
		Link:    messagePB.Link,
		Seen:    messagePB.Seen,
		Type:    int8(messagePB.Type),
	}

	var data []byte

	// SendMethods
	if len(messagePB.SendMethods) != 0 {
		data, err = json.Marshal(messagePB.SendMethods)
		if err != nil {
			return nil, errs.FromJSONMarshal(err, "SendMethods")
		}
		messageDB.SendMethods = data
	}

	// Details
	if len(messagePB.Details) != 0 {
		data, err = json.Marshal(messagePB.Details)
		if err != nil {
			return nil, errs.FromJSONMarshal(err, "Details")
		}
		messageDB.Details = data
	}

	return messageDB, nil
}

// GetMessagePB creates proto message from message model
func GetMessagePB(messageDB *Message) (*messaging.Message, error) {
	messagePB := &messaging.Message{
		MessageId:         fmt.Sprint(messageDB.ID),
		UserId:            fmt.Sprint(messageDB.UserID),
		Title:             messageDB.Title,
		Data:              messageDB.Message,
		Link:              messageDB.Link,
		CreateTimeSeconds: messageDB.CreatedAt.Unix(),
		Seen:              messageDB.Seen,
		Save:              true,
		Type:              messaging.MessageType(messageDB.Type),
		SendMethods:       make([]messaging.SendMethod, 0),
	}

	// SendMethod
	if len(messageDB.SendMethods) != 0 {
		data := make([]messaging.SendMethod, 0)
		err := json.Unmarshal(messageDB.SendMethods, &data)
		if err != nil {
			return nil, errs.FromJSONUnMarshal(err, "Details")
		}
		messagePB.SendMethods = data
	}

	// Details
	if len(messageDB.Details) != 0 {
		data := make(map[string]string)
		err := json.Unmarshal(messageDB.Details, &data)
		if err != nil {
			return nil, errs.FromJSONUnMarshal(err, "Details")
		}
		messagePB.Details = data
	}

	return messagePB, nil
}

package services

//go:generate mkdir -p ../protos/generated
//go:generate sh -c "protoc --go_out=paths=source_relative,plugins=grpc:../protos/generated -I=../protos ../protos/*.proto"

import (
	"context"

	empty "github.com/golang/protobuf/ptypes/empty"
	"github.com/pojntfx/nextcloud-talk-bot-framework/pkg/clients"
	nextcloudTalk "github.com/pojntfx/nextcloud-talk-bot-framework/pkg/protos/generated"
)

// NextcloudTalk is a Nextcloud Talk client.
type NextcloudTalk struct {
	nextcloudTalk.UnimplementedNextcloudTalkServer
	chatRequestChan  chan bool
	chatResponseChan chan chan clients.Chat
	writeChat        func(token, message string) error
}

// NewNextcloudTalk creates a new Nextcloud Talk Client.
func NewNextcloudTalk(chatRequestChan chan bool, chatResponseChan chan chan clients.Chat, writeChat func(token, message string) error) *NextcloudTalk {
	return &NextcloudTalk{
		chatRequestChan:  chatRequestChan,
		chatResponseChan: chatResponseChan,
		writeChat:        writeChat,
	}
}

// ReadChats reads the chats.
func (n *NextcloudTalk) ReadChats(req *empty.Empty, srv nextcloudTalk.NextcloudTalk_ReadChatsServer) error {
	n.chatRequestChan <- true

	readChan := <-n.chatResponseChan

	for chat := range readChan {
		if err := srv.Send(&nextcloudTalk.OutChat{
			ID:               int64(chat.ID),
			Token:            chat.Token,
			ActorID:          chat.ActorID,
			ActorDisplayName: chat.ActorDisplayName,
			Message:          chat.Message,
		}); err != nil {
			return err
		}
	}

	return nil
}

// WriteChat writes a chat.
func (n *NextcloudTalk) WriteChat(ctx context.Context, req *nextcloudTalk.InChat) (*empty.Empty, error) {
	if err := n.writeChat(req.GetToken(), req.GetMessage()); err != nil {
		return &empty.Empty{}, err
	}

	return &empty.Empty{}, nil
}

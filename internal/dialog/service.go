package dialog

import (
	"github.com/google/uuid"
)

type DialogService struct {
	repository *DialogRepository
}

func NewDialogService(repository *DialogRepository) *DialogService {
	return &DialogService{repository: repository}
}

func (d *DialogService) SendDirectMessage(senderId uuid.UUID, recipientId uuid.UUID, text string) (*Message, error) {
	return d.repository.SendDirectMessage(senderId, recipientId, text)
}

func (d *DialogService) ListLastDirectMessages(senderId, recipientId uuid.UUID, offset, limit uint) ([]*Message, error) {
	dialog, err := d.ResolvePrivateDialog(senderId, recipientId)
	if err != nil {
		return nil, err
	}

	messages, err := d.repository.GetPrivateDialogMessages(dialog.Id, offset, limit)

	if err != nil {
		return nil, err
	}

	return messages, nil
}

func (d *DialogService) ResolvePrivateDialog(senderId, recipientId uuid.UUID) (*Dialog, error) {
	return d.repository.GetPrivateDialog(senderId, recipientId)
}

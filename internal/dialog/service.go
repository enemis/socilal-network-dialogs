package dialog

type DialogService struct {
	repository *DialogRepository
}

func NewDialogService(repository *DialogRepository) *DialogService {
	return &DialogService{repository: repository}
}

func (d *DialogService) SendDirectMessage(senderId string, recipientId string, text string) (*Message, error) {
	return d.repository.SendDirectMessage(senderId, recipientId, text)
}

func (d *DialogService) ListLastDirectMessages(senderId, recipientId string, offset, limit uint) ([]*Message, error) {
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

func (d *DialogService) ResolvePrivateDialog(senderId, recipientId string) (*Dialog, error) {
	return d.repository.GetPrivateDialog(senderId, recipientId)
}

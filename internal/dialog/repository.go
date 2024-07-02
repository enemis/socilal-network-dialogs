package dialog

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"social-network-dialogs/internal/database"
	"time"
)

type DialogRepository struct {
	db *database.DatabaseStack
}

func NewDialogRepository(db *database.DatabaseStack) *DialogRepository {
	return &DialogRepository{
		db: db,
	}
}

func (r *DialogRepository) GetPrivateDialog(user1 uuid.UUID, user2 uuid.UUID) (*Dialog, error) {
	var dialog Dialog
	query := "SELECT d.* FROM dialog_participants dp1 " +
		"INNER JOIN dialog_participants dp2 ON dp1.dialog_id = dp2.dialog_id " +
		"INNER JOIN dialogs d ON dp1.dialog_id = d.id AND dp2.dialog_id = d.id AND d.type = 'direct' AND d.deleted_at IS NULL " +
		"WHERE dp1.user_id=$1 AND dp2.user_id=$2 " +
		"LIMIT 1;"

	err := r.db.GetReadConnection().Get(&dialog, query, user1, user2)
	if err != nil {
		dialog.Type = Direct
		dlg, err := r.CreateDialog(&dialog, user1, user2)
		if err != nil {
			return nil, err
		}
		return dlg, nil
	}
	return &dialog, nil
}

func (r *DialogRepository) CreateDialog(dialog *Dialog, participants ...uuid.UUID) (*Dialog, error) {
	dialog.Id = uuid.New()

	tx := r.db.GetWriteConnection().MustBeginTx(context.Background(), nil)

	query := "INSERT INTO dialogs (id, type, created_at) VALUES ($1, $2, $3)"

	_, err := tx.Exec(query, dialog.Id, dialog.Type, time.Now())
	if err != nil {
		err := tx.Rollback()
		return nil, errors.Wrap(err, "error create dialog")
	}

	query = "INSERT INTO dialog_participants (dialog_id, user_id) VALUES "
	for _, participant := range participants {
		query += fmt.Sprintf("('%s', '%s'),", dialog.Id, participant)
	}
	query = query[:len(query)-1] //remove last ,

	_, err = tx.Exec(query)
	if err != nil {
		err := tx.Rollback()
		return nil, errors.Wrap(err, "error create dialog")
	}

	err = tx.Commit()
	if err != nil {
		return nil, errors.Wrap(err, "error create dialog")
	}

	dialog, err = r.GetDialog(dialog.Id)
	if err != nil {
		return nil, err
	}

	return dialog, err
}

func (r *DialogRepository) GetDialog(dialogId uuid.UUID) (*Dialog, error) {
	var dialog Dialog
	tx := r.db.GetWriteConnection().MustBeginTx(context.Background(), nil)
	err := tx.Get(&dialog, "SELECT * FROM dialogs WHERE id=$1 and deleted_at is null LIMIT 1", dialogId)

	if err != nil {
		tx.Rollback()
		return nil, errors.New("dialog not found")
	}

	rows, err := tx.Queryx("SELECT * FROM dialog_paricipants WHERE dialog_id=$1 and deleted_at is null LIMIT 1", dialogId)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	participiants := make([]*Participant, 0)
	for rows.Next() {
		var participant Participant
		err = rows.StructScan(&participant)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		participiants = append(participiants, &participant)
	}
	tx.Commit()
	dialog.Participants = participiants

	return &dialog, nil
}

func (r *DialogRepository) SendDirectMessage(recipientId uuid.UUID, senderId uuid.UUID, text string) (*Message, error) {
	dialog, err := r.GetPrivateDialog(recipientId, senderId)
	if err != nil {
		return nil, err
	}

	messageId := uuid.New()

	query := "INSERT INTO messages (id, user_id, dialog_id, message, created_at) VALUES ($1, $2, $3, $4, $5)"

	_, err = r.db.GetWriteConnection().Exec(query, messageId, senderId, dialog.Id, text, time.Now())

	if err != nil {
		return nil, errors.Wrap(err, "error send message")
	}

	message, err := r.GetMessage(messageId)
	if err != nil {
		return nil, err
	}

	return message, err
}

func (r *DialogRepository) GetMessage(messageId uuid.UUID) (*Message, error) {
	var message Message
	err := r.db.GetReadConnection().Get(&message, "SELECT * FROM messages WHERE id=$1 and deleted_at is null LIMIT 1", messageId)

	if err != nil {
		return nil, errors.New("message not found")
	}

	return &message, err
}

func (r *DialogRepository) GetPrivateDialogMessages(dialogId uuid.UUID, offset, limit uint) ([]*Message, error) {
	var results = make([]*Message, 0, limit)
	//todo cursors
	query := "SELECT m.* FROM dialogs d " +
		"INNER JOIN messages m ON m.dialog_id = d.id " +
		"WHERE d.id=$1 AND d.deleted_at IS NULL AND d.type = 'direct' " +
		"ORDER BY m.created_at DESC LIMIT $2 OFFSET $3"

	rows, err := r.db.GetReadConnection().Queryx(query, dialogId, limit, offset)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var message Message
		err = rows.StructScan(&message)
		if err != nil {
			return nil, err
		}

		results = append(results, &message)
	}

	return results, nil
}

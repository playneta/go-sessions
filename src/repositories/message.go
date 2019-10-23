package repositories

import (
	"github.com/go-pg/pg"
	"github.com/playneta/go-sessions/src/models"
)

type (
	Message interface {
		Create(message *models.Message) error
		LastPublicMessages(limit int) ([]models.Message, error)
		LastPrivateMessages(user models.User, limit int) ([]models.Message, error)
	}

	messageRepository struct {
		db *pg.DB
	}
)

func NewMessage(db *pg.DB) Message {
	return &messageRepository{
		db: db,
	}
}

func (m *messageRepository) Create(message *models.Message) error {
	if _, err := m.db.Model(message).Relation("User").Insert(); err != nil {
		return err
	}

	if err := m.db.Model(message).
		Column("message.*").
		Relation("User").Relation("Receiver").
		Where("message.id=?", message.Id).First(); err != nil {
		return err
	}

	return nil
}

func (m *messageRepository) LastPublicMessages(limit int) ([]models.Message, error) {
	var messages []models.Message
	if err := m.db.Model(&messages).
		Column("message.*").
		Where("receiver_id IS NULL").
		Order("id desc").
		Relation("Receiver").
		Relation("User").
		Limit(limit).Select(); err != nil {
		if err == pg.ErrNoRows {
			return []models.Message{}, nil
		}

		return nil, err
	}

	return messages, nil
}

func (m *messageRepository) LastPrivateMessages(user models.User, limit int) ([]models.Message, error) {
	var messages []models.Message
	if err := m.db.Model(&messages).
		Column("message.*").
		Where("user_id=? and (receiver_id>0 or receiver_id=?)", user.ID, user.ID).
		Order("id desc").
		Relation("Receiver").
		Relation("User").
		Limit(limit).Select(); err != nil {
		if err == pg.ErrNoRows {
			return []models.Message{}, nil
		}

		return nil, err
	}

	return messages, nil
}

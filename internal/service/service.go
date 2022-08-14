package service

import "context"

type Reciver struct {
	msgStorer     MessagesStorer
	profileStorer ProfilesStorer
}

type MessagesStorer interface {
	Save(ctx context.Context, uid int64, msg string) error
	GetMessagesByUID(ctx context.Context, uid int64) ([]string, error)
}

type ProfilesStorer interface {
	SaveOrUpdate(ctx context.Context, uid int64, fName, lName, username string) error
}

func (r *Reciver) StoreMessage(ctx context.Context, uid int64, msg string) error {
	if err := r.msgStorer.Save(ctx, uid, msg); err != nil {
		return err
	}
	return nil
}

func (r *Reciver) StoreOrUpdateProfile(ctx context.Context, uid int64, fName, lName, username string) error {
	if err := r.profileStorer.SaveOrUpdate(ctx, uid, fName, lName, username); err != nil {
		return err
	}
	return nil
}

func (r *Reciver) GetMessagesByUID(ctx context.Context, uid int64) ([]string, error) {
	msgs, err := r.msgStorer.GetMessagesByUID(ctx, uid)
	if err != nil {
		return nil, err
	}
	return msgs, nil
}

func NewService(msgStorer MessagesStorer, profileStorer ProfilesStorer) *Reciver {
	return &Reciver{
		msgStorer:     msgStorer,
		profileStorer: profileStorer,
	}
}

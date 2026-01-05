package transactor

import (
	"context"

	mongodb "go.mongodb.org/mongo-driver/mongo"
)

type mongo struct {
	client *mongodb.Client
}

func NewMongo(client *mongodb.Client) Transactor {
	return &mongo{
		client: client,
	}
}

func (t *mongo) WithTransaction(
	ctx context.Context,
	fn func(ctx context.Context) error,
) error {
	if sc, ok := ctx.(mongodb.SessionContext); ok {
		return fn(sc)
	}

	session, err := t.client.StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, func(sc mongodb.SessionContext) (interface{}, error) {
		if err := fn(sc); err != nil {
			return nil, err
		}
		return nil, nil
	})

	return err
}

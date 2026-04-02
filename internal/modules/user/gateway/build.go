package gateway

import (
	"context"

	"reverie.jp/reverie/internal/platform/ulid"
)

func (g *gatewayImpl) BuildView(ctx context.Context, id ulid.ULID) (*UserView, error) {
	views, err := g.BuildListViews(ctx, []ulid.ULID{id})
	if err != nil {
		return nil, err
	}
	if len(views) == 0 {
		return nil, nil
	}
	return views[0], nil
}

func (g *gatewayImpl) BuildListViews(ctx context.Context, ids []ulid.ULID) ([]*UserView, error) {
	if len(ids) == 0 {
		return []*UserView{}, nil
	}

	users, err := g.repo.ListUsersByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}

	views := make([]*UserView, len(users))
	for i, u := range users {
		views[i] = &UserView{User: u}
	}

	return views, nil
}

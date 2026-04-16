package gateway

import (
	"context"
	"time"

	"reverie.jp/reverie/internal/domain/entity"
	"reverie.jp/reverie/internal/gen/sqlc"
	"reverie.jp/reverie/internal/modules/post/repository"
	usergw "reverie.jp/reverie/internal/modules/user/gateway"
	"reverie.jp/reverie/internal/platform/ulid"
)

type PostView struct {
	Post          *entity.Post
	Author        *entity.User
	ReplyCount    int64
	RepostCount   int64
	FavoriteCount int64
	IsFavorited   bool
}

type CreatePostParams struct {
	ID        ulid.ULID
	AuthorID  ulid.ULID
	ReplyToID *ulid.ULID
	RepostID  *ulid.ULID
	Text      string
}

type ListTimelineParams struct {
	Cursor *time.Time
	Limit  int32
}

type ListUserPostsParams struct {
	AuthorID ulid.ULID
	Cursor   *time.Time
	Limit    int32
}

type Gateway interface {
	CreatePost(ctx context.Context, params CreatePostParams) (*PostView, error)
	DeletePost(ctx context.Context, postID ulid.ULID, authorID ulid.ULID) error
	ListTimeline(ctx context.Context, params ListTimelineParams, requestorID ulid.ULID) ([]*PostView, error)
	ListUserPosts(ctx context.Context, params ListUserPostsParams, requestorID ulid.ULID) ([]*PostView, error)
}

type gatewayImpl struct {
	repo        repository.Repository
	userGateway usergw.Gateway
}

func New(q sqlc.Querier, userGateway usergw.Gateway) Gateway {
	return &gatewayImpl{
		repo:        repository.New(q),
		userGateway: userGateway,
	}
}

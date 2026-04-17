package repository

import (
	"context"
	"time"

	"reverie.jp/reverie/internal/domain/entity"
	"reverie.jp/reverie/internal/gen/sqlc"
	"reverie.jp/reverie/internal/platform/ulid"
)

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

type Repository interface {
	CreatePost(ctx context.Context, params CreatePostParams) (*entity.Post, error)
	DeletePost(ctx context.Context, postID ulid.ULID, authorID ulid.ULID) error
	GetPostByID(ctx context.Context, id ulid.ULID) (*entity.Post, error)
	ListTimeline(ctx context.Context, params ListTimelineParams) ([]*entity.Post, error)
	ListFollowingTimeline(ctx context.Context, params ListFollowingTimelineParams) ([]*entity.Post, error)
	ListUserPosts(ctx context.Context, params ListUserPostsParams) ([]*entity.Post, error)
	ListPostReposts(ctx context.Context, params ListPostRepostsParams) ([]*entity.Post, error)
	ListPostReplies(ctx context.Context, params ListPostRepliesParams) ([]*entity.Post, error)
	CountPostReplies(ctx context.Context, postID ulid.ULID) (int64, error)
	CountPostReposts(ctx context.Context, postID ulid.ULID) (int64, error)
	CountPostFavorites(ctx context.Context, postID ulid.ULID) (int64, error)
	GetPostFavorite(ctx context.Context, userID ulid.ULID, postID ulid.ULID) (bool, error)
	CreatePostFavorite(ctx context.Context, userID ulid.ULID, postID ulid.ULID) error
	DeletePostFavorite(ctx context.Context, userID ulid.ULID, postID ulid.ULID) error
}

type repositoryImpl struct {
	q sqlc.Querier
}

func New(q sqlc.Querier) Repository {
	return &repositoryImpl{q: q}
}

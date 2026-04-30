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
	Author        *usergw.UserView
	ReplyCount    int64
	RepostCount   int64
	FavoriteCount int64
	IsFavorited   bool
	RepostOf      *PostView
}

type CreatePostParams struct {
	ID        ulid.ULID
	AuthorID  ulid.ULID
	ReplyToPostID *ulid.ULID
	RepostPostID  *ulid.ULID
	Text      string
}

type ListTimelineParams struct {
	Cursor *time.Time
	Limit  int32
}

type ListFollowingTimelineParams struct {
	FollowerID ulid.ULID
	Cursor     *time.Time
	Limit      int32
}

type ListUserPostsParams struct {
	AuthorID ulid.ULID
	Cursor   *time.Time
	Limit    int32
}

type ListUserLikedPostsParams struct {
	UserID ulid.ULID
	Cursor *time.Time
	Limit  int32
}

type ListPostRepliesParams struct {
	Cursor *time.Time
	Limit  int32
}

type ListPostRepostsParams struct {
	Cursor *time.Time
	Limit  int32
}

type Gateway interface {
	// GetPost は authorCustomID と shortID で投稿を取得する。
	GetPost(ctx context.Context, authorCustomID string, shortID string, requestorID ulid.ULID) (*PostView, error)
	CreatePost(ctx context.Context, params CreatePostParams) (*PostView, error)
	// DeletePost は authorCustomID と shortID で投稿を特定して削除する。
	DeletePost(ctx context.Context, authorCustomID string, shortID string, authorID ulid.ULID) error
	LikePost(ctx context.Context, authorCustomID string, shortID string, userID ulid.ULID) (*PostView, error)
	UnlikePost(ctx context.Context, authorCustomID string, shortID string, userID ulid.ULID) (*PostView, error)
	ListTimeline(ctx context.Context, params ListTimelineParams, requestorID ulid.ULID) ([]*PostView, error)
	ListFollowingTimeline(ctx context.Context, params ListFollowingTimelineParams, requestorID ulid.ULID) ([]*PostView, error)
	ListUserPosts(ctx context.Context, params ListUserPostsParams, requestorID ulid.ULID) ([]*PostView, error)
	ListUserLikedPosts(ctx context.Context, params ListUserLikedPostsParams, requestorID ulid.ULID) ([]*PostView, error)
	ListPostReposts(ctx context.Context, authorCustomID string, shortID string, params ListPostRepostsParams, requestorID ulid.ULID) ([]*PostView, error)
	ListPostReplies(ctx context.Context, authorCustomID string, shortID string, params ListPostRepliesParams, requestorID ulid.ULID) ([]*PostView, error)
	ListPostLikes(ctx context.Context, authorCustomID string, shortID string, requestorID ulid.ULID, limit int32) ([]*usergw.UserView, error)
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

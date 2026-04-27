package handler

import (
	"context"
	"encoding/base64"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/timestamppb"
	"reverie.jp/reverie/internal/application/server/interceptor"
	postv1 "reverie.jp/reverie/internal/gen/pb/post/v1"
	timelinev1 "reverie.jp/reverie/internal/gen/pb/timeline/v1"
	userv1 "reverie.jp/reverie/internal/gen/pb/user/v1"
	"reverie.jp/reverie/internal/modules/post/usecase"
	"reverie.jp/reverie/internal/platform/xerrors"
)

func (h *Handler) ListFollowingTimeline(ctx context.Context, req *connect.Request[timelinev1.ListFollowingTimelineRequest]) (*connect.Response[timelinev1.ListFollowingTimelineResponse], error) {
	userID, ok := interceptor.UserIDFromContext(ctx)
	if !ok {
		return nil, xerrors.ErrUnauthenticated
	}

	outputs, err := h.listFollowingTimeline.Execute(ctx, usecase.ListTimelineInput{
		PageToken: req.Msg.PageToken,
		PageSize:  req.Msg.PageSize,
	}, userID)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&timelinev1.ListFollowingTimelineResponse{
		Posts:         toProtoPosts(outputs),
		NextPageToken: nextPageToken(outputs),
	}), nil
}

func (h *Handler) ListPublicTimeline(ctx context.Context, req *connect.Request[timelinev1.ListPublicTimelineRequest]) (*connect.Response[timelinev1.ListPublicTimelineResponse], error) {
	userID, ok := interceptor.UserIDFromContext(ctx)
	if !ok {
		return nil, xerrors.ErrUnauthenticated
	}

	outputs, err := h.listTimeline.Execute(ctx, usecase.ListTimelineInput{
		PageToken: req.Msg.PageToken,
		PageSize:  req.Msg.PageSize,
	}, userID)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&timelinev1.ListPublicTimelineResponse{
		Posts:         toProtoPosts(outputs),
		NextPageToken: nextPageToken(outputs),
	}), nil
}

func toProtoPost(o *usecase.PostOutput) *postv1.Post {
	if o == nil {
		return nil
	}
	p := &postv1.Post{
		Id:          o.ID.String(),
		Text:        o.Text,
		ReplyCount:  int32(o.ReplyCount),
		RepostCount: int32(o.RepostCount),
		LikeCount:   int32(o.FavoriteCount),
		IsLiked:     o.IsFavorited,
		CreateTime:  timestamppb.New(o.CreateTime),
	}
	if o.Author != nil && o.Author.User != nil {
		u := o.Author.User
		bio := u.Biography
		p.Author = &userv1.User{
			Id:           u.ID.String(),
			CustomId:     u.CustomID,
			DisplayName:  u.DisplayName,
			Biography:    &bio,
			IsPrivate:    u.IsPrivate,
			IsMe:         o.Author.IsMe,
			IsFollowing:  o.Author.IsFollowing,
			IsFollowedBy: o.Author.IsFollowedBy,
			CreateTime:   timestamppb.New(u.CreateTime),
		}
	}
	if o.ReplyToID != nil {
		s := o.ReplyToID.String()
		p.ReplyToId = &s
	}
	if o.RepostID != nil {
		s := o.RepostID.String()
		p.RepostId = &s
	}
	if o.RepostOf != nil {
		p.RepostOf = toProtoPost(o.RepostOf)
	}
	return p
}

func toProtoPosts(outputs []*usecase.PostOutput) []*postv1.Post {
	posts := make([]*postv1.Post, len(outputs))
	for i, o := range outputs {
		posts[i] = toProtoPost(o)
	}
	return posts
}

func nextPageToken(posts []*usecase.PostOutput) string {
	if len(posts) == 0 {
		return ""
	}
	last := posts[len(posts)-1]
	raw := last.CreateTime.UTC().Format("2006-01-02T15:04:05.999999999Z")
	return base64.StdEncoding.EncodeToString([]byte(raw))
}

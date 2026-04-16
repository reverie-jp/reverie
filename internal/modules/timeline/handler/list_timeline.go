package handler

import (
	"context"

	"connectrpc.com/connect"
	"reverie.jp/reverie/internal/application/server/interceptor"
	timelinev1 "reverie.jp/reverie/internal/gen/pb/timeline/v1"
	postv1 "reverie.jp/reverie/internal/gen/pb/post/v1"
	"reverie.jp/reverie/internal/modules/post/usecase"
	"reverie.jp/reverie/internal/platform/xerrors"
	"google.golang.org/protobuf/types/known/timestamppb"
	userv1 "reverie.jp/reverie/internal/gen/pb/user/v1"
	"encoding/base64"
)

func (h *Handler) ListFollowingTimeline(ctx context.Context, req *connect.Request[timelinev1.ListFollowingTimelineRequest]) (*connect.Response[timelinev1.ListFollowingTimelineResponse], error) {
	userID, ok := interceptor.UserIDFromContext(ctx)
	if !ok {
		return nil, xerrors.ErrUnauthenticated
	}

	// フォロー機能が未実装のためパブリックタイムラインと同じ動作
	outputs, err := h.listTimeline.Execute(ctx, usecase.ListTimelineInput{
		PageToken: req.Msg.PageToken,
		PageSize:  req.Msg.PageSize,
	}, userID)
	if err != nil {
		return nil, err
	}

	posts := toProtoPosts(outputs)

	return connect.NewResponse(&timelinev1.ListFollowingTimelineResponse{
		Posts:         posts,
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

	posts := toProtoPosts(outputs)

	return connect.NewResponse(&timelinev1.ListPublicTimelineResponse{
		Posts:         posts,
		NextPageToken: nextPageToken(outputs),
	}), nil
}

func toProtoPosts(outputs []*usecase.PostOutput) []*postv1.Post {
	posts := make([]*postv1.Post, len(outputs))
	for i, o := range outputs {
		p := &postv1.Post{
			Id:          o.ID.String(),
			Text:        o.Text,
			ReplyCount:  int32(o.ReplyCount),
			RepostCount: int32(o.RepostCount),
			LikeCount:   int32(o.FavoriteCount),
			IsLiked:     o.IsFavorited,
			CreateTime:  timestamppb.New(o.CreateTime),
		}
		if o.Author != nil {
			p.Author = &userv1.User{
				Id:          o.Author.ID.String(),
				CustomId:    o.Author.CustomID,
				DisplayName: o.Author.DisplayName,
				IsPrivate:   o.Author.IsPrivate,
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
		posts[i] = p
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

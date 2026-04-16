package handler

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
	postv1 "reverie.jp/reverie/internal/gen/pb/post/v1"
	userv1 "reverie.jp/reverie/internal/gen/pb/user/v1"
	"reverie.jp/reverie/internal/modules/post/usecase"
)

func toProtoPost(out *usecase.PostOutput) *postv1.Post {
	if out == nil {
		return nil
	}

	p := &postv1.Post{
		Id:            out.ID.String(),
		Text:          out.Text,
		ReplyCount:    int32(out.ReplyCount),
		RepostCount:   int32(out.RepostCount),
		FavoriteCount: int32(out.FavoriteCount),
		IsFavorited:   out.IsFavorited,
		CreateTime:    timestamppb.New(out.CreateTime),
	}

	if out.Author != nil {
		p.Author = &userv1.User{
			Id:          out.Author.ID.String(),
			CustomId:    out.Author.CustomID,
			DisplayName: out.Author.DisplayName,
			IsPrivate:   out.Author.IsPrivate,
		}
	}

	if out.ReplyToID != nil {
		s := out.ReplyToID.String()
		p.ReplyToId = &s
	}

	if out.RepostID != nil {
		s := out.RepostID.String()
		p.RepostId = &s
	}

	return p
}

func nextCursor(posts []*usecase.PostOutput) *string {
	if len(posts) == 0 {
		return nil
	}
	last := posts[len(posts)-1]
	s := last.CreateTime.UTC().Format(time.RFC3339Nano)
	return &s
}

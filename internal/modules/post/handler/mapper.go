package handler

import (
	"encoding/base64"

	"google.golang.org/protobuf/types/known/timestamppb"
	postv1 "reverie.jp/reverie/internal/gen/pb/post/v1"
	userv1 "reverie.jp/reverie/internal/gen/pb/user/v1"
	"reverie.jp/reverie/internal/modules/post/usecase"
	usergw "reverie.jp/reverie/internal/modules/user/gateway"
)

func toProtoPost(out *usecase.PostOutput) *postv1.Post {
	if out == nil {
		return nil
	}

	p := &postv1.Post{
		Id:          out.ID.String(),
		Text:        out.Text,
		ReplyCount:  int32(out.ReplyCount),
		RepostCount: int32(out.RepostCount),
		LikeCount:   int32(out.FavoriteCount),
		IsLiked:     out.IsFavorited,
		CreateTime:  timestamppb.New(out.CreateTime),
	}

	if out.Author != nil && out.Author.User != nil {
		u := out.Author.User
		bio := u.Biography
		p.Author = &userv1.User{
			Id:           u.ID.String(),
			CustomId:     u.CustomID,
			DisplayName:  u.DisplayName,
			Biography:    &bio,
			IsPrivate:    u.IsPrivate,
			IsMe:         out.Author.IsMe,
			IsFollowing:  out.Author.IsFollowing,
			IsFollowedBy: out.Author.IsFollowedBy,
			CreateTime:   timestamppb.New(u.CreateTime),
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

	if out.RepostOf != nil {
		p.RepostOf = toProtoPost(out.RepostOf)
	}

	return p
}

func toProtoUserFromView(v *usergw.UserView) *userv1.User {
	if v == nil || v.User == nil {
		return nil
	}
	u := v.User
	bio := u.Biography
	return &userv1.User{
		Id:           u.ID.String(),
		CustomId:     u.CustomID,
		DisplayName:  u.DisplayName,
		Biography:    &bio,
		IsPrivate:    u.IsPrivate,
		IsMe:         v.IsMe,
		IsFollowing:  v.IsFollowing,
		IsFollowedBy: v.IsFollowedBy,
		CreateTime:   timestamppb.New(u.CreateTime),
	}
}

func nextPageToken(posts []*usecase.PostOutput) string {
	if len(posts) == 0 {
		return ""
	}
	last := posts[len(posts)-1]
	raw := last.CreateTime.UTC().Format("2006-01-02T15:04:05.999999999Z")
	return base64.StdEncoding.EncodeToString([]byte(raw))
}

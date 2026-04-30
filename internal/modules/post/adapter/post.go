package adapter

import (
	"encoding/base64"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/timestamppb"

	postv1 "reverie.jp/reverie/internal/gen/pb/post/v1"
	useradapter "reverie.jp/reverie/internal/modules/user/adapter"
	"reverie.jp/reverie/internal/modules/post/usecase"
)

func ToPost(out *usecase.PostOutput) *postv1.Post {
	if out == nil {
		return nil
	}
	p := &postv1.Post{
		Id:          out.ID.String(),
		ShortId:     out.ShortID,
		Text:        out.Text,
		Author:      useradapter.ToUser(out.Author),
		ReplyCount:  int32(out.ReplyCount),
		RepostCount: int32(out.RepostCount),
		LikeCount:   int32(out.FavoriteCount),
		IsLiked:     out.IsFavorited,
		CreateTime:  timestamppb.New(out.CreateTime),
	}
	if out.ReplyToPostID != nil {
		s := out.ReplyToPostID.String()
		p.ReplyToPostId = &s
	}
	if out.RepostPostID != nil {
		s := out.RepostPostID.String()
		p.RepostPostId = &s
	}
	if out.RepostOf != nil {
		p.RepostOf = ToPost(out.RepostOf)
	}
	return p
}

func ToPosts(outs []*usecase.PostOutput) []*postv1.Post {
	posts := make([]*postv1.Post, len(outs))
	for i, o := range outs {
		posts[i] = ToPost(o)
	}
	return posts
}

func NextPageToken(posts []*usecase.PostOutput) string {
	if len(posts) == 0 {
		return ""
	}
	last := posts[len(posts)-1]
	raw := last.CreateTime.UTC().Format("2006-01-02T15:04:05.999999999Z")
	return base64.StdEncoding.EncodeToString([]byte(raw))
}

func ToPostResponse(out *usecase.PostOutput) *connect.Response[postv1.GetPostResponse] {
	return connect.NewResponse(&postv1.GetPostResponse{Post: ToPost(out)})
}

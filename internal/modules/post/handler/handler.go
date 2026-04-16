package handler

import (
	"reverie.jp/reverie/internal/gen/pb/post/v1/postv1connect"
	"reverie.jp/reverie/internal/modules/post/usecase"
)

type Handler struct {
	postv1connect.UnimplementedPostServiceHandler
	createPost    *usecase.CreatePost
	deletePost    *usecase.DeletePost
	listTimeline  *usecase.ListTimeline
	listUserPosts *usecase.ListUserPosts
}

func New(
	createPost *usecase.CreatePost,
	deletePost *usecase.DeletePost,
	listTimeline *usecase.ListTimeline,
	listUserPosts *usecase.ListUserPosts,
) *Handler {
	return &Handler{
		createPost:    createPost,
		deletePost:    deletePost,
		listTimeline:  listTimeline,
		listUserPosts: listUserPosts,
	}
}

package handler

import (
	"reverie.jp/reverie/internal/gen/pb/post/v1/postv1connect"
	"reverie.jp/reverie/internal/modules/post/usecase"
)

type Handler struct {
	postv1connect.UnimplementedPostServiceHandler
	getPost            *usecase.GetPost
	createPost         *usecase.CreatePost
	deletePost         *usecase.DeletePost
	likePost           *usecase.LikePost
	unlikePost         *usecase.UnlikePost
	listUserPosts      *usecase.ListUserPosts
	listPostReplies    *usecase.ListPostReplies
	listPostReposts    *usecase.ListPostReposts
	listPostLikes      *usecase.ListPostLikes
	listUserLikedPosts *usecase.ListUserLikedPosts
}

func New(
	getPost *usecase.GetPost,
	createPost *usecase.CreatePost,
	deletePost *usecase.DeletePost,
	likePost *usecase.LikePost,
	unlikePost *usecase.UnlikePost,
	listUserPosts *usecase.ListUserPosts,
	listPostReplies *usecase.ListPostReplies,
	listPostReposts *usecase.ListPostReposts,
	listPostLikes *usecase.ListPostLikes,
	listUserLikedPosts *usecase.ListUserLikedPosts,
) *Handler {
	return &Handler{
		getPost:            getPost,
		createPost:         createPost,
		deletePost:         deletePost,
		likePost:           likePost,
		unlikePost:         unlikePost,
		listUserPosts:      listUserPosts,
		listPostReplies:    listPostReplies,
		listPostReposts:    listPostReposts,
		listPostLikes:      listPostLikes,
		listUserLikedPosts: listUserLikedPosts,
	}
}

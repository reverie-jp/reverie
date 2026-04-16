package post

import (
	"reverie.jp/reverie/internal/gen/pb/post/v1/postv1connect"
	"reverie.jp/reverie/internal/gen/sqlc"
	postgw "reverie.jp/reverie/internal/modules/post/gateway"
	"reverie.jp/reverie/internal/modules/post/handler"
	"reverie.jp/reverie/internal/modules/post/usecase"
	usergw "reverie.jp/reverie/internal/modules/user/gateway"
)

func InitModule(q sqlc.Querier, userGateway usergw.Gateway) postv1connect.PostServiceHandler {
	postGateway := postgw.New(q, userGateway)
	getPost := usecase.NewGetPost(postGateway)
	createPost := usecase.NewCreatePost(postGateway)
	deletePost := usecase.NewDeletePost(postGateway)
	likePost := usecase.NewLikePost(postGateway)
	unlikePost := usecase.NewUnlikePost(postGateway)
	listUserPosts := usecase.NewListUserPosts(postGateway)
	return handler.New(getPost, createPost, deletePost, likePost, unlikePost, listUserPosts)
}

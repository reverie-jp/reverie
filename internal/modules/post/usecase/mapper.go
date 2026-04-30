package usecase

import postgw "reverie.jp/reverie/internal/modules/post/gateway"

func toPostOutput(view *postgw.PostView) *PostOutput {
	if view == nil {
		return nil
	}

	out := &PostOutput{
		ID:            view.Post.ID,
		ShortID:       view.Post.ShortID,
		Text:          view.Post.Text,
		Author:        view.Author,
		ReplyToPostID:     view.Post.ReplyToPostID,
		RepostPostID:      view.Post.RepostPostID,
		ReplyCount:    view.ReplyCount,
		RepostCount:   view.RepostCount,
		FavoriteCount: view.FavoriteCount,
		IsFavorited:   view.IsFavorited,
		CreateTime:    view.Post.CreateTime,
	}

	if view.RepostOf != nil {
		out.RepostOf = toPostOutput(view.RepostOf)
	}

	return out
}

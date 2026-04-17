package usecase

import postgw "reverie.jp/reverie/internal/modules/post/gateway"

func toPostOutput(view *postgw.PostView) *PostOutput {
	if view == nil {
		return nil
	}

	out := &PostOutput{
		ID:            view.Post.ID,
		Text:          view.Post.Text,
		ReplyToID:     view.Post.ReplyToID,
		RepostID:      view.Post.RepostID,
		ReplyCount:    view.ReplyCount,
		RepostCount:   view.RepostCount,
		FavoriteCount: view.FavoriteCount,
		IsFavorited:   view.IsFavorited,
		CreateTime:    view.Post.CreateTime,
	}

	if view.Author != nil {
		out.Author = &PostAuthorOutput{
			ID:          view.Author.ID,
			CustomID:    view.Author.CustomID,
			DisplayName: view.Author.DisplayName,
			IsPrivate:   view.Author.IsPrivate,
		}
	}

	if view.RepostOf != nil {
		out.RepostOf = toPostOutput(view.RepostOf)
	}

	return out
}

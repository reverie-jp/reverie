package usecase

import (
	"context"

	postgw "reverie.jp/reverie/internal/modules/post/gateway"
	"reverie.jp/reverie/internal/platform/ulid"
	"reverie.jp/reverie/internal/platform/xerrors"
)

type CreatePost struct {
	postGateway postgw.Gateway
}

func NewCreatePost(postGateway postgw.Gateway) *CreatePost {
	return &CreatePost{postGateway: postGateway}
}

func (uc *CreatePost) Execute(ctx context.Context, input CreatePostInput, authorID ulid.ULID) (*PostOutput, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	params := postgw.CreatePostParams{
		ID:       ulid.New(),
		AuthorID: authorID,
		Text:     input.Text,
	}

	if input.ReplyToID != nil {
		id, err := ulid.Parse(*input.ReplyToID)
		if err != nil {
			return nil, xerrors.ErrInvalidArgument.WithMessage("invalid reply_to_id")
		}
		params.ReplyToID = &id
	}

	if input.RepostID != nil {
		id, err := ulid.Parse(*input.RepostID)
		if err != nil {
			return nil, xerrors.ErrInvalidArgument.WithMessage("invalid repost_id")
		}
		params.RepostID = &id
	}

	view, err := uc.postGateway.CreatePost(ctx, params)
	if err != nil {
		return nil, err
	}

	return toPostOutput(view), nil
}

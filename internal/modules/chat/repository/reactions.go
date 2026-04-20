package repository

import (
	"context"

	"reverie.jp/reverie/internal/gen/sqlc"
	"reverie.jp/reverie/internal/platform/ulid"
)

type ReactionRow struct {
	Emoji  string
	Count  int64
	IsMine bool
}

type ReactionByMessageRow struct {
	MessageID ulid.ULID
	Emoji     string
	Count     int64
	IsMine    bool
}

func (r *repositoryImpl) AddReaction(ctx context.Context, messageID, userID ulid.ULID, emoji string) error {
	return r.q.AddMessageReaction(ctx, sqlc.AddMessageReactionParams{
		ID:        ulid.New(),
		MessageID: messageID,
		UserID:    userID,
		Emoji:     emoji,
	})
}

func (r *repositoryImpl) RemoveReaction(ctx context.Context, messageID, userID ulid.ULID, emoji string) error {
	return r.q.RemoveMessageReaction(ctx, sqlc.RemoveMessageReactionParams{
		MessageID: messageID,
		UserID:    userID,
		Emoji:     emoji,
	})
}

func (r *repositoryImpl) ListReactionsByMessage(ctx context.Context, messageID, userID ulid.ULID) ([]*ReactionRow, error) {
	rows, err := r.q.ListReactionsByMessage(ctx, messageID, userID)
	if err != nil {
		return nil, err
	}
	result := make([]*ReactionRow, len(rows))
	for i, row := range rows {
		result[i] = &ReactionRow{Emoji: row.Emoji, Count: row.Count, IsMine: row.IsMine}
	}
	return result, nil
}

func (r *repositoryImpl) ListReactionsByMessages(ctx context.Context, messageIDs []ulid.ULID, userID ulid.ULID) ([]*ReactionByMessageRow, error) {
	ids := make([]string, len(messageIDs))
	for i, id := range messageIDs {
		ids[i] = id.String()
	}
	rows, err := r.q.ListReactionsByMessages(ctx, ids, userID)
	if err != nil {
		return nil, err
	}
	result := make([]*ReactionByMessageRow, len(rows))
	for i, row := range rows {
		result[i] = &ReactionByMessageRow{
			MessageID: row.MessageID,
			Emoji:     row.Emoji,
			Count:     row.Count,
			IsMine:    row.IsMine,
		}
	}
	return result, nil
}

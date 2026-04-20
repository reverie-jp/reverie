package handler

import (
	"context"
	"encoding/base64"

	"connectrpc.com/connect"
	"reverie.jp/reverie/internal/application/server/interceptor"
	chatv1 "reverie.jp/reverie/internal/gen/pb/chat/v1"
	"reverie.jp/reverie/internal/platform/xerrors"
)

func (h *Handler) ListMessages(ctx context.Context, req *connect.Request[chatv1.ListMessagesRequest]) (*connect.Response[chatv1.ListMessagesResponse], error) {
	requestorID, ok := interceptor.UserIDFromContext(ctx)
	if !ok {
		return nil, xerrors.ErrUnauthenticated
	}

	outputs, err := h.listMessages.Execute(ctx, requestorID, req.Msg.RoomId, req.Msg.PageToken, req.Msg.PageSize)
	if err != nil {
		return nil, err
	}

	messages := make([]*chatv1.ChatMessage, len(outputs))
	for i, o := range outputs {
		messages[i] = toProtoMessage(o)
	}

	var nextPageToken string
	if len(outputs) > 0 {
		last := outputs[len(outputs)-1]
		raw := last.CreateTime.UTC().Format("2006-01-02T15:04:05.999999999Z")
		nextPageToken = base64.StdEncoding.EncodeToString([]byte(raw))
	}

	return connect.NewResponse(&chatv1.ListMessagesResponse{
		Messages:      messages,
		NextPageToken: nextPageToken,
	}), nil
}

func (h *Handler) SendMessage(ctx context.Context, req *connect.Request[chatv1.SendMessageRequest]) (*connect.Response[chatv1.SendMessageResponse], error) {
	requestorID, ok := interceptor.UserIDFromContext(ctx)
	if !ok {
		return nil, xerrors.ErrUnauthenticated
	}

	out, err := h.sendMessage.Execute(ctx, requestorID, req.Msg.RoomId, req.Msg.Content)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&chatv1.SendMessageResponse{
		Message: toProtoMessage(out),
	}), nil
}

func (h *Handler) AddMessageReaction(ctx context.Context, req *connect.Request[chatv1.AddMessageReactionRequest]) (*connect.Response[chatv1.AddMessageReactionResponse], error) {
	requestorID, ok := interceptor.UserIDFromContext(ctx)
	if !ok {
		return nil, xerrors.ErrUnauthenticated
	}

	out, err := h.addMessageReaction.Execute(ctx, requestorID, req.Msg.MessageId, req.Msg.Emoji)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&chatv1.AddMessageReactionResponse{
		Message: toProtoMessage(out),
	}), nil
}

func (h *Handler) RemoveMessageReaction(ctx context.Context, req *connect.Request[chatv1.RemoveMessageReactionRequest]) (*connect.Response[chatv1.RemoveMessageReactionResponse], error) {
	requestorID, ok := interceptor.UserIDFromContext(ctx)
	if !ok {
		return nil, xerrors.ErrUnauthenticated
	}

	out, err := h.removeMessageReaction.Execute(ctx, requestorID, req.Msg.MessageId, req.Msg.Emoji)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&chatv1.RemoveMessageReactionResponse{
		Message: toProtoMessage(out),
	}), nil
}

package api

import "context"

type ChatRequest struct {
	Message string `json:"message"`
}

type ChatResponse struct {
	Reply string `json:"reply"`
}

func (c *CoreClient) SendChatMessage(ctx context.Context, s *Session, message string) (string, error) {
	body := ChatRequest{Message: message}
	var resp ChatResponse
	req := &request{method: "POST", path: "/chat", body: body}
	if err := c.doWithAuth(ctx, s, req, &resp); err != nil {
		return "", err
	}
	return resp.Reply, nil
}

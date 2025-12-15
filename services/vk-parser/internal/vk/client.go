package vk

import (
	"fmt"

	"github.com/SevereCloud/vksdk/v3/api"
	"github.com/SevereCloud/vksdk/v3/object"
)

type Client struct {
	vk *api.VK
	groupID int
}

func NewClient(token string, groupID int) *Client {
	vk := api.NewVK(token)
	return &Client{vk: vk, groupID: groupID}
}

func (c *Client) GetWallPosts(count int) ([]object.WallWallpost, error) {
	params := api.Params{
		"owner_id": -c.groupID,
		"count": count,
		"filter": "owner",
		"extended": 0,
	}

	resp, err := c.vk.WallGet(params)
	if err != nil {
		return nil, err
	}

	return resp.Items, nil
}

func (c *Client) GetPostAttachments(postID int) ([]object.WallWallpostAttachment, error) {
	params := api.Params{
		"posts": []string{fmt.Sprintf("%d_%d", -c.groupID, postID)},
		"extended": 1,
	}

	resp, err := c.vk.WallGetByID(params)
	if err != nil {
		return nil, err
	}

	return resp.Items[0].Attachments, nil
}

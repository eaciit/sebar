package sebarclient

type Client struct {
	Host, UserID, Secret string
}

func (c *Client) Connect() error {
	return nil
}

func (c *Client) Close() {
}

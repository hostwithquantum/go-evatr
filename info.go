package evatr

import "context"

// GetStatusMessages returns all status message descriptions.
func (c *Client) GetStatusMessages(ctx context.Context) ([]StatusMessage, error) {
	var result []StatusMessage
	if err := c.doRequest(ctx, "GET", "/v1/info/statusmeldungen", nil, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// GetEUMemberStates returns EU member states and their VIES availability.
func (c *Client) GetEUMemberStates(ctx context.Context) ([]EUMemberState, error) {
	var result []EUMemberState
	if err := c.doRequest(ctx, "GET", "/v1/info/eu_mitgliedstaaten", nil, &result); err != nil {
		return nil, err
	}
	return result, nil
}

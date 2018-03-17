package adn

func (c *Client) streamStatus(token string, uid UID, status bool) (err error) {
	var anos []Annotation
	anos, err = c.GetAnnotations(token, uid)
	if err == nil {
		if len(anos) == 0 {
			anos = append(anos, Annotation{
				Type: StreamAnnotation,
				Value: Stream{
					Online: status,
				},
			})
		} else {
			for idx := range anos {
				if anos[idx].Type == StreamAnnotation {
					anos[idx].Value = Stream{
						Online: status,
					}
				}
			}
		}
		err = c.putAnnotations(token, anos)
	}
	return
}

func (c *Client) StreamOnline(token string, uid UID) (err error) {
	err = c.streamStatus(token, uid, true)
	return
}

func (c *Client) StreamOffline(token string, uid UID) (err error) {
	err = c.streamStatus(token, uid, false)
	return
}

func (c *Client) EnsureStreamChat(token string) (chnl *Channel, err error) {
	chnl, err = c.CreateChannel(token, Channel{
		Type: StreamAnnotation,
		Writers: ChannelMembers{
			AnyUser: true,
		},
		Readers: ChannelMembers{
			AnyUser: true,
		},
	})
	return
}

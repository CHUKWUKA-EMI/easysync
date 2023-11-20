package updatechannel

import channel "github.com/chukwuka-emi/easysync/Channel"

type request struct {
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Type        channel.Type `json:"type"`
}

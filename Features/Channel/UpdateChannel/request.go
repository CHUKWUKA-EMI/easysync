package updatechannel

import channel "github.com/chukwuka-emi/easysync/Features/Channel"

type request struct {
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Type        channel.Type `json:"type"`
}

package exchanges

import client "github.com/lerenn/cryptellation/clients/go"

func (m ServiceInfoResponseMessage) ToModel() client.ServiceInfo {
	return client.ServiceInfo{
		APIVersion: m.Payload.ApiVersion,
		BinVersion: m.Payload.BinVersion,
	}
}

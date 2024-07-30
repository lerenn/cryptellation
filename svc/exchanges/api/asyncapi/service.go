package asyncapi

import "cryptellation/pkg/client"

func (m ServiceInfoResponseMessage) ToModel() client.ServiceInfo {
	return client.ServiceInfo{
		APIVersion: m.Payload.ApiVersion,
		BinVersion: m.Payload.BinVersion,
	}
}

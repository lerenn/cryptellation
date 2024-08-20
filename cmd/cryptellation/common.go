package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/lerenn/cryptellation/pkg/client"
	"github.com/lerenn/cryptellation/pkg/utils"
)

type serviceInformator interface {
	ServiceInfo(ctx context.Context) (client.ServiceInfo, error)
}

func displayServiceInfo(svcClient serviceInformator) error {
	var info client.ServiceInfo
	elapsed, err := utils.ElapsedTime(func() (err error) {
		info, err = svcClient.ServiceInfo(context.TODO())
		return err
	})
	if err != nil {
		return err
	}

	b, err := json.MarshalIndent(info, "", " ")
	if err != nil {
		return err
	}

	fmt.Printf("Infos: %+v\n", string(b))
	fmt.Printf("Elapsed: %q\n", elapsed)
	return nil
}

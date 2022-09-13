package consulapi

import (
	"github.com/hashicorp/consul/api"
	"github.com/pnocera/novaco/internal/settings"
)

// ConsulAPI is a wrapper for the consul api
var sets = settings.GetSettings()
var config *api.Config = api.DefaultConfig()
var client *api.Client

func Init() {

	var err error = nil

	config.Address = sets.GetConsulAddress()

	client, err = api.NewClient(config)

	if err != nil {

		sets.Logger.Error("error creating consul client", err)

	}

}

// GetKV returns the value of a key from the consul kv store

func GetKV(key string) (string, error) {

	kv := client.KV()

	pair, _, err := kv.Get(key, nil)

	if err != nil {

		return "", err

	}

	if pair == nil {

		return "", nil

	}

	return string(pair.Value), nil

}

// PutKV puts a key value pair in the consul kv store

func PutKV(key, value string) error {

	kv := client.KV()

	p := &api.KVPair{Key: key, Value: []byte(value)}

	_, err := kv.Put(p, nil)

	return err

}

// DeleteKV deletes a key from the consul kv store

func DeleteKV(key string) error {

	kv := client.KV()

	_, err := kv.Delete(key, nil)

	return err

}

// GetService returns the address of a service

func GetService(service string) (string, error) {

	services, _, err := client.Catalog().Service(service, "", nil)

	if err != nil {

		return "", err

	}

	if len(services) == 0 {

		return "", nil

	}

	return services[0].ServiceAddress, nil

}

// GetServicePort returns the port of a service

func GetServicePort(service string) (int, error) {

	services, _, err := client.Catalog().Service(service, "", nil)

	if err != nil {

		return 0, err

	}

	if len(services) == 0 {

		return 0, nil

	}

	return services[0].ServicePort, nil

}

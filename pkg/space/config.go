package space

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	clientv3 "go.etcd.io/etcd/client/v3"
	"strings"
)

type AppConf interface {
	ErrorDomain() string
	Dynamic() (bool, []string, string)
}

type ConfigLoader struct {
	conf AppConf
	Env  string
}

func NewConfigLoader(conf AppConf) *ConfigLoader {
	cl := &ConfigLoader{conf: conf}
	cl.loadFromFile()
	cl.loadFromEtcd()
	errDomain = cl.conf.ErrorDomain()
	return cl
}

func (cl *ConfigLoader) loadFromFile() {
	viper.AddConfigPath("./")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath("../../configs")
	viper.AddConfigPath("../../../configs")
	viper.SetConfigType("toml")

	viper.BindEnv("env", "APP_ENV")
	env := viper.GetString("env")
	if strings.TrimSpace(env) == "" {
		env = "dev"
	}
	cl.Env = env
	viper.SetConfigName(fmt.Sprintf("app.%s", env))

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = viper.Unmarshal(cl.conf)
	if err != nil {
		panic(err)
	}
	fmt.Println("Loaded config file:", viper.ConfigFileUsed())
}

func (cl *ConfigLoader) loadFromEtcd() {
	isDynamic, endPoints, key := cl.conf.Dynamic()
	if !isDynamic {
		return
	}

	cli, err := clientv3.New(clientv3.Config{
		Endpoints: endPoints,
	})
	if err != nil {
		panic(err)
	}

	response, err := cli.Get(context.Background(), key)
	if err != nil {
		panic(err)
	}
	for _, kv := range response.Kvs {
		err = json.Unmarshal(kv.Value, cl.conf)
		if err != nil {
			panic(err)
		}
	}

	go func() {
		defer cli.Close()
		rch := cli.Watch(context.Background(), key)
		for watchResp := range rch {
			for _, ev := range watchResp.Events {
				err = json.Unmarshal(ev.Kv.Value, cl.conf)
				if err != nil {
					fmt.Println("json err", err)
				}
			}
		}
	}()
}

package space

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"path"
	"strings"

	"github.com/spf13/viper"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type AppConf interface {
}

type injectConf struct {
	Meta struct {
		AppCode string
	}
	Server struct {
		Http struct {
			Addr string
		}
		Grpc struct {
			Addr string
		}
	}
	Data struct {
		Database struct {
			Host     string
			Port     int
			User     string
			Password string
			Dbname   string
		}
		Redis struct {
			Addr string
			Db   int
		}
	}
	Oss struct {
		Endpoint string
		Key      string
		Secret   string
	}
	Etcd struct {
		EnableDynamicConf bool
		Endpoint          string
	}
	Mq struct {
		Addr string
	}
}

type ConfigLoader struct {
	appConf    AppConf
	injectConf injectConf
	Env        string
}

func NewConfigLoader(conf AppConf) *ConfigLoader {
	cl := &ConfigLoader{appConf: conf}
	cl.loadFromFile()
	cl.loadFromEtcd()
	errDomain = cl.injectConf.Meta.AppCode
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
	err = viper.Unmarshal(cl.appConf)
	if err != nil {
		panic(err)
	}
	err = viper.Unmarshal(&cl.injectConf)
	if err != nil {
		panic(err)
	}
	slog.Info("Loaded config file:" + viper.ConfigFileUsed())
}

func (cl *ConfigLoader) loadFromEtcd() {
	if !cl.injectConf.Etcd.EnableDynamicConf {
		return
	}

	cli, err := clientv3.New(clientv3.Config{
		Endpoints: []string{cl.injectConf.Etcd.Endpoint},
	})
	if err != nil {
		panic(err)
	}

	key := path.Join("/space-appConf", cl.Env, cl.injectConf.Meta.AppCode, "app.json")
	response, err := cli.Get(context.Background(), strings.ToLower(key))
	if err != nil {
		panic(err)
	}
	for _, kv := range response.Kvs {
		err = json.Unmarshal(kv.Value, cl.appConf)
		if err != nil {
			panic(err)
		}
	}

	go func() {
		defer cli.Close()
		rch := cli.Watch(context.Background(), key)
		for watchResp := range rch {
			for _, ev := range watchResp.Events {
				err = json.Unmarshal(ev.Kv.Value, cl.appConf)
				if err != nil {
					fmt.Println("json err", err)
				}
			}
		}
	}()
}

package space

import (
	"context"
	"encoding/json"
	"fmt"
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
		Addr    string
		GroupId string
	}
}

type ConfigLoader struct {
	logger     *Logger
	appConf    AppConf
	injectConf injectConf
	Env        string
}

func NewConfigLoader(conf AppConf, logger *Logger) *ConfigLoader {
	cl := &ConfigLoader{appConf: conf, logger: logger}
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
	cl.logger.Info("Loaded config file: " + viper.ConfigFileUsed())
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

	key := strings.ToLower(path.Join("/space-conf", cl.Env, cl.injectConf.Meta.AppCode, "app.json"))
	cl.logger.Info("Watching etcd key: " + key)
	response, err := cli.Get(context.Background(), key)
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
				cl.logger.Info(fmt.Sprintf("etcd changed, key: %s, value: %s", ev.Kv.Key, ev.Kv.Value))
				err = json.Unmarshal(ev.Kv.Value, cl.appConf)
				if err != nil {
					fmt.Println("json err", err)
				}
			}
		}
	}()
}

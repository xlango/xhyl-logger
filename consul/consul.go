package consul

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"log"
	"logconnection/conf"
	"net/http"
)

func NewConsulClient(addr string) (*api.Client, error) {
	config := api.DefaultConfig()
	config.Address = addr
	return api.NewClient(config)
}

func RegisterServer() {
	if conf.GlobalConfig.ConsulAddress == "" {
		return
	}
	client, err := NewConsulClient(conf.GlobalConfig.ConsulAddress)
	if err != nil {
		log.Fatal("consul client error : ", err)
	}

	registration := new(api.AgentServiceRegistration)
	registration.ID = conf.GlobalConfig.ConsulRegisterId               // 服务节点的名称
	registration.Name = conf.GlobalConfig.ConsulRegisterName           // 服务名称
	registration.Port = conf.GlobalConfig.ConsulRegisterPort           // 服务端口
	registration.Tags = []string{conf.GlobalConfig.ConsulRegisterTags} // tag，可以为空
	registration.Address = conf.GlobalConfig.ConsulRegisterAddress     // 服务 IP

	checkPort := conf.GlobalConfig.Port
	if conf.GlobalConfig.ConsulCheckPort != 0 {
		checkPort = conf.GlobalConfig.ConsulCheckPort
	}

	registration.Check = &api.AgentServiceCheck{ // 健康检查
		HTTP:                           fmt.Sprintf("http://%s:%d%s", registration.Address, checkPort, "/check"),
		Timeout:                        conf.GlobalConfig.ConsulCheckTimeout,
		Interval:                       conf.GlobalConfig.ConsulCheckInterval,            // 健康检查间隔
		DeregisterCriticalServiceAfter: conf.GlobalConfig.DeregisterCriticalServiceAfter, //check失败后30秒删除本服务，注销时间，相当于过期时间
		// GRPC:     fmt.Sprintf("%v:%v/%v", IP, r.Port, r.Service),// grpc 支持，执行健康检查的地址，service 会传到 Health.Check 函数中
	}

	err = client.Agent().ServiceRegister(registration)
	if err != nil {
		log.Fatal("register server error : ", err)
	}

	http.HandleFunc("/check", ConsulCheck)
	http.ListenAndServe(fmt.Sprintf(":%d", checkPort), nil)
}

var count int64

// consul 服务端会自己发送请求，来进行健康检查
func ConsulCheck(w http.ResponseWriter, r *http.Request) {

	s := "consulCheck" + fmt.Sprint(count) + "remote:" + r.RemoteAddr + " " + r.URL.String()
	fmt.Println(s)
	fmt.Fprintln(w, s)
	count++
}

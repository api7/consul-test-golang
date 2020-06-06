package main

import (
    "github.com/gin-gonic/gin"

    "fmt"
    consulapi "github.com/hashicorp/consul/api"
    "log"
    "net/http"
)

func main() {
    r := gin.Default()

    //健康检查
    r.GET("/healthz", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "message": "consul test healthz",
        })
    })

    //注册服务到consul
    ConsulRegister()

    http.ListenAndServe(":8080", r)
}

func ConsulRegister() {
    // 创建连接consul服务配置
    config := consulapi.DefaultConfig()
    config.Address = "127.0.0.1:8500"
    client, err := consulapi.NewClient(config)
    if err != nil {
        log.Fatal("consul client error : ", err)
    }

    //服务信息
    registration := new(consulapi.AgentServiceRegistration)
    registration.ID = "1"
    registration.Name = "go-consul-test"
    registration.Port = 8080
    registration.Tags = []string{"go-consul-test"}
    registration.Address = "127.0.0.1"

    //健康检查
    check := new(consulapi.AgentServiceCheck)
    check.HTTP = fmt.Sprintf("http://%s:%d/healthz", registration.Address, registration.Port)
    check.Timeout = "5s"
    check.Interval = "5s"
    check.DeregisterCriticalServiceAfter = "30s" // 故障检查失败30s后 consul自动将注册服务删除
    registration.Check = check

    //注册
    err = client.Agent().ServiceRegister(registration)
}

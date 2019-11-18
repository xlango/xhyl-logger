package es

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"gopkg.in/olivere/elastic.v5"
	"logconnection/conf"
	"logconnection/proto/model"
	"math/rand"
	"time"
)

func NewEsClient(addr string) (*elastic.Client, error) {
	client, err := elastic.NewClient(
		elastic.SetURL(addr),
		elastic.SetSniff(false),
	)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func InsertLog(index string, typeName string, logInfo string, level model.Level) error {
	client, err := NewEsClient(conf.GlobalConfig.EsHost)
	if err != nil {
		logs.Error("error : es client : ", err.Error())
		return err
	}

	id := time.Now().Format(fmt.Sprintf("%v%d", "20060102150405", rand.Intn(9999)))
	log, err := json.Marshal(model.LogInfo{
		Level:   level.String(),
		Time:    time.Now().Format("2006-01-02 15:04:05"),
		Content: logInfo,
	})

	_, err = client.Index().
		Index(index).
		Type(typeName).
		BodyJson(string(log)).
		Id(id).
		Do(context.Background())

	if err != nil {
		logs.Error("error : es : ", err.Error())
		return err
	}

	return nil
}

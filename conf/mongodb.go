package conf

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type MongoClient struct {
	Database   string //数据库
	Collection string //集合
}

func (this *MongoClient) Insert(v interface{}) bool {
	mongo, err := mgo.Dial(GlobalConfig.MongoHost) // 建立连接

	defer mongo.Close()

	if err != nil {
		return false
	}

	client := mongo.DB(this.Database).C(this.Collection) //选择数据库和集合

	//插入数据
	cErr := client.Insert(v)

	if cErr != nil {
		return false
	}
	return true
}

func (this *MongoClient) FindOne(arg map[string]interface{}) (rs interface{}) {
	mongo, err := mgo.Dial(GlobalConfig.MongoHost) // 建立连接

	defer mongo.Close()

	if err != nil {
		return false
	}

	client := mongo.DB(this.Database).C(this.Collection) //选择数据库和集合

	//查找id为 s20180907
	cErr := client.Find(arg).One(&rs)

	if cErr != nil {
		return nil
	}

	return rs
}

func (this *MongoClient) FindAll(arg map[string]interface{}, rstype interface{}, pageIndex int, pageSize int) (rs []interface{}) {
	mongo, err := mgo.Dial(GlobalConfig.MongoHost) // 建立连接

	defer mongo.Close()

	if err != nil {
		return nil
	}

	client := mongo.DB(this.Database).C(this.Collection) //选择数据库和集合

	//每次最多输出15条数据
	iter := client.Find(arg).Sort("_id").Skip((pageIndex - 1) * pageSize).Limit(pageSize).Iter()

	for iter.Next(&rstype) {
		rs = append(rs, rstype)
	}

	if err := iter.Close(); err != nil {
		return nil
	}

	return rs
}

func (this *MongoClient) Delete(arg map[string]interface{}) bool {
	mongo, err := mgo.Dial(GlobalConfig.MongoHost) // 建立连接

	defer mongo.Close()

	if err != nil {
		return false
	}

	client := mongo.DB(this.Database).C(this.Collection) //选择数据库和集合

	//只更新一条 批量UpdateAll
	cErr := client.Remove(arg)

	if cErr != nil {
		return false
	}

	return true
}

func (this *MongoClient) Update(arg map[string]interface{}, v map[string]interface{}) bool {
	mongo, err := mgo.Dial(GlobalConfig.MongoHost) // 建立连接

	defer mongo.Close()

	if err != nil {
		return false
	}

	client := mongo.DB(this.Database).C(this.Collection) //选择数据库和集合

	//只更新一条
	cErr := client.Update(arg, bson.M{"$set": v})

	if cErr != nil {
		return false
	}

	return true
}

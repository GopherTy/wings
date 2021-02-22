package db

import (
	"crypto/sha512"
	"encoding/hex"
	"errors"

	_ "github.com/go-sql-driver/mysql" // 模块注册导入数据库驱动
	"github.com/go-xorm/xorm"

	"github.com/gopherty/wings/common/conf"
)

var db *xorm.Engine

// Register 注册器
type Register struct {
}

// Name .
func (Register) Name() string {
	return "Common.DB"
}

// Regist 实现 IRegister 接口，以注册获取初始化好的 db 对象。
func (Register) Regist() (err error) {
	cnf := conf.Instance()

	// 初始化日志对象
	// 检查数据库配置内容是否为空。
	if cnf.DB.Driver == "" || cnf.DB.Source == "" {
		return errors.New("driver or source is empty")
	}

	engine, err := xorm.NewEngine(cnf.DB.Driver, cnf.DB.Source)
	if err != nil {
		return err
	}

	// 设置数据库最大连接数和空闲数
	engine.SetMaxOpenConns(cnf.DB.MaxOpenConns)
	engine.SetMaxIdleConns(cnf.DB.MaxIdleConns)

	// 是否开启 SQL 日志
	engine.ShowSQL(cnf.DB.ShowSQL)

	// 是否开启缓存
	if cnf.DB.Cached != 0 {
		cacher := xorm.NewLRUCacher(xorm.NewMemoryStore(), cnf.DB.Cached)
		engine.SetDefaultCacher(cacher)
	}

	db = engine

	// 创建用户相关的表
	err = createTable(&Administrator{})
	if err != nil {
		return
	}
	// 同步表结构
	db.Sync2(&Administrator{})

	// generate administrator
	ok, err := db.Get(&Administrator{
		User: cnf.User.Name,
	})
	if err != nil {
		return
	}
	if ok {
		return
	}

	var passwd string
	h := sha512.New()
	_, err = h.Write([]byte(cnf.User.Passwd))
	if err != nil {
		return
	}
	passwd = hex.EncodeToString(h.Sum(nil))[:18]
	_, err = db.InsertOne(&Administrator{
		User:     cnf.User.Name,
		Password: passwd,
	})
	return
}

func createTable(beans ...interface{}) (err error) {
	var exists bool
	for _, bean := range beans {
		exists, err = db.IsTableExist(bean)
		if err != nil {
			return
		}
		if !exists {
			err = db.CreateTables(bean)
			if err != nil {
				return
			}
		}
	}
	return
}

// Engine 获取 db 对象
func Engine() *xorm.Engine {
	return db
}

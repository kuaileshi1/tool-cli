// @Title mysql初始化
// @Description 请填写文件描述（需要改）
// @Author shigx 2022/3/24 4:39 下午
package mysql

import (
	"fmt"
	"github.com/pkg/errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var _ Repo = (*dbRepo)(nil)

type Repo interface {
	GetDb() *gorm.DB
	CloseDb() error
}

type dbRepo struct {
	DbConn *gorm.DB
}

// @Description 连接mysql数据库
// @Auth shigx
// @Date 2022/3/24 4:57 下午
// @param
// @return
func New(dbAddr, dbUser, dbPass, dbName string) (Repo, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=%t&loc=%s",
		dbUser,
		dbPass,
		dbAddr,
		dbName,
		true,
		"Local")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})

	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("[db connection failed] Database name: %s", dbName))
	}
	db.Set("gorm:table_options", "CHARSET=utf8mb4")
	// db = db.Debug()

	return &dbRepo{
		DbConn: db,
	}, nil
}

func (d *dbRepo) GetDb() *gorm.DB {
	return d.DbConn
}

func (d *dbRepo) CloseDb() error {
	sqlDB, err := d.DbConn.DB()
	if err != nil {
		return err
	}

	return sqlDB.Close()
}

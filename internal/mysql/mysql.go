// Package mysql
// @Title mysql初始化
// @Description 请填写文件描述（需要改）
// @Author shigx 2022/3/24 4:39 下午
package mysql

import (
	"database/sql"
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

type Config struct {
	Addr     string // 数据库链接地址
	User     string // 用户名
	Password string // 密码
	DbName   string // 数据库名
}

type dbRepo struct {
	DbConn *gorm.DB
}

// New
//
//	@Description: 连接mysql数据库
//	@Auth shigx 2024-05-14 16:48:43
//	@param config
//	@return Repo
//	@return error
func New(config *Config) (Repo, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=%t&loc=%s",
		config.User,
		config.Password,
		config.Addr,
		config.DbName,
		true,
		"Local")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})

	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("[db connection failed] Database name: %s", config.DbName))
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

// GetTableComment
//
//	@Description: 查询表备注信息
//	@Auth shigx 2024-05-14 17:02:02
//	@param db
//	@param dbName
//	@param tableName
//	@return string
//	@return error
func GetTableComment(db *gorm.DB, dbName string, tableName string) (string, error) {
	// sql := "SELECT `table_comment` FROM `information_schema`.`tables` WHERE"
	var comment string
	if err := db.Table("information_schema.tables").
		Select("table_comment").
		Where("table_schema = ? and table_name = ?", dbName, tableName).
		Take(&comment).Error; err != nil {
		return "", err
	}

	return comment, nil
}

// TableColumn @Description 表字段信息定义
// @Auth shigx
// @Date 2022/3/24 10:42 下午
// @param
// @return
type TableColumn struct {
	OrdinalPosition int64          `gorm:"column:ORDINAL_POSITION"` // 字段顺序
	ColumnName      string         `gorm:"column:COLUMN_NAME"`      // 字段名称
	ColumnType      string         `gorm:"column:COLUMN_TYPE"`      // 字段类型
	DataType        string         `gorm:"column:DATA_TYPE"`        // 数据类型
	ColumnKey       sql.NullString `gorm:"column:COLUMN_KEY"`       // 字段键
	IsNullable      string         `gorm:"column:IS_NULLABLE"`      // 是否允许为空
	Extra           sql.NullString `gorm:"column:EXTRA"`            // 额外信息
	ColumnComment   sql.NullString `gorm:"column:COLUMN_COMMENT"`   // 字段备注
	ColumnDefault   sql.NullString `gorm:"column:COLUMN_DEFAULT"`   // 字段默认值
}

// GetTableColumn
//
//	@Description: 返回表字段信息
//	@Auth shigx 2024-05-14 17:04:42
//	@param db
//	@param dbName
//	@param tableName
//	@return []TableColumn
//	@return error
func GetTableColumn(db *gorm.DB, dbName string, tableName string) ([]TableColumn, error) {
	ret := make([]TableColumn, 0)
	err := db.Table("information_schema.columns").
		Select(`ORDINAL_POSITION`, `COLUMN_NAME`, `COLUMN_TYPE`, `DATA_TYPE`, `COLUMN_KEY`, `IS_NULLABLE`, `EXTRA`, `COLUMN_COMMENT`, `COLUMN_DEFAULT`).
		Where("table_schema = ? and table_name = ?", dbName, tableName).
		Order("ORDINAL_POSITION ASC").
		Find(&ret).
		Error

	return ret, err
}

package model

import (
	"fmt"
	"time"

	"github.com/AnnonaOrg/osenv"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	Self *gorm.DB
}

var DB *Database
var DBIsReady bool

func Init() error {
	DBIsReady = false
	db, err := GetSelfDB()
	if err != nil {
		return err
	}
	DB = &Database{
		Self: db,
	}
	return nil
}

func GetSelfDB() (*gorm.DB, error) {
	return InitSelfDB()
}

func InitSelfDB() (*gorm.DB, error) {
	switch dbtype := osenv.GetServerDbType(); {
	case dbtype == "mysql":
		return openDBMysql()
	case dbtype == "postgres":
		return openDBPostgreSQL()
	default:
		if len(dbtype) > 0 {
			return nil, fmt.Errorf("配置中数据库类型(%s)不支持", dbtype)
		}
		return nil, fmt.Errorf("配置中数据库类型(%s)未配置", dbtype)
	}
}

func openDBMysql() (*gorm.DB, error) { //user, pass, host, port, dbname string
	user := osenv.GetServerDbUsername() //  viper.GetString("db.username")
	pass := osenv.GetServerDbPassword() // viper.GetString("db.password")
	host := osenv.GetServerDbHost()     // viper.GetString("db.host")
	port := osenv.GetServerDbPort()     // viper.GetString("db.port")
	dbname := osenv.GetServerDbName()   // viper.GetString("db.name")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, pass, host, port, dbname)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("Database connection failed:(%s)", dsn)
	} else {
		DBIsReady = true
	}
	//set for db connection
	setupDB(db)
	return db, nil
	// dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	// db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
}

func openDBPostgreSQL() (*gorm.DB, error) { //username, password, host, port, name string
	user := osenv.GetServerDbUsername() //  viper.GetString("db.username")
	pass := osenv.GetServerDbPassword() // viper.GetString("db.password")
	host := osenv.GetServerDbHost()     // viper.GetString("db.host")
	port := osenv.GetServerDbPort()     // viper.GetString("db.port")
	dbname := osenv.GetServerDbName()   // viper.GetString("db.name")
	dsn := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s  sslmode=disable TimeZone=Asia/Shanghai", user, pass, host, port, dbname)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		// log.Errorf(err, "Database connection failed.:(%s)", dsn)
		return nil, fmt.Errorf("Database connection failed:(%s)", dsn)
	} else {
		DBIsReady = true
	}
	//set for db connection
	setupDB(db)
	return db, nil
	// dsn := "user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai"
	// db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

func setupDB(db *gorm.DB) error {
	sqlDB, err := db.DB()
	// sqlDB, err := db.DB()
	if err != nil {
		// log.Errorf(err, "Database connection setupDB failed.")
		return err
	}

	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(2)

	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(50)

	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(time.Hour)
	return nil
}

func Close() error {
	// DBIsReady = false
	sqlDB, err := DB.Self.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

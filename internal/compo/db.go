package compo

import (
	"awesomeProject/config"
	"awesomeProject/pkg/path"
	"context"
	"fmt"
	"github.com/go-redis/redis/extra/redisotel"
	"github.com/go-redis/redis/v8"
	"github.com/sony/sonyflake"
	"go.uber.org/zap"
	"gopkg.in/natefinch/lumberjack.v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type Data struct {
	db  *gorm.DB
	rdb *redis.Client
	sf  *sonyflake.Sonyflake
}

func NewData(logger *zap.Logger, db *gorm.DB, rdb *redis.Client, sf *sonyflake.Sonyflake) (*Data, func(), error) {
	cleanup := func() {
		logger.Info("cleaning up database")
	}
	return &Data{
		db:  db,
		rdb: rdb,
		sf:  sf,
	}, cleanup, nil
}

func NewDB(conf *config.Configuration, gLog *zap.Logger) *gorm.DB {
	if conf.Database.Driver != "postgres" {
		panic(conf.Database.Driver + " driver is not supported")
	}

	var writer io.Writer
	var logMode logger.LogLevel

	// 是否启用日志文件
	if conf.Database.EnableFileLogWriter {
		logFileDir := conf.Log.RootDir
		if !filepath.IsAbs(logFileDir) {
			logFileDir = filepath.Join(path.WorkPath(), logFileDir)
		}
		// 自定义 Writer
		writer = &lumberjack.Logger{
			Filename:   filepath.Join(logFileDir, conf.Database.LogFilename),
			MaxSize:    conf.Log.MaxSize,
			MaxBackups: conf.Log.MaxBackups,
			MaxAge:     conf.Log.MaxAge,
			Compress:   conf.Log.Compress,
		}
	} else {
		// 默认 Writer
		writer = os.Stdout
	}

	switch conf.Database.LogMode {
	case "silent":
		logMode = logger.Silent
	case "error":
		logMode = logger.Error
	case "warn":
		logMode = logger.Warn
	case "info":
		logMode = logger.Info
	default:
		logMode = logger.Info
	}

	newLogger := logger.New(
		log.New(writer, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,                        // 慢查询 SQL 阈值
			Colorful:                  !conf.Database.EnableFileLogWriter, // 禁用彩色打印
			IgnoreRecordNotFoundError: false,                              // 忽略ErrRecordNotFound（记录未找到）错误
			LogLevel:                  logMode,                            // Log lever
		},
	)

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		conf.Database.Host,
		conf.Database.UserName,
		conf.Database.Password,
		conf.Database.Database,
		strconv.Itoa(conf.Database.Port),
	)
	if db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: conf.Database.TablePrefix,
			//SingularTable: true,

		},
		DisableForeignKeyConstraintWhenMigrating: true,      // 禁用自动创建外键约束
		Logger:                                   newLogger, // 使用自定义 Logger
	}); err != nil {
		gLog.Error("failed opening connection to err:", zap.Any("err", err))
		panic("failed to connect database")
	} else {
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(conf.Database.MaxIdleConnections)
		sqlDB.SetMaxOpenConns(conf.Database.MaxOpenConnections)
		return db
	}
}

// NewRedis .
func NewRedis(c *config.Configuration, gLog *zap.Logger) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     c.Redis.Host + ":" + c.Redis.Port,
		Password: c.Redis.Password, // no password set
		DB:       c.Redis.DB,       // use default DB
	})

	client.AddHook(redisotel.TracingHook{})
	if err := client.Ping(context.Background()).Err(); err != nil {
		gLog.Error("redis connect failed, err:", zap.Any("err", err))
		panic("failed to connect redis")
	}

	return client
}

type contextTxKey struct{}

func (d *Data) ExecTx(ctx context.Context, fn func(ctx context.Context) error) error {
	return d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		ctx = context.WithValue(ctx, contextTxKey{}, tx)
		return fn(ctx)
	})
}

func (d *Data) DB(ctx context.Context) *gorm.DB {
	tx, ok := ctx.Value(contextTxKey{}).(*gorm.DB)
	if ok {
		return tx
	}
	return d.db
}

func (d *Data) RDB() *redis.Client {
	return d.rdb
}

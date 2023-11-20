package database

import (
	"fmt"
	"log"

	"github.com/sky-mirror/boot"
	"github.com/urfave/cli/v2"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

var defaultOption ConnectOption

func init() {
	boot.Register(&defaultOption)
}

var _ boot.Beforer = &ConnectOption{}
var _ boot.Afterer = &ConnectOption{}

// ConnectOption defines a generic connect option for all dialects.
type ConnectOption struct {
	Dialect string
	Host    string
	Port    int // optional, if you append port in host, this option is unnecessary.
	DBName  string
	User    string
	Pass    string
	Config  gorm.Config
	Silence bool

	Testing bool
}

// ConnStr generates connection string.
func (opt *ConnectOption) ConnStr() string {
	switch opt.Dialect {
	case "sqlite3":
		return opt.Host
	case "mysql":
		return fmt.Sprintf(
			"%s:%s@(%s)/%s?charset=utf8&parseTime=True&loc=Local",
			opt.User, opt.Pass, opt.Host, opt.DBName)
	case "postgres":
		return fmt.Sprintf("host=%v port=%v user=%v "+
			"dbname=%v password=%v sslmode=disable",
			opt.Host, opt.Port, opt.User, opt.DBName, opt.Pass)
	case "mssql":
		return fmt.Sprintf(
			"server=%v;user id=%v;password=%v;port=%v;database=%v;",
			opt.Host,
			opt.User,
			opt.Pass,
			opt.Port,
			opt.DBName,
		)
	default:
		log.Panicln("bad dialect: " + opt.Dialect)
	}

	return ""
}

// Dialector generates gorm Dialector.
func (opt *ConnectOption) Dialector() gorm.Dialector {
	dsn := opt.ConnStr()
	switch opt.Dialect {
	case "sqlite3":
		return sqlite.Open(dsn)
	case "mysql":
		return mysql.Open(dsn)
	case "postgres":
		return postgres.Open(dsn)
	case "mssql":
		return sqlserver.Open(dsn)
	default:
		log.Panicln("bad dialect: " + opt.Dialect)
	}

	return nil
}

// CliFlags returns cli flag list.
func (opt *ConnectOption) CliFlags() []cli.Flag {
	var flags []cli.Flag
	flags = append(flags, &cli.StringFlag{
		Name:        "db-dialect",
		Usage:       "[mysql|sqlite3|postgres|mssql]",
		EnvVars:     []string{"DB_DIALECT"},
		Required:    true,
		Destination: &opt.Dialect,
	})
	flags = append(flags, &cli.StringFlag{
		Name:        "db-host",
		Usage:       "mysql -> host, sqlite3 -> filepath",
		EnvVars:     []string{"DB_HOST"},
		Destination: &opt.Host,
	})
	flags = append(flags, &cli.IntFlag{
		Name:        "db-port",
		EnvVars:     []string{"DB_PORT"},
		Destination: &opt.Port,
	})
	flags = append(flags, &cli.StringFlag{
		Name:        "db-name",
		EnvVars:     []string{"DB_NAME"},
		Destination: &opt.DBName,
	})
	flags = append(flags, &cli.StringFlag{
		Name:        "db-user",
		EnvVars:     []string{"DB_USER"},
		Destination: &opt.User,
	})
	flags = append(flags, &cli.StringFlag{
		Name:        "db-password",
		EnvVars:     []string{"DB_PASSWORD"},
		Destination: &opt.Pass,
	})
	flags = append(flags, &cli.BoolFlag{
		Name:        "db-silence-logger",
		EnvVars:     []string{"DB_SILENCE_LOGGER"},
		Destination: &opt.Silence,
	})

	return flags
}

// Before checks config and initializes.
func (opt *ConnectOption) Before(c *cli.Context) error {
	switch opt.Dialect {
	case "sqlite3":
	case "mysql", "postgres", "mssql":
		if opt.Host == "" {
			panic("bad conn setting, host")
		}

		if opt.DBName == "" {
			panic("bad conn setting, name")
		}

		if opt.User == "" {
			panic("bad conn setting, user")
		}

		if opt.Pass == "" {
			panic("bad conn setting, pass")
		}
	default:
		panic("bad dialect value:" + opt.Dialect)
	}

	Initialize(Default, *opt)
	return nil
}

// After closes connection.
func (opt *ConnectOption) After() {
	Finalize()
}

package cliflag

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"
)

var totalOrder int

type AConfig struct {
	Enabled     bool
	Path        string
	Name        string
	beforeOrder int
	afterOrder  int
}

func (c *AConfig) CliFlags() (flags []cli.Flag) {
	flags = append(flags, &cli.BoolFlag{
		Name:        "a-enabled",
		Destination: &c.Enabled,
	})
	flags = append(flags, &cli.StringFlag{
		Name:        "a-path",
		Destination: &c.Path,
		Required:    true,
	})
	flags = append(flags, &cli.StringFlag{
		Name:        "a-name",
		Destination: &c.Name,
		EnvVars:     []string{"NAME"},
	})
	return
}

func (c *AConfig) Before(ctx *cli.Context) error {
	totalOrder++
	c.beforeOrder = totalOrder
	return nil
}

func (c *AConfig) After() {
	totalOrder++
	c.afterOrder = totalOrder
}

type BConfig struct {
	Enabled bool
	Timeout time.Duration
	Path    string
}

func (c *BConfig) CliFlags() (flags []cli.Flag) {
	flags = append(flags, &cli.BoolFlag{
		Name:        "b-enabled",
		Destination: &c.Enabled,
	})
	flags = append(flags, &cli.DurationFlag{
		Name:        "b-timeout",
		Destination: &c.Timeout,
	})
	flags = append(flags, &cli.StringFlag{
		Name:        "a-path", // to test duplicated name.
		Destination: &c.Path,
		Required:    true,
	})
	return
}

type CConfig struct {
	Name        string
	beforeOrder int
	afterOrder  int
}

func (c *CConfig) CliFlags() (flags []cli.Flag) {
	flags = append(flags, &cli.StringFlag{
		Name:        "c-name",
		Destination: &c.Name,
		EnvVars:     []string{"NAME"},
	})
	return
}

func (c *CConfig) Before(ctx *cli.Context) error {
	totalOrder++
	c.beforeOrder = totalOrder
	return nil
}

func (c *CConfig) After() {
	totalOrder++
	c.afterOrder = totalOrder
}

func TestRequired(t *testing.T) {
	s := assert.New(t)

	var a AConfig
	app := &cli.App{
		Name:  "greet",
		Usage: "say a greeting",
		Action: func(c *cli.Context) error {
			fmt.Println("Greetings")
			return nil
		},
		Flags: a.CliFlags(),
	}

	s.NoError(app.Run([]string{"", "-a-enabled=true", "-a-path=/tmp"}))
	s.Error(app.Run([]string{""}))
}

func TestDuplicated(t *testing.T) {
	s := assert.New(t)

	var a AConfig
	var b BConfig
	app := &cli.App{
		Name:  "greet",
		Usage: "say a greeting",
		Action: func(c *cli.Context) error {
			fmt.Println("Greetings")
			return nil
		},
		Flags: append(a.CliFlags(), b.CliFlags()...),
	}

	s.Panics(func() { app.Run([]string{"", "-a-path=/tmp"}) })
	s.Equal("", a.Path)
	s.Equal("", b.Path)
}

func TestOrderingAndSameEnvVars(t *testing.T) {
	s := assert.New(t)

	var a AConfig
	var c CConfig

	Register(&a)
	Register(&c)

	app := &cli.App{
		Name:  "greet",
		Usage: "say a greeting",
		Action: func(c *cli.Context) error {
			fmt.Println("Greetings")
			return nil
		},
		Flags:  Globals(),
		Before: Initialize,
		After:  Finalize,
	}

	os.Setenv("NAME", "reindeer")
	s.NoError(app.Run([]string{"", "-a-path=/tmp"}))
	s.Equal(1, a.beforeOrder)
	s.Equal(2, c.beforeOrder)
	s.Equal(3, c.afterOrder)
	s.Equal(4, a.afterOrder)
	s.Equal("reindeer", a.Name)
	s.Equal("reindeer", c.Name)
}

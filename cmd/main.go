package main

import (
	"context"
	"github.com/mylxsw/go-toolkit/container"
	"github.com/mylxsw/go-toolkit/log"
	"github.com/web-skeleton/framework/internal"
	"gopkg.in/urfave/cli.v1"
	"gopkg.in/urfave/cli.v1/altsrc"
	"io/ioutil"
	"os"
	"strings"
)

var logger = log.Module("main")

var Version = ""

type config struct {
	Data     internal.Data
	Skeleton string
	Output   string
}

func main() {
	serverFlags := []cli.Flag{
		cli.StringFlag{
			Name:  "conf",
			Value: "",
			Usage: "配置文件路径",
		},
		altsrc.NewStringFlag(cli.StringFlag{
			Name:  "log_level",
			Value: "DEBUG",
			Usage: "日志输出级别",
		}),
		altsrc.NewBoolTFlag(cli.BoolTFlag{
			Name:  "log_colorful",
			Usage: "日志是否彩色输出",
		}),
		altsrc.NewStringFlag(cli.StringFlag{
			Name:  "skeleton",
			Value: "./skeleton",
			Usage: "项目骨架模板目录",
		}),
		altsrc.NewStringFlag(cli.StringFlag{
			Name:  "output",
			Value: "project",
			Usage: "输出路径，默认为项目名称",
		}),
		altsrc.NewStringFlag(cli.StringFlag{
			Name:  "template_vars",
			Value: "./vars.json",
			Usage: "模板变量文件",
		}),
	}

	app := cli.NewApp()
	app.Name = ""
	app.Version = Version
	app.Authors = []cli.Author{
		{
			Name:  "mylxsw",
			Email: "mylxsw@aicode.cc",
		},
	}
	app.Action = handler
	app.Before = func(c *cli.Context) error {
		conf := c.String("conf")
		if conf == "" {
			return nil
		}

		inputSource, err := altsrc.NewYamlSourceFromFile(conf)
		if err != nil {
			return err
		}

		return altsrc.ApplyInputSourceValues(c, inputSource, serverFlags)
	}
	app.Flags = serverFlags

	if err := app.Run(os.Args); err != nil {
		logger.Emergency(err.Error())
	}
}

func handler(c *cli.Context) error {
	log.SetDefaultLevel(log.GetLevelByName(c.String("log_level")))
	log.SetDefaultColorful(c.Bool("log_colorful"))

	logger.Infof("version=%s", Version)

	ctx, cancel := context.WithCancel(context.Background())
	cc := container.NewWithContext(ctx)

	cc.MustBindValue("version", Version)
	cc.MustSingleton(func() *cli.Context {
		return c
	})

	defer cc.MustResolve(func() {
		cancel()
	})

	// init configuration
	cc.MustSingleton(func() *config {
		conf := config{}
		conf.Output = c.String("output")
		conf.Skeleton = strings.TrimRight(c.String("skeleton"), "/")

		source, err := ioutil.ReadFile(c.String("template_vars"))
		if err != nil {
			logger.Errorf("open template_vars file failed: %s", err)
			os.Exit(2)
		}

		conf.Data, err = internal.NewData(source)
		if err != nil {
			logger.Errorf("parse template_vars failed: %s", err)
			os.Exit(2)
		}

		return &conf
	})

	return cc.ResolveWithError(func(conf *config) error {
		return internal.Artisan(cc, conf.Skeleton, conf.Output + ".zip", conf.Data)
	})
}
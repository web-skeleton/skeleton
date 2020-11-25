package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/mylxsw/asteria/level"
	"github.com/mylxsw/asteria/log"
	"github.com/mylxsw/container"
	"github.com/web-skeleton/skeleton/internal"
	"gopkg.in/urfave/cli.v1"
	"gopkg.in/urfave/cli.v1/altsrc"
	"gopkg.in/yaml.v2"
)

var logger = log.Module("main")

var Version = ""
var GitCommit = ""

type config struct {
	Data     internal.Data
	Skeleton string
	Output   string
}

type Instruction struct {
	Vars []InstructionVar `json:"vars" yaml:"vars"`
}

type InstructionVar struct {
	Name    string   `json:"name" yaml:"name"`
	Desc    string   `json:"desc" yaml:"desc"`
	Default string   `json:"default" yaml:"default"`
	Options []string `json:"options,omitempty" yaml:"options,omitempty"`
}

func main() {
	serverFlags := []cli.Flag{
		cli.StringFlag{
			Name:  "conf",
			Value: "",
			Usage: "configuration file",
		},
		altsrc.NewStringFlag(cli.StringFlag{
			Name:  "log_level",
			Value: "DEBUG",
			Usage: "log level",
		}),
		altsrc.NewStringFlag(cli.StringFlag{
			Name:  "skeleton",
			Value: "./skeleton",
			Usage: "项目骨架模板目录，只有.sk扩展名的文件才会被作为模板解析，其它文件直接复制到目标文件",
		}),
		altsrc.NewStringFlag(cli.StringFlag{
			Name:  "output",
			Value: "",
			Usage: "输出路径，默认为项目名称",
		}),
		altsrc.NewStringFlag(cli.StringFlag{
			Name:  "template_vars",
			Value: "",
			Usage: "模板变量文件，留空则自动解析 skeleton 目录下的 skeleton.yaml",
		}),
	}

	app := cli.NewApp()
	app.Name = "skeleton"
	app.Version = fmt.Sprintf("%s (%s)", Version, GitCommit)
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
	log.All().LogLevel(level.GetLevelByName(c.String("log_level")))

	logger.Infof("version=%s", fmt.Sprintf("%s (%s)", Version, GitCommit))

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

		conf.Skeleton = strings.TrimRight(c.String("skeleton"), "/")
		conf.Output = c.String("output")
		if conf.Output == "" {
			conf.Output = filepath.Base(conf.Skeleton)
		}

		var source []byte
		var err error
		if c.String("template_vars") != "" {
			source, err = ioutil.ReadFile(c.String("template_vars"))
			if err != nil {
				logger.Errorf("open template_vars file failed: %s", err)
				os.Exit(2)
			}

			conf.Data, err = internal.NewData(source)
			if err != nil {
				logger.Errorf("parse template_vars failed: %s", err)
				os.Exit(2)
			}
		} else {
			instructionBytes, err := ioutil.ReadFile(filepath.Join(conf.Skeleton, "skeleton.yaml"))
			if err != nil {
				logger.Errorf("can not read skeleton.yaml file: %v", err)
				os.Exit(2)
			}

			var instruction Instruction
			if err := yaml.Unmarshal(instructionBytes, &instruction); err != nil {
				logger.Errorf("parse skeleton.yaml failed: %v", err)
				os.Exit(2)
			}

			qs := make([]*survey.Question, 0)
			for _, q := range instruction.Vars {
				var prompt survey.Prompt
				if len(q.Options) > 0 {
					prompt = &survey.Select{Message: q.Desc, Default: q.Default, Options: q.Options}
				} else {
					prompt = &survey.Input{Message: q.Desc, Default: q.Default}
				}

				qs = append(qs, &survey.Question{
					Name:     q.Name,
					Prompt:   prompt,
					Validate: survey.Required,
				})
			}

			data := make(map[string]interface{})
			if err := survey.Ask(qs, &data); err != nil {
				if err == terminal.InterruptErr {
					logger.Warningf("interrupt received")
					os.Exit(0)
				}

				logger.Errorf("parse question failed: %v", err)
				os.Exit(2)
			}

			conf.Data = data
		}

		return &conf
	})

	return cc.ResolveWithError(func(conf *config) error {
		parsedFiles, err := internal.ParseSkeleton(conf.Skeleton, conf.Data)
		if err != nil {
			return err
		}

		return internal.GenerateZip(parsedFiles, conf.Output+".zip")
	})
}

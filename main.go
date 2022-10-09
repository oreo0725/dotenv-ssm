package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

var (
	svc ssmiface.SSMAPI

	Version = "undefined"
)

func init() {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	svc = ssm.New(sess)
}

func GetParameter(name *string) (*ssm.GetParameterOutput, error) {
	results, err := svc.GetParameter(&ssm.GetParameterInput{
		Name:           name,
		WithDecryption: aws.Bool(true),
	})

	return results, err
}

func main() {
	cli.VersionFlag = &cli.BoolFlag{
		Name:    "print-version",
		Aliases: []string{"v"},
		Usage:   "print only the version",
	}

	app := &cli.App{
		Name: "ssm2dotenv",
		Usage: fmt.Sprintf(`ssm2dotenv is a tool to inject SSM parameters into .env file. Given the ssm parameter key, it will fetch the value and replace the value in .env file.
		`),
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "input", Aliases: []string{"i"}, Usage: "input file path", Required: true},
			&cli.StringFlag{Name: "output", Aliases: []string{"o"}, Usage: "output file path", Required: true},
			&cli.StringFlag{Name: "env", Aliases: []string{"e"}, Usage: "the environment variable to be replaced in the input file"},
		},
		Version:   Version,
		UsageText: "ssm2dotenv [--input] [--output] [arguments...]",
		Action: func(c *cli.Context) error {
			input, output := c.String("input"), c.String("output")
			env := c.String("env")
			// read input file
			envContent, err := os.ReadFile(input)
			if err != nil {
				return errors.Wrap(err, "failed to read input file")
			}
			items, err := GetEnvItems(envContent, env)
			if err != nil {
				return err
			}
			// write output items into env file format
			var envLines []string
			for _, item := range items {
				envLines = append(envLines, fmt.Sprintf("%s=%s", item.Name, item.Value))
			}
			outputContent := []byte(strings.Join(envLines, "\n"))
			return os.WriteFile(output, outputContent, 0644)
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatalf("fail: %v", err)
	} else {
		fmt.Println("finish.")
	}
}

type EnvItem struct {
	Name           string
	Value          string
	OriginValue    string
	IsSSMParameter bool
}

func GetEnvItems(envContent []byte, env string) (map[string]EnvItem, error) {
	items := make(map[string]EnvItem)

	lines := strings.Split(string(envContent), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}

		parts := strings.Split(line, "=")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid line: %s", line)
		}

		name := parts[0]
		value, originValue := parts[1], parts[1]

		isSSMParameter := strings.HasPrefix(originValue, "ssm://")
		if isSSMParameter {
			parameterName := strings.TrimPrefix(originValue, "ssm://")
			if env != "" {
				parameterName = strings.ReplaceAll(parameterName, "${env}", env)
			}
			parameter, err := GetParameter(&parameterName)
			if err != nil {
				return nil, errors.Wrap(err, name)
			}
			value = *parameter.Parameter.Value
		}

		items[name] = EnvItem{
			Name:           name,
			Value:          value,
			OriginValue:    originValue,
			IsSSMParameter: isSSMParameter,
		}
	}

	return items, nil
}

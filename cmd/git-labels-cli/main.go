package main

import (
	"context"
	"flag"
	"io/ioutil"
	"log"
	"os"

	"github.com/bvieira/gitlabels/gitlabels"
)

func main() {
	cfgPath := flag.String("cfg", "cfg.yaml", "config filepath")
	token := flag.String("token", "", "github auth token")
	flag.Parse()

	logger := log.New(os.Stdout, "", 0)

	if *token == "" {
		logger.Fatalf("missing '-token' parameter")
	}

	cfg, cerr := config(*cfgPath)
	if cerr != nil {
		logger.Fatalf("could not load config from path:[%s], error:[%v]", *cfgPath, cerr)
	}

	err := gitlabels.NewService(gitlabels.NewGithubRepository(*token), logger).Exec(context.Background(), cfg)
	if err != nil {
		logger.Fatalf("could not exec gitlabels using config:[%s], error:[%v]", *cfgPath, err)
	}
}

func config(file string) (gitlabels.Config, error) {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return gitlabels.Config{}, err
	}
	return gitlabels.ParseConfig(content)
}

package main

import (
	"log"
	"os"
	"path/filepath"

	cc "github.com/ivanpirog/coloredcobra"
	root "github.com/mrtnhwtt/kittypass/cli"
	"github.com/mrtnhwtt/kittypass/internal/config"
)

func main() {
	conf := config.InitializeConfig()

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	logPath := conf.GetString("log_path")
	dir := filepath.Dir(logPath)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(file)

	cmd := root.NewRootCmd(conf)
	cc.Init(&cc.Config{
		RootCmd:  cmd,
		Headings: cc.HiCyan + cc.Bold + cc.Underline,
		Commands: cc.HiYellow + cc.Bold,
		Example:  cc.Italic,
		ExecName: cc.Bold,
		Flags:    cc.Bold,
	})
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

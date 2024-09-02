package main

import (
	"log"
	"os"

	cc "github.com/ivanpirog/coloredcobra"
	root "github.com/mrtnhwtt/kittypass/cli"
	"github.com/mrtnhwtt/kittypass/internal/config"
)

func main() {
	conf := config.InitializeConfig() // add this to an app struct ?

	homedir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	file, err := os.OpenFile(homedir+"/.kittypass.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666) //TODO: replace with path set in config
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

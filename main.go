package main

import (
	"bufio"
	"flag"
	"os"
	"regexp"
	"runtime"

	log "github.com/Crosse/gosimplelogger"
)

func main() {
	var fonts []string

	var filename = flag.String("fromFile", "", "text file containing fonts to install")
	var debug = flag.Bool("debug", false, "Enable debug logging")
	var dryrun = flag.Bool("dry-run", false, "Don't actually download or install anything")
	flag.Parse()

	if *filename == "" && len(flag.Args()) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	if *debug {
		log.LogLevel = log.LogDebug
	} else {
		log.LogLevel = log.LogInfo
	}

	if *filename != "" {
		fd, err := os.Open(*filename)
		if err != nil {
			log.Fatal(err)
		}
		scanner := bufio.NewScanner(fd)
		for scanner.Scan() {
			line := scanner.Text()
			skip, err := regexp.MatchString(`^(#.*|\s*)?$`, line)
			if err != nil {
				log.Errorf("error reading %s: $v", *filename, err)
				continue
			}

			if !skip {
				fonts = append(fonts, line)
			}
		}
		if err = scanner.Err(); err != nil {
			log.Fatal(err)
		}
	}

	for _, v := range flag.Args() {
		fonts = append(fonts, v)
	}

	for _, v := range fonts {
		if *dryrun {
			log.Infof("Would install font(s) from %v", v)
			continue
		}

		log.Debugf("Installing font from %v", v)
		if err := InstallFont(v); err != nil {
			log.Error(err)
		}
	}

	log.Infof("Installed %v fonts", installedFonts)

	if runtime.GOOS == "windows" {
		log.Info("You will need to logoff and logon before the installed font(s) will be available.")
	}
}

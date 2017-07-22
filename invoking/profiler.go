package main

import (
	"fmt"
	"os"
	"runtime/pprof"

	log "github.com/cihub/seelog"
	FLAGS "github.com/jessevdk/go-flags"
)

// Opts options to invoke the profiler
type Opts struct {
	CPU  bool   `short:"c" long:"cpu" description:"run the cpu profiler"`
	Mem  bool   `short:"m" long:"memory" description:"run the memory profiler"`
	Path string `short:"p" long:"path" description:"the path to write the profile to"`
}

var opts Opts
var parser = FLAGS.NewParser(&opts, FLAGS.Default)

// startCPUProfiling starts the cpu profiler and writes to a file
func startCPUProfiling(path string) error {
	f, err := os.Create(path)
	if err != nil {
		log.Critical(err)
		return err
	}
	pprof.StartCPUProfile(f)
	return nil
}

// startMemProfiling starts the memory profiler
func startMemProfiling(path string) error {
	f, err := os.Create(path)
	if err != nil {
		log.Critical(err)
		return err
	}
	pprof.WriteHeapProfile(f)
	return nil
}

// StartProfiling starts the profiler
func StartProfiling() {
	if opts.CPU {
		log.Info("Starting cpu profile")
		startCPUProfiling(opts.Path)
	} else {
		log.Info("Starting memory profile")
		startMemProfiling(opts.Path)
	}

}

// StopProfiling stops the profiler
func StopProfiling() {
	log.Info("Stoping profiler")
	pprof.StopCPUProfile()
}

func validateFlags(opts *Opts) (err error) {
	if opts.CPU && opts.Mem {
		return fmt.Errorf("you cannot profile both memory and cpu, please choose one")
	}

	if len(opts.Path) > 0 {
		return nil
	}

	if len(opts.Path) == 0 && opts.CPU {
		opts.Path = "./cpu-profile.pprof"
		return nil
	}

	if len(opts.Path) == 0 && opts.Mem {
		opts.Path = "./mem-profile.pprof"
		return nil
	}

	return fmt.Errorf("Must provide a profiler type")
}

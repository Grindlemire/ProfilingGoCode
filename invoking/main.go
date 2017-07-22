package main

import (
	"io"
	"os"
	"syscall"

	"github.com/vrecan/death"
	"github.schq.secious.com/UltraViolet/ProfilingGoCode/logger"

	log "github.com/cihub/seelog"
)

func main() {
	_, err := parser.Parse()
	if err != nil {
		Exit(1)
	}
	logger.SetupLogger("/var/log/persistent/profiler")
	defer log.Flush()

	err = validateFlags(&opts)
	if err != nil {
		log.Critical("Error validating profiler flags: ", err)
		Exit(1)
	}

	var goRoutines []io.Closer
	death := death.NewDeath(syscall.SIGINT, syscall.SIGTERM)
	StartProfiling()
	defer StopProfiling()

	generatedNums := make(chan int, 100)

	generator := NewGenerator(generatedNums)
	generator.Start()

	writer := NewWriter(generatedNums, "./output.txt")
	writer.Start()

	death.WaitForDeath(goRoutines...)
}

// Exit flushes logs and returns an exit status
func Exit(status int) {
	log.Flush()
	os.Exit(status)
}

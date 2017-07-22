package logger

import (
	"fmt"
	"os"

	log "github.com/cihub/seelog"
)

// SetupLogger sets up the logger to log to the correct path.
func SetupLogger(logpath string) {
	log.RegisterReceiver("stderr", &StdErrReceiver{})
	seelogConf := `
<seelog type="asynctimer" asyncinterval="1000000">
   <outputs formatid="all">
           <filter levels="info" formatid="fmtinfo" >
              <rollingfile type="size" filename="` + logpath + `" maxsize="20000000" maxrolls="5" />
           </filter>
             <filter levels="warn" formatid="fmtwarn">
              <rollingfile type="size" filename="` + logpath + `" maxsize="20000000" maxrolls="5" />
            </filter>
            <filter levels="error,critical" formatid="fmterror">
            <custom name="stderr" formatid="test"/>
              <rollingfile type="size" filename="` + logpath + `" maxsize="20000000" maxrolls="5" />
            </filter>
            <!-- <filter levels="debug">
              <rollingfile type="size" filename="` + logpath + `" maxsize="20000000" maxrolls="5" />
            </filter>  -->
    </outputs>
    <formats>
        <format id="fmtinfo" format="%EscM(32)[%Level]%EscM(0) [%Date %Time] [%File] %Msg%n"/>
        <format id="fmterror" format="%EscM(31)[%LEVEL]%EscM(0) [%Date %Time] [%FuncShort @ %File.%Line] %Msg%n"/>
         <format id="fmtwarn" format="%EscM(33)[%LEVEL]%EscM(0) [%Date %Time] [%FuncShort @ %File.%Line] %Msg%n"/>
         <format id="test" format="%Msg%n"/>
        <format id="all" format="%EscM(2)[%LEVEL]%EscM(0) [%Date %Time] [%FuncShort @ %File.%Line] %Msg%n"/>
    </formats>
</seelog>
`
	logger, err := log.LoggerFromConfigAsBytes([]byte(seelogConf))
	if err != nil {
		panic(err)
	}
	log.ReplaceLogger(logger)
}

type StdErrReceiver struct {
}

func (r *StdErrReceiver) ReceiveMessage(msg string, level log.LogLevel, context log.LogContextInterface) error {
	fmt.Fprintf(os.Stderr, msg)
	return nil
}

func (r *StdErrReceiver) AfterParse(initArgs log.CustomReceiverInitArgs) error {
	return nil
}

func (r *StdErrReceiver) Flush() {

}

func (r *StdErrReceiver) Close() error {
	return nil
}

package ffmpeg

import (
	"bytes"
	"errors"
	"github.com/je4/goffmpeg/utils"
	"regexp"
	"strings"
)

// Configuration ...
type Configuration struct {
	FfmpegBin             string
	FfprobeBin            string
	PrefixCommandBin      string   // use ffmpeg on different machine (ssh, wsl, ...)
	PrefixCommandParam    []string // parameter (user, privatekey, ...)
	PrefixCommandCmdAsOne bool     // append ffmpeg command as one parameter
}

/*
initialize configuration to use prefix command
actually it is only supported to append ffmpeg at the end of the command
*/
func PrefixConfigure(ffmpegBin string, ffprobeBin, prefixCommand string, cmdAsParam bool) (Configuration, error) {
	conf, err := Configure()
	if err != nil && (ffmpegBin == "" || ffprobeBin == "") {
		return Configuration{}, err
	}
	conf.FfmpegBin = ffmpegBin
	conf.FfprobeBin = ffprobeBin
	conf.PrefixCommandCmdAsOne = cmdAsParam

	// explode command but take care of double quotes
	r := regexp.MustCompile("'.+'|\".+\"|\\S+")
	parts := r.FindAllString(prefixCommand, -1)
	if len(parts) < 1 {
		return Configuration{}, errors.New("invalid prefix command: " + prefixCommand)
	}
	conf.PrefixCommandBin = parts[0]
	conf.PrefixCommandParam = parts[1:]

	_, err = utils.TestCmd(conf.PrefixCommandBin, conf.FfmpegBin)
	if err != nil {
		return Configuration{}, err
	}

	_, err = utils.TestCmd(conf.PrefixCommandBin, conf.FfprobeBin)
	if err != nil {
		return Configuration{}, err
	}
	return conf, nil
}

// Configure Get and set FFmpeg and FFprobe bin paths
func Configure() (Configuration, error) {
	var outFFmpeg bytes.Buffer
	var outProbe bytes.Buffer

	execFFmpegCommand := utils.GetFFmpegExec()
	execFFprobeCommand := utils.GetFFprobeExec()

	outFFmpeg, err := utils.TestCmd(execFFmpegCommand[0], execFFmpegCommand[1])
	if err != nil {
		return Configuration{}, err
	}

	outProbe, err = utils.TestCmd(execFFprobeCommand[0], execFFprobeCommand[1])
	if err != nil {
		return Configuration{}, err
	}

	ffmpeg := strings.Replace(outFFmpeg.String(), utils.LineSeparator(), "", -1)
	fprobe := strings.Replace(outProbe.String(), utils.LineSeparator(), "", -1)

	cnf := Configuration{ffmpeg, fprobe, "", []string{}, false}
	return cnf, nil
}

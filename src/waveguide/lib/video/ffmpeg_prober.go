package video

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"waveguide/lib/log"
	"waveguide/lib/util"
)

type FFProbe struct {
	Path string
}

func (f *FFProbe) Init() (err error) {
	_, err = os.Stat(f.Path)
	return
}

func (f *FFProbe) VideoNeedsEncoding(ifname string, wanted Info) (needs bool, err error) {
	ext := filepath.Ext(ifname)
	outbuff := new(util.Buffer)
	errbuff := new(util.Buffer)
	var args []string
	args = append(args, "-show_streams")
	args = append(args, "-of")
	args = append(args, "json")
	args = append(args, "-i")
	args = append(args, ifname)

	var probe ProbeInfo

	cmd := exec.Command(f.Path, args...)
	cmd.Stdout = outbuff
	cmd.Stderr = errbuff
	err = cmd.Start()
	if err != nil {
		return
	}
	err = cmd.Wait()
	if err == nil {
		err = json.NewDecoder(outbuff).Decode(&probe)
		if err == nil {
			needs = !wanted.Matches(probe.Streams)
			if wanted.Ext != "" {
				needs = needs || (wanted.Ext != ext)
			}
		} else {
			log.Errorf("failed to decode json: %s", err.Error())
		}
	} else {
		log.Errorf("%s %s failed: %s %s", f.Path, args, outbuff.String(), errbuff.String())
	}
	return
}

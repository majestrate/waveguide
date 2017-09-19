package video

import (
	"os"
	"os/exec"
	"waveguide/lib/log"
)

type FFMPEGEncoder struct {
	Path   string
	Params []string
}

func (enc *FFMPEGEncoder) Init() (err error) {
	_, err = os.Stat(enc.Path)
	return err
}

func (enc *FFMPEGEncoder) EncodeFile(infile, outfile string) error {
	args := []string{"-i", infile}
	args = append(args, enc.Params...)
	args = append(args, outfile)
	cmd := exec.Command(enc.Path, args...)
	err := cmd.Run()
	if err != nil {
		out, _ := cmd.CombinedOutput()
		log.Errorf("%s failed: %s", args, string(out))
	}
	return err
}

package video

import (
	"os"
	"os/exec"
	"waveguide/lib/log"
	"waveguide/lib/util"
)

type FFMPEGEncoder struct {
	Path   string
	Params []string
}

func (enc *FFMPEGEncoder) Init() (err error) {
	_, err = os.Stat(enc.Path)
	return err
}

func (enc *FFMPEGEncoder) Transcode(infile, outfile string) error {
	args := []string{"-i", infile}
	args = append(args, "-c:a")
	args = append(args, "copy")
	args = append(args, "-c:v")
	args = append(args, "copy")
	args = append(args, outfile)
	cmd := exec.Command(enc.Path, args...)
	outbuff := new(util.Buffer)
	errbuff := new(util.Buffer)
	cmd.Stdout = outbuff
	cmd.Stderr = errbuff
	err := cmd.Start()
	if err != nil {
		log.Fatalf("%s failed to exec: %s", enc.Path, err)
	}
	err = cmd.Wait()
	if err != nil {
		log.Errorf("%s %s failed: %s %s", enc.Path, args, outbuff.String(), errbuff.String())
	}
	return err
}

func (enc *FFMPEGEncoder) EncodeFile(infile, outfile string) error {
	args := []string{"-i", infile}
	args = append(args, enc.Params...)
	args = append(args, outfile)
	cmd := exec.Command(enc.Path, args...)
	outbuff := new(util.Buffer)
	errbuff := new(util.Buffer)
	cmd.Stdout = outbuff
	cmd.Stderr = errbuff
	err := cmd.Start()
	if err != nil {
		log.Fatalf("%s failed to exec: %s", enc.Path, err)
	}
	err = cmd.Wait()
	if err != nil {
		log.Errorf("%s %s failed: %s %s", enc.Path, args, outbuff.String(), errbuff.String())
	}
	return err
}

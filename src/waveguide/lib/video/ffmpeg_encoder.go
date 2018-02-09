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

func (enc *FFMPEGEncoder) exec(args ...string) error {
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

func (enc *FFMPEGEncoder) Thumbnail(infile, outfile string) error {
	// -ss 01:23:45 -i input -vframes 1 -q:v 2 output.jpg
	return enc.exec("-ss", "00:00:00", "-i", infile, "-vframes", "1", "-filter:v", "scale=480:-1", "-q:v", "2", outfile)
}

func (enc *FFMPEGEncoder) Transcode(infile, outfile string) error {
	return enc.exec("-i", infile, "-c:a", "copy", "-c:v", "copy", outfile)
}

func (enc *FFMPEGEncoder) EncodeFile(infile, outfile string) error {
	args := []string{"-i", infile}
	args = append(args, enc.Params...)
	args = append(args, outfile)
	return enc.exec(args...)
}

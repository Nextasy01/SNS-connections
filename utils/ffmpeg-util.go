package utils

import (
	"os"
	"os/exec"
	"strings"
)

type FFmpegConverter struct{}

func (ffmpeg *FFmpegConverter) ConvertToReelsFormat(filename string) error {
	cmd := exec.Command("ffmpeg", "-i", filename, "-c:v", "libx264", "-c:a", "aac", "-vf",
		"crop='min(in_w,1920)':'min(in_h,1920*16/9)',scale='if(gt(a,9/16),min(iw,1920),-1)':'if(gt(a,9/16),-1,min(ih,1920)*16/9)',setsar=1",
		"-r", "60", "-g", "60", "-pix_fmt", "yuv420p", "-b:v", "25M", "-movflags", "faststart", strings.Replace(filename, ".", "-converted.", 1))

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

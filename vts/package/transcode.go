package vts

import (
	"bytes"
	"common"
	"fmt"
	"log"
	"os/exec"
)

type Transcoder interface {
	Transcode(task common.MqTask) error
}

type FFMpegTranscoder struct {
}

func (t *FFMpegTranscoder) Transcode(task common.MqTask) error {
	// ffmpeg -i input.mp4 -c:v libx264 -c:a aac -strict experimental -b:a 192k -b:v 400k -f mp4 output.mp4
	newFile := fmt.Sprintf("%s_q.mp4", task.Key)
	cmd := exec.Command("ffmpeg", "-i", task.Key, "-c:v", "libx264", "-c:a", "aac", "-strict", "experimental", "-b:a", "192k", "-b:v", "400k", "-f", "mp4", newFile)
	_ = newFile

	errBuff := bytes.Buffer{}
	cmd.Stderr = &errBuff
	err := cmd.Run()
	if err != nil {
		log.Printf("Error %s\n", errBuff.String())
		return err
	}
	return nil
}

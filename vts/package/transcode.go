package vts

import (
	"bytes"
	"common"
	"fmt"
	"log"
	"os/exec"
	"sync"
)

type Transcoder interface {
	Transcode(task common.MqTask) error
}

type FFMpegTranscoder struct {
}

func (t *FFMpegTranscoder) Transcode(task common.MqTask) error {
	var wg sync.WaitGroup
	for _, resln := range task.Resolutions {
		wg.Add(1)
		go func() error {
			defer wg.Done()
			log.Printf("[Debug] compressing %s with resolution %d\n", task.Key, resln)
			err := t.compress(task.Key, fmt.Sprintf("%s_%d.mp4", task.Key, resln), resln)
			if err != nil {
				return err
			}
			return nil
		}()
	}
	wg.Wait()
	return nil
}

func (t *FFMpegTranscoder) compress(inpFile, outFile string, resolution int) error {
	// ffmpeg -i v.mp4 -vf scale=2048:-2 v2.mp4

	cmd := exec.Command("ffmpeg", "-i", inpFile, "-vf", fmt.Sprintf("scale=%d:-2", resolution), outFile, "-y")

	errBuff := bytes.Buffer{}
	cmd.Stderr = &errBuff
	err := cmd.Run()
	if err != nil {
		log.Printf("[Error] for i: %s o: %s with r: %d: %s\n", inpFile, outFile, resolution, errBuff.String())
		return err
	}
	return nil
}

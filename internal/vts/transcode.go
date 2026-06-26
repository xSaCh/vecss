package vts

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"sync"
	"vecss/internal/domain"
)

type Transcoder interface {
	Transcode(task domain.MqTask, inputFilePath string) ([]string, error)
}

type FFMpegTranscoder struct {
}

func (t *FFMpegTranscoder) Transcode(task domain.MqTask, inputFilePath string) ([]string, error) {
	var wg sync.WaitGroup
	paths := []string{}
	for _, resln := range task.Resolutions {
		wg.Add(1)
		go func() error {
			defer wg.Done()
			log.Printf("[Debug] compressing %s with resolution %d\n", inputFilePath, resln)
			outputFilePath := fmt.Sprintf("%s_%d.mp4", inputFilePath, resln)

			err := t.compress(inputFilePath, outputFilePath, resln)
			if err != nil {
				return err
			}

			paths = append(paths, outputFilePath)
			return nil
		}()
	}
	wg.Wait()
	return paths, nil
}

func (t *FFMpegTranscoder) compress(inpFile, outFile string, resolution int) error {
	// ffmpeg -i v.mp4 -vf scale=2048:-2 v2.mp4

	cmd := exec.Command("ffmpeg", "-i", inpFile, "-vf", fmt.Sprintf("scale=%d:-2", resolution), outFile, "-y")
	log.Printf("[Debug] running command: %s\n", cmd.String())
	errBuff := bytes.Buffer{}
	cmd.Stderr = &errBuff
	err := cmd.Run()
	if err != nil {
		log.Printf("[Error] for i: %s o: %s with r: %d: %s\n", inpFile, outFile, resolution, errBuff.String())
		return err
	}
	return nil
}

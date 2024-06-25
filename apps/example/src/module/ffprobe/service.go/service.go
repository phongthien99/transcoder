package service

import (
	"bufio"
	"strings"

	ffmpeg "github.com/u2takey/ffmpeg-go"
	"go.uber.org/zap"
)

type IProbeService interface {
	GetMultimediaInfo(filename string) (string, error)
	FastConvert(input string, output string, config map[string]any) error
}
type probeService struct {
	logger *zap.SugaredLogger
}

// FastConvert implements IProbeService.
func (p *probeService) FastConvert(input string, output string, config map[string]any) error {

	cmd := ffmpeg.Input(input)
	args := ffmpeg.KwArgs{}
	for key, value := range config {
		args[key] = value

	}
	cmd = cmd.Output(output, args)

	cmdSTring := cmd.OverWriteOutput().Compile()

	stderr, err := cmdSTring.StderrPipe()
	if err != nil {
		p.logger.Error("Error creating stderr pipe:", err)
		return err
	}

	if err := cmdSTring.Start(); err != nil {
		p.logger.Error("Error starting ffmpeg command:", err)
		return err
	}

	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			line := scanner.Text()

			if strings.Contains(line, "frame=") && strings.Contains(line, "time=") && strings.Contains(line, "bitrate=") {
				progress := parseProgressLine(line)
				if progress != nil {

					p.logger.Infof("progress: %+v", *progress)
				}
			}
		}

		if err := scanner.Err(); err != nil {
			p.logger.Error("Error reading ffmpeg stdout:", err)
		}
	}()

	if err := cmdSTring.Wait(); err != nil {
		p.logger.Error("Error during ffmpeg command:", err)
		return err
	}

	p.logger.Info("Conversion completed successfully.")
	return nil
}

// GetMultimediaInfo fetches information about a multimedia file using ffprobe.
func (p *probeService) GetMultimediaInfo(filename string) (string, error) {
	// Execute FFprobe command to get information
	info, err := ffmpeg.Probe(filename)
	if err != nil {
		p.logger.Error(err)
		return "", err
	}

	return info, nil

}

func NewProbeService(logger *zap.SugaredLogger) IProbeService {
	return &probeService{
		logger: logger,
	}
}

type Progress struct {
	Progress        float64
	CurrentBitrate  string
	FramesProcessed string
	CurrentTime     string
	Speed           string
}

func parseProgressLine(line string) *Progress {
	progress := &Progress{}

	fields := strings.Fields(line)
	for _, field := range fields {
		if strings.HasPrefix(field, "frame=") {
			progress.FramesProcessed = strings.TrimPrefix(field, "frame=")
		} else if strings.HasPrefix(field, "time=") {
			progress.CurrentTime = strings.TrimPrefix(field, "time=")
		} else if strings.HasPrefix(field, "bitrate=") {
			progress.CurrentBitrate = strings.TrimPrefix(field, "bitrate=")
		} else if strings.HasPrefix(field, "speed=") {
			progress.Speed = strings.TrimPrefix(field, "speed=")

		}
	}

	return progress
}

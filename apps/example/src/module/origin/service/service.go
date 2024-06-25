package service

import (
	"errors"
	"example/src/module/origin/adapter"
	"example/src/module/origin/model"
	"path"
	"sync"

	"github.com/spf13/afero"
	"go.uber.org/zap"
)

type IFileUploadService interface {
	UploadFile(path string, content []byte, base64Config string) error
	ReadFile(path string, base64Config string) ([]byte, error)
}

// Struct cho FileUploadService
type FileUploadService struct {
	Uploaders map[string]afero.Fs
	mu        sync.RWMutex // RWMutex để bảo vệ sự thay đổi của Uploaders map
	logger    *zap.SugaredLogger
}

// NewFileUploadService khởi tạo một instance của FileUploadService với Uploader tương ứng
func NewFileUploadService(logger *zap.SugaredLogger) (IFileUploadService, error) {
	s := &FileUploadService{
		Uploaders: make(map[string]afero.Fs),
		logger:    logger,
	}

	return s, nil
}

// GetUploader trả về một uploader từ map Uploaders
func (s *FileUploadService) GetUploader(backend string) (afero.Fs, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	uploader, exists := s.Uploaders[backend]
	return uploader, exists
}

// SetUploader thiết lập một uploader mới vào map Uploaders
func (s *FileUploadService) SetUploader(backend string, uploader afero.Fs) afero.Fs {

	s.Uploaders[backend] = uploader
	return uploader
}

// initUploaderIfNeeded kiểm tra và khởi tạo uploader nếu không tồn tại
func (s *FileUploadService) initUploaderIfNeeded(key string, config model.Config) (afero.Fs, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.Uploaders[key]; !exists {
		switch config.FS {
		case model.S3:
			fs, err := adapter.NewS3Fs(config.S3Config)
			if err != nil {
				s.logger.Error(err)
				return nil, err
			}
			s.SetUploader(key, fs)
			return fs, nil
		case model.FTP:
			fs, err := adapter.NewSftps(config.FTPConfig)
			if err != nil {
				s.logger.Error(err)
				return nil, err
			}
			s.SetUploader(key, fs)
		// case model.HTTP:
		// 	s.SetUploader(key, NewHTTPUploader(config.HTTPConfig))
		case model.OS:
			fs, err := adapter.NewOSFs(config.OSConfig.BasePath)
			if err != nil {
				s.logger.Error(err)
				return nil, err
			}
			s.SetUploader(key, fs)
			return fs, nil
			// default:
			// 	fmt.Printf("unsupported fs: %s", config.FS)
		}
	}
	return nil, errors.New("Not init uploader")
}

// UploadFile uploads a file to the configured backend
func (s *FileUploadService) UploadFile(filePath string, content []byte, base64Config string) error {
	var uploader afero.Fs
	uploader, exists := s.GetUploader(base64Config)
	if !exists {
		config, err := model.ParseBase64Config(base64Config) // Thay Config{} bằng config thực tế của bạn nếu có
		if err != nil {
			s.logger.Error(err)
			return err
		}
		uploader, err = s.initUploaderIfNeeded(base64Config, config)
		if err != nil {
			s.logger.Error(err)
			return err
		}
	}
	// Tạo đường dẫn đầy đủ nếu cần
	err := s.ensureDirExists(uploader, filePath)
	if err != nil {
		s.logger.Errorf("Error ensuring directory exists: %s", err)
		return err
	}

	return afero.WriteFile(uploader, filePath, content, 0755)
}

// ReadFile reads a file from the configured backend
func (s *FileUploadService) ReadFile(filePath string, base64Config string) ([]byte, error) {

	uploader, exists := s.GetUploader(base64Config)
	if !exists {
		config, err := model.ParseBase64Config(base64Config) // Thay Config{} bằng config thực tế của bạn nếu có
		if err != nil {
			s.logger.Error(err)
			return nil, err
		}
		uploader, err = s.initUploaderIfNeeded(base64Config, config)
		if err != nil {
			s.logger.Error(err)
			return nil, err
		}
	}

	return afero.ReadFile(uploader, filePath)
}

func (s *FileUploadService) ensureDirExists(fs afero.Fs, filePath string) error {
	dir := path.Dir(filePath)

	// Kiểm tra xem thư mục cha đã tồn tại chưa, nếu chưa thì tạo mới
	exist, err := afero.DirExists(fs, dir)
	if err != nil {
		return err
	}

	if !exist {
		err := fs.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}
	}

	return nil
}

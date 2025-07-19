package services

type Container struct {
	FileUpload *FileUploadService
}

var Services *Container

func InitServices(bucketName, region string) error {
	fileUploadService, err := NewFileUploadService(bucketName, region)
	if err != nil {
		return err
	}

	Services = &Container{
		FileUpload: fileUploadService,
	}
	return nil
}
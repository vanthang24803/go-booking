package cloudinary

import (
	"context"
	"mime/multipart"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/may20xx/booking/config"
	"github.com/may20xx/booking/pkg/log"
)

type Cloudinary interface {
	UploadFile(file multipart.File) (*uploader.UploadResult, error)
	DeleteFile(publicID string) (*uploader.DestroyResult, error)
}

type CloudinaryService struct {
	cld *cloudinary.Cloudinary
}

func NewCloudinaryService() (*CloudinaryService, error) {

	settings := config.GetConfig()

	cld, err := cloudinary.NewFromParams(settings.CloudinaryCloudName, settings.CloudinaryAPIKey, settings.CloudinaryAPISecret)

	if err != nil {
		log.Msg.Panic("Cloudinary error: ", err)
		return nil, err
	}
	return &CloudinaryService{cld: cld}, nil
}

func (s *CloudinaryService) UploadFile(file multipart.File) (*uploader.UploadResult, error) {
	ctx := context.Background()
	uploadParams := uploader.UploadParams{
		Format: "webp",
	}

	result, err := s.cld.Upload.Upload(ctx, file, uploadParams)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *CloudinaryService) DeleteFile(publicID string) (*uploader.DestroyResult, error) {
	ctx := context.Background()

	result, err := s.cld.Upload.Destroy(ctx, uploader.DestroyParams{PublicID: publicID})
	if err != nil {
		return nil, err
	}

	return result, nil
}

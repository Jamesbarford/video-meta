package video

type VideoMetaService interface {
	CreateVideoMeta(videoId int, metaData *[]VideoMetaDataPayload) error
	ReadVideoMeta(videoId int) (*[]VideoMetaData, error)
	UpdateVideoMeta(metaId int, metaData *VideoMetaData) (*VideoMetaData, error)
	DeleteVideoMeta(metaDataId int) error
}

type VideoMetaServiceImpl struct {
	repo VideoMetaRepository
}

func NewVideoMetaService(repo VideoMetaRepository) VideoMetaService {
	return &VideoMetaServiceImpl{
		repo: repo,
	}
}

func (service *VideoMetaServiceImpl) CreateVideoMeta(videoId int, metaData *[]VideoMetaDataPayload) error {
	return service.repo.CreateVideoMeta(videoId, metaData)
}

func (service *VideoMetaServiceImpl) ReadVideoMeta(videoId int) (*[]VideoMetaData, error) {
	return service.repo.ReadVideoMeta(videoId)
}

func (service *VideoMetaServiceImpl) UpdateVideoMeta(metaId int, metaData *VideoMetaData) (*VideoMetaData, error) {
	return service.repo.UpdateVideoMeta(metaId, &metaData.Key, &metaData.Value)
}

func (service *VideoMetaServiceImpl) DeleteVideoMeta(metaDataId int) error {
	return service.repo.DeleteVideoMeta(metaDataId)
}

package video

import (
	"fmt"
	"log"
	"strconv"

	"github.com/Jamesbarford/video-meta/server/database"
	"github.com/pkg/errors"
)

type VideoMetaRepository interface {
	CreateVideoMeta(videoId int, metaData *[]VideoMetaDataPayload) error
	ReadVideoMeta(videoId int) (*[]VideoMetaData, error)
	UpdateVideoMeta(metaId int, key *string, value *string) (*VideoMetaData, error)
	DeleteVideoMeta(metaId int) error
}

type VideoMetaRepositoryImpl struct {
	db *database.DbConnection
}

func NewVideoMetaRepository(db *database.DbConnection) VideoMetaRepository {
	return &VideoMetaRepositoryImpl{
		db: db,
	}
}

func (repo *VideoMetaRepositoryImpl) CreateVideoMeta(videoId int, metaData *[]VideoMetaDataPayload) error {
	if len(*metaData) == 0 {
		return errors.New("No values in metadata array")
	}
	query := "INSERT INTO meta_data (video_id, key, value) VALUES "
	values := []interface{}{}
	argc := 1
	strVideoId := strconv.Itoa(videoId)

	for _, r := range *metaData {
		query += fmt.Sprintf("($%d, $%d, $%d),", argc, argc+1, argc+2)
		values = append(values, strVideoId, r.Key, r.Value)
		argc += 3
	}

	query = query[:len(query)-1]

	rows, err := repo.db.Query(query, values...)

	if err != nil {
		log.Printf("Error preparing statement for creating video meta with video id %d: %v", videoId, err)
		return errors.Wrap(err, "Failed to create VideoMetaData")
	}
	defer rows.Close()

	return nil
}

func (repo *VideoMetaRepositoryImpl) ReadVideoMeta(videoId int) (*[]VideoMetaData, error) {
	query := "SELECT key, value, id FROM meta_data WHERE video_id = $1"
	rows, err := repo.db.Query(query, videoId)

	if err != nil {
		log.Printf("Error preparing statement for reading video meta with ID %d: %v", videoId, err)
		return nil, errors.Wrap(err, "failed to prepare statement for reading video meta")
	}

	defer rows.Close()
	var metaData []VideoMetaData
	for rows.Next() {
		var key, value string
		var metaId int
		err := rows.Scan(&key, &value, &metaId)
		if err != nil {
			log.Printf("Error scanning row for video meta with ID %d: %v", videoId, err)
			return nil, errors.Wrap(err, "failed to scan row for video meta")
		}
		metaData = append(metaData, VideoMetaData{Key: key, Value: value, Id: metaId})
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating rows for video meta with ID %d: %v", videoId, err)
		return nil, errors.Wrap(err, "error iterating rows for video meta")
	}

	return &metaData, nil
}

func (repo *VideoMetaRepositoryImpl) UpdateVideoMeta(metaId int, key *string, value *string) (*VideoMetaData, error) {
	query := "UPDATE meta_data SET value = $1, key = $2 WHERE id = $3"

	if _, err := repo.db.Query(query, *value, *key, metaId); err != nil {
		log.Printf("Error updating meta_data for video with ID %d: %v", metaId, err)
		return nil, errors.Wrap(err, "error updating VideoMetaData")
	}

	// Return the updated metadata
	return &VideoMetaData{Key: *key, Value: *value, Id: metaId}, nil
}

func (repo *VideoMetaRepositoryImpl) DeleteVideoMeta(metaId int) error {
	query := "DELETE FROM meta_data WHERE id = $1"
	rows, err := repo.db.Query(query, metaId)

	if err != nil {
		log.Printf("Error deleting meta_data for video with ID %d: %v", metaId, err)
		return errors.Wrap(err, "error deleting VideoMetaData")
	}
	defer rows.Close()

	return nil
}

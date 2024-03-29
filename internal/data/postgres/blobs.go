package postgres

import (
	"github.com/Masterminds/squirrel"
	"github.com/kenjitheman/blobman/internal/data"
	"gitlab.com/distributed_lab/kit/pgdb"
	"math"
)

const blobsTableName = "blobs"

type BlobsDB struct {
	db  *pgdb.DB
	sql squirrel.Sqlizer
}

func NewBlobs(db *pgdb.DB) data.Blobs {
	return &BlobsDB{
		db:  db.Clone(),
		sql: squirrel.Select("*").From(blobsTableName),
	}
}

func (b *BlobsDB) New() data.Blobs {
	return NewBlobs(b.db)
}

func (b *BlobsDB) GetBlobById(id string) (*data.Blob, error) {
	var result data.Blob
	b.sql = squirrel.Select("*").From(blobsTableName).Where("id = ?", id)
	err := b.db.Get(&result, b.sql)
	return &result, err
}

func (b *BlobsDB) CreateBlob(blob *data.Blob, owner data.Owner) error {
	b.sql = squirrel.Insert(blobsTableName).SetMap(map[string]interface{}{
		"id":       blob.ID,
		"owner_id": owner.ID,
		"value":    blob.Value,
	})
	return b.db.Exec(b.sql)
}

func (b *BlobsDB) DeleteBlob(id string) error {
	b.sql = squirrel.Delete(blobsTableName).Where(squirrel.Eq{"id": id})
	return b.db.Exec(b.sql)
}

func (b *BlobsDB) GetTotalPages(limit uint64) (uint64, error) {
	var res []data.Blob
	err := b.db.Select(&res, b.sql.(squirrel.SelectBuilder))
	if err != nil {
		return 0, err
	}
	return uint64(math.Ceil(float64(len(res)) / float64(limit))), nil
}

func (b *BlobsDB) GetPage(pageParams pgdb.OffsetPageParams) data.Blobs {
	b.sql = pageParams.ApplyTo(b.sql.(squirrel.SelectBuilder), "id")
	return b
}

func (b *BlobsDB) GetBlobsList() ([]data.Blob, error) {
	var result []data.Blob
	b.sql = squirrel.Select("*").FromSelect(b.sql.(squirrel.SelectBuilder), "blobs_by_page")
	err := b.db.Select(&result, b.sql)
	return result, err
}

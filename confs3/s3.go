package confs3

import (
	"context"
	"errors"
	"io"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/go-courier/envconf"
	"github.com/minio/minio-go"
)

type S3Endpoint interface {
	Endpoint() string
	AccessKeyID() string
	SecretAccessKey() string
	BucketName() string
	Secure() bool
}

var (
	ErrUnknownContentType = errors.New("unknown content type")
	ErrInvalidObject      = errors.New("invalid object key")
)

type ObjectDB struct {
	Endpoint        string                                                             `env:",upstream"`
	AccessKeyID     string                                                             `env:""`
	SecretAccessKey envconf.Password                                                   `env:""`
	BucketName      string                                                             `env:""`
	Secure          bool                                                               `env:""`
	PresignedValues func(db *ObjectDB, key string, expiresIn time.Duration) url.Values `env:"-"`
}

func (db *ObjectDB) LivenessCheck() map[string]string {
	key := db.BucketName + "." + db.Endpoint
	m := map[string]string{
		key: "ok",
	}

	c, err := db.Client()

	if err != nil {
		m[key] = err.Error()
	} else {
		if _, err := c.GetBucketLocation(db.BucketName); err != nil {
			m[key] = err.Error()
		}
	}

	return m
}

func (db *ObjectDB) Client() (*minio.Client, error) {
	client, err := minio.New(db.Endpoint, db.AccessKeyID, db.SecretAccessKey.String(), db.Secure)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (db *ObjectDB) PublicURL(meta *ObjectMeta) *url.URL {
	u := &url.URL{}
	u.Scheme = "http"
	if db.Secure {
		u.Scheme += "s"
	}

	u.Host = db.Endpoint
	u.Path = db.BucketName + "/" + meta.Key()
	return u
}

func (db *ObjectDB) ProtectURL(meta *ObjectMeta, expiresIn time.Duration) (*url.URL, error) {
	c, err := db.Client()
	if err != nil {
		return nil, err
	}
	values := url.Values{}
	if db.PresignedValues != nil {
		values = db.PresignedValues(db, meta.Key(), expiresIn)
	}

	u, err := c.PresignedGetObject(db.BucketName, meta.Key(), expiresIn, values)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (db *ObjectDB) PutObject(ctx context.Context, fileReader io.Reader, meta ObjectMeta) error {
	if ctx == nil {
		ctx = context.Background()
	}

	c, err := db.Client()
	if err != nil {
		return err
	}

	if meta.Size == 0 {
		if canLen, ok := fileReader.(interface{ Len() int }); ok {
			meta.Size = int64(canLen.Len())
		}
	}

	_, err = c.PutObjectWithContext(ctx, db.BucketName, meta.Key(), fileReader, meta.Size, minio.PutObjectOptions{
		ContentType: meta.ContentType,
	})

	return err
}

func (db *ObjectDB) ReadObject(writer io.Writer, group string, objectID uint64) error {
	c, err := db.Client()
	if err != nil {
		return err
	}

	object, err := c.GetObject(db.BucketName, (&ObjectMeta{Group: group, ObjectID: objectID}).Key(), minio.GetObjectOptions{})
	if err != nil {
		return err
	}
	defer object.Close()

	_, err = io.Copy(writer, object)
	if err != nil {
		return err
	}

	return err
}

func (db *ObjectDB) PresignedPutObject(group string, objectID uint64, expiresIn time.Duration) (string, error) {
	c, err := db.Client()
	if err != nil {
		return "", err
	}
	presignedURL, err := c.PresignedPutObject(db.BucketName, (&ObjectMeta{Group: group, ObjectID: objectID}).Key(), expiresIn)
	if err != nil {
		return "", err
	}
	return presignedURL.String(), nil
}

func (db *ObjectDB) DeleteObject(group string, objectID uint64) error {
	c, err := db.Client()
	if err != nil {
		return err
	}

	return c.RemoveObject(db.BucketName, (&ObjectMeta{Group: group, ObjectID: objectID}).Key())
}

func (db *ObjectDB) StatsObject(group string, objectID uint64) (*ObjectMeta, error) {
	c, err := db.Client()
	if err != nil {
		return nil, err
	}

	object, err := c.GetObject(db.BucketName, (&ObjectMeta{Group: group, ObjectID: objectID}).Key(), minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	defer object.Close()

	info, err := object.Stat()
	if err != nil {
		return nil, err
	}

	om, err := ParseObjectMetaFromKey(info.Key)
	if err != nil {
		return nil, err
	}

	om.ContentType = info.ContentType
	om.ETag = info.ETag
	om.Size = info.Size

	return om, err
}

func (db *ObjectDB) ListObjectByGroup(group string) ([]*ObjectMeta, error) {
	c, err := db.Client()
	if err != nil {
		return nil, err
	}

	chDone := make(chan struct{})
	defer close(chDone)

	metas := make([]*ObjectMeta, 0)

	for obj := range c.ListObjectsV2(db.BucketName, group, true, chDone) {
		om, err := ParseObjectMetaFromKey(obj.Key)
		if err != nil {
			continue
		}

		om.ContentType = obj.ContentType
		om.ETag = obj.ETag
		om.Size = obj.Size

		metas = append(metas, om)
	}

	return metas, nil
}

func ParseObjectMetaFromKey(key string) (*ObjectMeta, error) {
	parts := strings.Split(key, "/")
	if len(parts) != 2 {
		return nil, ErrInvalidObject
	}
	group := parts[0]

	objectID, err := strconv.ParseUint(parts[1], 10, 64)
	if err != nil {
		return nil, ErrInvalidObject
	}

	om := &ObjectMeta{
		ObjectID: objectID,
		Group:    group,
	}

	return om, nil
}

type ObjectMeta struct {
	ObjectID    uint64 `json:"objectID"`
	Group       string `json:"group"`
	Size        int64  `json:"size"`
	ContentType string `json:"contentType"`
	ETag        string `json:"etag"`
}

func (meta ObjectMeta) Key() string {
	return meta.Group + "/" + strconv.FormatUint(meta.ObjectID, 10)
}

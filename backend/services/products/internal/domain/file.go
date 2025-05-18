package domain

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/m11ano/mipt-webdev-course/backend/services/products/internal/infra/config"
	"github.com/m11ano/mipt-webdev-course/backend/services/products/internal/infra/storage"
)

type FileTarget int

const (
	FileTargetProductPreview FileTarget = 1
	FileTargetProductSlider  FileTarget = 2
)

func (s FileTarget) String() string {
	switch s {
	case FileTargetProductPreview:
		return "product_preview"
	case FileTargetProductSlider:
		return "product_slider"
	default:
		return ""
	}
}

var FileTargetMap = map[string]FileTarget{
	"product_preview": FileTargetProductPreview,
	"product_slider":  FileTargetProductSlider,
}

type File struct {
	ID                  uuid.UUID
	Target              FileTarget
	AssignedToTarget    bool
	StorageFileKey      *string
	UploadedToStorage   bool
	ToDeleteFromStorage bool

	CreatedAt time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time
}

func NewFile(target FileTarget, assignedToTarget bool) *File {
	return &File{
		ID:               uuid.New(),
		Target:           target,
		AssignedToTarget: assignedToTarget,
		CreatedAt:        time.Now(),
	}
}

func (i *File) GetURL(cfg *config.Config) string {
	var builder strings.Builder
	builder.WriteString(storage.GetBucketURL(storage.BucketProducts, cfg))
	builder.WriteString("/")

	if i.StorageFileKey != nil {
		builder.WriteString(*i.StorageFileKey)
	}

	return builder.String()
}

package sharedstorage

import (
	"fmt"
	"strings"

	"github.com/Watari995/musclead/internal/myerror"
)

type ImageKind string

const (
	ImageKindProfile ImageKind = "profiles"
	ImageKindMeal    ImageKind = "meals"
)

func ValidateUserOwnedImagePath(kind ImageKind, userID, path string) error {
	prefix := fmt.Sprintf("%s/%s/", kind, userID)
	if !strings.HasPrefix(path, prefix) {
		return myerror.NewBadRequestError().SetMessage("path does not belong to current user")
	}
	if strings.Contains(path, "..") {
		return myerror.NewBadRequestError().SetMessage("path traversal is not allowed")
	}
	return nil
}

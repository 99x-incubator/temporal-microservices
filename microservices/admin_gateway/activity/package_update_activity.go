package activity

import (
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"

	"99x.io/admin_gateway/shared"
)

func UpdatePackageActivity(ctx context.Context, packageId string, userId string) (string, error) {

	reqBody := `{"user_id": "` + userId + `", "package_id": "` + packageId + `"}`
	_, err := http.Post(shared.UPDATE_PACKAGE_SERVICE, "application/json", ioutil.NopCloser(strings.NewReader(reqBody)))
	if err != nil {
		return "", errors.New("failed to update package")
	}

	return packageId, nil
}

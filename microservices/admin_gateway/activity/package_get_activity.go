package activity

import (
	"context"
	"errors"
	"io/ioutil"
	"net/http"

	"99x.io/admin_gateway/shared"
)

func GetPackageActivity(ctx context.Context, userId string) (string, error) {
	resp, err := http.Get(shared.GET_PACKAGE_SERVICE + "?userID=" + userId)
	if err != nil || resp.StatusCode != http.StatusOK {
		return "", errors.New("failed to get current package")
	}
	defer resp.Body.Close()
	currentPackage, _ := ioutil.ReadAll(resp.Body)
	return string(currentPackage), nil
}

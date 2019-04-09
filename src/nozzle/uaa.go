package nozzle

import "github.com/cloudfoundry-incubator/uaago"

//go:generate counterfeiter ./ UAA
type UAA interface {
	GetAuthToken() (string, error)
}

type uaa struct {
	uaaUser           string
	uaaPass           string
	skipSSLValidation bool
	uaaClient         *uaago.Client
}

func NewUAA(uaaURL string, uaaUser string, uaaPass string, skipSSLValidation bool) (UAA, error) {
	uaaClient, err := uaago.NewClient(uaaURL)
	if err != nil {
		return nil, err
	}

	return &uaa{
		uaaUser:           uaaUser,
		uaaPass:           uaaPass,
		skipSSLValidation: skipSSLValidation,
		uaaClient:         uaaClient,
	}, nil
}

func (uaa *uaa) GetAuthToken() (string, error) {
	token, _, err := uaa.uaaClient.GetAuthTokenWithExpiresIn(uaa.uaaUser, uaa.uaaPass, uaa.skipSSLValidation)
	return token, err
}

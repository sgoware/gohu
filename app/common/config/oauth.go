package config

import "errors"

func GetClientSecret(clientId string) (clientSecret string, err error) {
	if agolloClient == nil {
		return "", errors.New(emptyConfigClientErr)
	}
	v, err := agolloClient.GetViper("oauth.yaml")
	if err != nil {
		return "", errors.New(getViperErr)
	}

	clientSecret = v.GetString("Client." + clientId + ".Secret")
	if clientSecret == "" {
		return "", errors.New(viperEmptyKeyErr)
	}
	return clientSecret, nil
}

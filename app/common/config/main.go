package config

import "errors"

func GetDomain() (string, error) {
	if agolloClient == nil {
		return "", errors.New(emptyConfigClientErr)
	}
	v, err := agolloClient.GetViper("gohu.yaml")
	if err != nil {
		return "", errors.New(getViperErr)
	}
	return v.GetString("Domain"), nil
}

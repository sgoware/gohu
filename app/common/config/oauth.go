package config

func GetClientSecret(clientId string) (clientSecret string, err error) {
	if agolloClient == nil {
		return "", errEmptyConfigClient
	}
	v, err := agolloClient.GetViper("oauth.yaml")
	if err != nil {
		return "", errGetViper
	}

	clientSecret = v.GetString("Client." + clientId + ".Secret")
	if clientSecret == "" {
		return "", errViperEmptyKey
	}
	return clientSecret, nil
}

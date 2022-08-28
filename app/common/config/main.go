package config

func GetDomain() (string, error) {
	if agolloClient == nil {
		return "", errEmptyConfigClient
	}
	v, err := agolloClient.GetViper("gohu.yaml")
	if err != nil {
		return "", errGetViper
	}
	return v.GetString("Domain"), nil
}

func GetMainDomain() (string, error) {
	if agolloClient == nil {
		return "", errEmptyConfigClient
	}
	v, err := agolloClient.GetViper("gohu.yaml")
	if err != nil {
		return "", errGetViper
	}
	return v.GetString("MainDomain"), nil
}

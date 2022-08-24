package config

func (c *Agollo) GetDomain() (string, error) {
	v, err := c.GetViper("gohu.yaml")
	if err != nil {
		return "", err
	}
	return v.GetString("Domain"), nil
}

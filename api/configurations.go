package api

import (
	"encoding/base64"
	"errors"
	"strings"

	"github.com/golang-jwt/jwt/v4"
)

type Configuration struct {
	ConsumerKey    string
	ConsumerSecret string
	Host           string
	Base64Token    string
	BearerToken    string
	ViewToken      string
	Sid            string
	UserId         string
	EditToken      string
	EditSid        string
	EditRid        string
	ServerId       string
	LoginParams    string
	NeoFinKey      string
}

func NewConfiguration(consumerKey, consumerSecret, host, accessToken, neoFinKey string) *Configuration {
	config := &Configuration{
		ConsumerKey:    consumerKey,
		ConsumerSecret: consumerSecret,
		Host:           host,
		BearerToken:    accessToken,
		NeoFinKey:      neoFinKey,
	}
	config.Base64Token = config.convertBase64()
	return config
}

func (c *Configuration) convertBase64() string {
	base64String := c.ConsumerKey + ":" + c.ConsumerSecret
	base64Token := base64.StdEncoding.EncodeToString([]byte(base64String))
	return base64Token
}

func (c *Configuration) extractUserId(viewToken string) (string, error) {
	if viewToken == "" {
		return "", errors.New("view Token hasn't been generated. Kindly call the login function and try to generate OTP")
	}
	parser := jwt.NewParser(jwt.WithoutClaimsValidation())
	token, _, err := parser.ParseUnverified(viewToken, jwt.MapClaims{})
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("Unable to parse claims")
	}

	userid, ok := claims["sub"].(string)
	if !ok {
		return "", errors.New("User ID not found in token")
	}

	c.UserId = userid
	return userid, nil
}

func (c *Configuration) getDomain(sessionInit bool) (string, error) {
	hostList := []string{"prod", "uat"}
	host := strings.ToLower(strings.TrimSpace(c.Host))

	for _, h := range hostList {
		if host == h {
			if sessionInit {
				if host == "prod" {
					return SessionProdBaseUrl, nil
				}
				return SessionUatBaseUrl, nil
			}
			if host == "prod" {
				return ProdBaseUrl, nil
			}
			return UatBaseUrl, nil
		}
	}
	return "", errors.New("Either UAT or PROD in Environment accepted")
}

func (c *Configuration) getUrlDetails(apiInfo string) (string, error) {
	domainInfo, err := c.getDomain(false)
	if err != nil {
		return "", err
	}

	host := strings.ToLower(strings.TrimSpace(c.Host))
	if host == "prod" {
		return domainInfo + ProdUrl[apiInfo], nil
	}
	return domainInfo + UatUrl[apiInfo], nil
}

func (c *Configuration) getNeoFinKey() string {
	host := strings.ToLower(strings.TrimSpace(c.Host))
	if host == "prod" {
		if c.NeoFinKey != "" {
			return c.NeoFinKey
		}
		return "X6Nk8cQhUgGmJ2vBdWw4sfzrz4L5En"
	}
	if c.NeoFinKey != "" {
		return c.NeoFinKey
	}
	return "bQJNkL5z8m4aGcRgjDvXhHfSx7VpZnE"
}

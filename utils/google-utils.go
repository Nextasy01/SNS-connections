package utils

import (
	"github.com/Nextasy01/SNS-connections/entity"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
)

type EnvReader interface {
	ReadFromEnv() (string, string, error)
	GetAppEnv() (string, error)
}

type GoogleEnvReader struct{}

func (g *GoogleEnvReader) ReadFromEnv() (string, string, string, string, error) {
	envFile, err := godotenv.Read(".env")

	if err != nil {
		return "", "", "", "", err
	}

	return envFile["Google_Client_id"], envFile["Google_Secret_key"], envFile["Google_Service_Account"], envFile["Google_Service_Account_Private_Key"], nil
}

func (g *GoogleEnvReader) GetAppEnv() (string, error) {
	envFile, err := godotenv.Read(".env")

	if err != nil {
		return "", err
	}

	return envFile["APP_ENV"], nil

}

func (g *GoogleEnvReader) GetProdRedirectUrlEnv() (string, error) {
	envFile, err := godotenv.Read(".env")

	if err != nil {
		return "", err
	}

	return envFile["PROD_REDIRECT_URI"], nil

}

func NewToken(acc *entity.GoogleAccount) *oauth2.Token {
	token := new(oauth2.Token)
	token.AccessToken = acc.AccessToken
	token.RefreshToken = acc.RefreshToken
	token.Expiry = acc.ExpiresAt
	token.TokenType = acc.TokenType

	return token
}

func NewConfig() (*oauth2.Config, error) {
	g := new(GoogleEnvReader)
	client_id, secret_key, _, _, err := g.ReadFromEnv()

	if err != nil {
		return nil, err
	}

	conf := &oauth2.Config{
		ClientID:     client_id,
		ClientSecret: secret_key,
		Scopes:       []string{"email", "profile", "https://www.googleapis.com/auth/youtube", "https://www.googleapis.com/auth/youtube.upload", "https://www.googleapis.com/auth/youtube.readonly"},
		Endpoint:     google.Endpoint,
		RedirectURL:  "http://localhost:9000",
	}

	return conf, nil
}

func NewConfigServiceAccount() (*jwt.Config, error) {
	g := new(GoogleEnvReader)
	_, _, email, private_key, err := g.ReadFromEnv()

	if err != nil {
		return nil, err
	}

	conf := &jwt.Config{
		Email:      email,
		PrivateKey: []byte(private_key),
		Scopes: []string{
			"https://www.googleapis.com/auth/drive",
			"https://www.googleapis.com/auth/drive.appdata",
			"https://www.googleapis.com/auth/drive.readonly",
		},
		TokenURL: google.JWTTokenURL,
	}

	return conf, nil
}

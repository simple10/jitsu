package authorization

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jitsucom/jitsu/configurator/random"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
	"io/ioutil"
	"net/url"
	"time"
)

const (
	BoxyHQ = "boxyhq"
)

type SSOProvider interface {
	GetUser(code string) (*UserEntity, error)
	AuthLink() string
	AccessTokenTTL() time.Duration
	Name() string
}

type UserEntity struct {
	ID          string
	Email       string
	AccessToken string
}

type UserInfoEntity struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

type BoxyHQProvider struct {
	ssoConfig *SSOConfig
}

func NewBoxyHQProvider(ssoConfig *SSOConfig) (*BoxyHQProvider, error) {
	return &BoxyHQProvider{ssoConfig}, nil
}

func (bp *BoxyHQProvider) GetUser(code string) (*UserEntity, error) {
	conf := &clientcredentials.Config{
		ClientID:       "dummy",
		ClientSecret:   "dummy",
		EndpointParams: url.Values{"tenant": {bp.ssoConfig.Tenant}, "product": {bp.ssoConfig.Product}, "grant_type": {"authorization_code"}, "code": {code}},
		TokenURL:       bp.ssoConfig.Host + "/api/oauth/token",
		AuthStyle:      oauth2.AuthStyleInParams,
	}

	ctx := context.Background()
	token, err := conf.Token(ctx)
	if err != nil || token == nil {
		return nil, fmt.Errorf("can't get token from sso server")
	}

	userInfo, err := bp.GetUserInfo(conf)
	if err != nil {
		return nil, fmt.Errorf("can't get user info from sso server %v", err)
	}

	userEntity := UserEntity{
		ID:          userInfo.ID,
		Email:       userInfo.Email,
		AccessToken: token.AccessToken,
	}

	return &userEntity, nil
}

func (bp *BoxyHQProvider) GetUserInfo(conf *clientcredentials.Config) (*UserInfoEntity, error) {
	client := conf.Client(context.Background())
	response, err := client.Get(bp.ssoConfig.Host + "/api/oauth/userinfo")
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	userInfoEntity := &UserInfoEntity{}
	err = json.Unmarshal(body, userInfoEntity)
	if err != nil {
		return nil, err
	}

	return userInfoEntity, nil
}

func (bp *BoxyHQProvider) AuthLink() string {
	authLink := fmt.Sprintf(
		"%s?response_type=code&provider=saml&client_id=dummy&tenant=%s&product=%s&state=%s",
		bp.ssoConfig.Host+"/api/oauth/authorize",
		bp.ssoConfig.Tenant,
		bp.ssoConfig.Product,
		random.String(10),
	)

	return authLink
}

func (bp *BoxyHQProvider) Name() string {
	return BoxyHQ
}

func (bp *BoxyHQProvider) AccessTokenTTL() time.Duration {
	return bp.ssoConfig.AccessTokenTTLSeconds
}
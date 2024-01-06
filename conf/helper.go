package conf

import (
	"fmt"
	"net/url"
	"time"

	"github.com/VaalaCat/frp-panel/common"
	"github.com/VaalaCat/frp-panel/utils"
	v1 "github.com/fatedier/frp/pkg/config/v1"
	"github.com/sirupsen/logrus"
)

func RPCListenAddr() string {
	cfg := Get()
	return fmt.Sprintf(":%d", cfg.Master.RPCPort)
}

func RPCCallAddr() string {
	cfg := Get()
	return fmt.Sprintf("%s:%d", cfg.Master.RPCHost, cfg.Master.RPCPort)
}

func InternalFRPServerToken() string {
	cfg := Get()
	return utils.MD5(fmt.Sprintf("%s:%d:%s",
		cfg.Master.InternalFRPServerHost,
		cfg.Master.InternalFRPServerPort,
		cfg.App.GlobalSecret))
}

func JWTSecret() string {
	cfg := Get()
	return utils.SHA1(fmt.Sprintf("%s:%d:%s", cfg.Master.APIHost, cfg.Master.APIPort, cfg.App.GlobalSecret))
}

func MasterAPIListenAddr() string {
	cfg := Get()
	return fmt.Sprintf(":%d", cfg.Master.APIPort)
}

func ServerAPIListenAddr() string {
	cfg := Get()
	return fmt.Sprintf(":%d", cfg.Server.APIPort)
}

func FRPsAuthOption() v1.HTTPPluginOptions {
	cfg := Get()
	authUrl, err := url.Parse(fmt.Sprintf("http://%s:%d%s", cfg.Master.InternalFRPAuthServerHost,
		cfg.Master.InternalFRPAuthServerPort,
		cfg.Master.InternalFRPAuthServerPath))
	if err != nil {
		logrus.WithError(err).Fatalf("parse auth url error")
	}
	return v1.HTTPPluginOptions{
		Name: "multiuser",
		Ops:  []string{"Login"},
		Addr: authUrl.Host,
		Path: authUrl.Path,
	}
}

func GetCommonJWT(uid string) string {
	token, _ := utils.GetJwtTokenFromMap(JWTSecret(),
		time.Now().Unix(),
		int64(Get().App.CookieAge),
		map[string]string{common.UserIDKey: uid})
	return token
}

func GetCommonJWTWithExpireTime(uid string, expSec int) string {
	token, _ := utils.GetJwtTokenFromMap(JWTSecret(),
		time.Now().Unix(),
		int64(expSec),
		map[string]string{common.UserIDKey: uid})
	return token
}
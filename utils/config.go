package utils

import "github.com/spf13/viper"

func GetHttpsEnable () bool {
	return viper.GetBool("server.restful.httpsenable")
}

func GetHttpsCert () string {
	return viper.GetString("server.restful.httpsclientcert")
}

func GetCertCommonName() string {
	return viper.GetString("server.restful.certcommonname")
}

func GetClientPort () string {
	return viper.GetString("server.restful.clientport")
}
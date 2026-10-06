package models

type Config struct {
	AppName       string `json:"AppName"`
	ServiceName   string `json:"ServiceName"`
	ServiceDesc   string `json:"ServiceDesc"`
	WebServerAddr string `json:"WebServerAddr"`
}

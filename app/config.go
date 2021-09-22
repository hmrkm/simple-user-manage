package main

type Config struct {
	MysqlUser              string `split_words:"true"`
	MysqlPassword          string `split_words:"true"`
	MysqlDatabase          string `split_words:"true"`
	AuthenticationEndpoint string `split_words:"true"`
	HttpRetryNumber        int    `split_words:"true"`
	HttpSleepSecond        int    `split_words:"true"`
}

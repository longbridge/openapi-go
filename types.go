package openapi

type Market string

const (
	MarketUS Market = "US"
	MarketUK Market = "UK"
	MarketHK Market = "HK"
	MarketCN Market = "CN"
	MarketSG Market = "SG"
)

type Language string

const (
	LanguageZHCN Language = "zh-CN"
	LanguageZHHK Language = "zh-HK"
	LanguageEN   Language = "en"
)

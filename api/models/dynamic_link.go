package models

type DynamicLinkInfo struct {
	Host                    string                  `json:"host"`
	Link                    string                  `json:"link"`
	AndroidParameters       AndroidParameters       `json:"androidParameters,omitempty"`
	IosParameters           IosParameters           `json:"iosParameters,omitempty"`
	OtherPlatformParameters OtherPlatformParameters `json:"otherPlatformParameters,omitempty"`
	AnalyticsInfo           AnalyticsInfo           `json:"analyticsInfo,omitempty"`
	SocialMetaTagInfo       SocialMetaTagInfo       `json:"socialMetaTagInfo,omitempty"`
}

type AndroidParameters struct {
	AndroidPackageName           string `json:"androidPackageName,omitempty"`
	AndroidFallbackLink          string `json:"androidFallbackLink,omitempty"`
	AndroidMinPackageVersionCode string `json:"androidMinPackageVersionCode,omitempty"`
}

type IosParameters struct {
	IosFallbackLink     string `json:"iosFallbackLink,omitempty"`
	IosIpadFallbackLink string `json:"iosIpadFallbackLink,omitempty"`
	IosAppStoreId       string `json:"iosAppStoreId,omitempty"`
}

type OtherPlatformParameters struct {
	FallbackURL string `json:"ofl,omitempty"`
}

type AnalyticsInfo struct {
	MarketingParameters    MarketingParameters    `json:"marketingParameters,omitempty"`
	ItunesConnectAnalytics ItunesConnectAnalytics `json:"itunesConnectAnalytics,omitempty"`
}

type MarketingParameters struct {
	UtmSource   string `json:"utmSource,omitempty"`
	UtmMedium   string `json:"utmMedium,omitempty"`
	UtmCampaign string `json:"utmCampaign,omitempty"`
	UtmTerm     string `json:"utmTerm,omitempty"`
	UtmContent  string `json:"utmContent,omitempty"`
}

type ItunesConnectAnalytics struct {
	At string `json:"at,omitempty"`
	Ct string `json:"ct,omitempty"`
	Mt string `json:"mt,omitempty"`
	Pt string `json:"pt,omitempty"`
}

type SocialMetaTagInfo struct {
	SocialTitle       string `json:"socialTitle,omitempty"`
	SocialDescription string `json:"socialDescription,omitempty"`
	SocialImageLink   string `json:"socialImageLink,omitempty"`
}

type Suffix struct {
	Option string `json:"option,omitempty"` // "SHORT" or "UNGUESSABLE"
}

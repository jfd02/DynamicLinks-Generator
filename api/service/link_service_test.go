package service

import (
	"os"
	"testing"

	"dynamic-links-generator/api/models"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log := zerolog.New(os.Stdout).With().Timestamp().Logger()
	_ = log
	os.Exit(m.Run())
}

func TestParseLongDynamicLink(t *testing.T) {
	tests := []struct {
		name     string
		longLink string
		want     models.CreateDynamicLinkRequest
		wantErr  bool
	}{
		{
			name: "complete link with all parameters",
			longLink: "https://example.com?link=https://target.com" +
				"&apn=com.android.app" +
				"&afl=https://android-fallback.com" +
				"&amv=123" +
				"&isi=123456789" +
				"&ifl=https://ios-fallback.com" +
				"&ipfl=https://ipad-fallback.com" +
				"&ofl=https://other-platform-fallback.com" +
				"&utm_source=source" +
				"&utm_medium=medium" +
				"&utm_campaign=campaign" +
				"&utm_term=term" +
				"&utm_content=content" +
				"&at=at" +
				"&ct=ct" +
				"&mt=mt" +
				"&pt=pt" +
				"&st=social title" +
				"&sd=social description" +
				"&si=https://social-image.com" +
				"&path=SHORT",
			want: models.CreateDynamicLinkRequest{
				DynamicLinkInfo: models.DynamicLinkInfo{
					Host: "example.com",
					Link: "https://target.com",
					AndroidParameters: models.AndroidParameters{
						AndroidPackageName:           "com.android.app",
						AndroidFallbackLink:          "https://android-fallback.com",
						AndroidMinPackageVersionCode: "123",
					},
					IosParameters: models.IosParameters{
						IosAppStoreId:       "123456789",
						IosFallbackLink:     "https://ios-fallback.com",
						IosIpadFallbackLink: "https://ipad-fallback.com",
					},
					OtherPlatformParameters: models.OtherPlatformParameters{
						FallbackURL: "https://other-platform-fallback.com",
					},
					AnalyticsInfo: models.AnalyticsInfo{
						MarketingParameters: models.MarketingParameters{
							UtmSource:   "source",
							UtmMedium:   "medium",
							UtmCampaign: "campaign",
							UtmTerm:     "term",
							UtmContent:  "content",
						},
						ItunesConnectAnalytics: models.ItunesConnectAnalytics{
							At: "at",
							Ct: "ct",
							Mt: "mt",
							Pt: "pt",
						},
					},
					SocialMetaTagInfo: models.SocialMetaTagInfo{
						SocialTitle:       "social title",
						SocialDescription: "social description",
						SocialImageLink:   "https://social-image.com",
					},
				},
				Suffix: models.Suffix{
					Option: "SHORT",
				},
			},
			wantErr: false,
		},
		{
			name:     "invalid URL format",
			longLink: "not a valid url",
			want:     models.CreateDynamicLinkRequest{},
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &linkService{}
			got, err := service.ParseLongDynamicLink(tt.longLink)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"dynamic-links-generator/api/apperrors"
	"dynamic-links-generator/api/models"
	"dynamic-links-generator/api/repository"
	"dynamic-links-generator/config"
	"dynamic-links-generator/utils"

	"github.com/rs/zerolog/log"
)

type LinkService interface {
	CreateDynamicLink(ctx context.Context, params models.CreateDynamicLinkRequest) (*models.ShortLinkResponse, error)
	ParseLongDynamicLink(longLink string) (models.CreateDynamicLinkRequest, error)
	ResolveShortPath(ctx context.Context, rawURL string) (*models.LongLinkResponse, error)
	PrepareDynamicLinkRequest(input map[string]any) (models.CreateDynamicLinkRequest, error)
}

type linkService struct {
	repo repository.LinkRepository
	cfg  *config.Config
}

func NewLinkService(repo repository.LinkRepository, cfg *config.Config) *linkService {
	return &linkService{
		repo: repo,
		cfg:  cfg,
	}
}

func (s *linkService) getLongLinkFromHostAndPath(
	ctx context.Context,
	host string,
	path string,
) (*models.LongLinkResponse, error) {
	rawQueryStr, err := s.repo.GetQueryParamsByHostAndPath(ctx, host, path)
	if err != nil {
		return nil, err
	}

	longLink := fmt.Sprintf("%s://%s/%s", s.cfg.URLScheme, host, path)
	if rawQueryStr != "" {
		longLink += "?" + rawQueryStr
	}

	log.Debug().
		Str("path", path).
		Str("long_link", longLink).
		Msg("Link retrieved from service")

	return &models.LongLinkResponse{
		LongLink: longLink,
	}, nil
}

func (s *linkService) CreateDynamicLink(ctx context.Context, params models.CreateDynamicLinkRequest) (*models.ShortLinkResponse, error) {
	warnings := []models.Warning{}

	log.Debug().
		Str("params", fmt.Sprintf("%+v", params)).
		Msg("Dynamic link parameters")

	host, err := utils.CleanHost(params.DynamicLinkInfo.Host)
	if err != nil {
		log.Error().
			Str("host", params.DynamicLinkInfo.Host).
			Msg("Invalid host")
		return nil, fmt.Errorf("invalid host: %w", err)
	}

	if !utils.IsDomainAllowed(s.cfg.DomainAllowList, params.DynamicLinkInfo.Link) {
		log.Error().
			Str("link", params.DynamicLinkInfo.Link).
			Msg("Domain link not in allow list")
		return nil, apperrors.ErrDomainLinkNotAllowed
	}

	isi := params.DynamicLinkInfo.IosParameters.IosAppStoreId

	if isi != "" {
		if !utils.IsNumericString(isi) {
			return nil, apperrors.ErrInvalidAppStoreID
		}
	}

	queryParams := url.Values{}
	queryParams.Add("link", params.DynamicLinkInfo.Link)

	addParam := func(key, value string) {
		if value != "" {
			queryParams.Add(key, value)
		}
	}

	addParam("apn", params.DynamicLinkInfo.AndroidParameters.AndroidPackageName)
	addParam("afl", params.DynamicLinkInfo.AndroidParameters.AndroidFallbackLink)
	addParam("amv", params.DynamicLinkInfo.AndroidParameters.AndroidMinPackageVersionCode)

	addParam("ifl", params.DynamicLinkInfo.IosParameters.IosFallbackLink)
	addParam("ipfl", params.DynamicLinkInfo.IosParameters.IosIpadFallbackLink)
	addParam("isi", isi)

	addParam("ofl", params.DynamicLinkInfo.OtherPlatformParameters.FallbackURL)

	addParam("st", params.DynamicLinkInfo.SocialMetaTagInfo.SocialTitle)
	addParam("sd", params.DynamicLinkInfo.SocialMetaTagInfo.SocialDescription)

	si := params.DynamicLinkInfo.SocialMetaTagInfo.SocialImageLink

	addParam("si", si)

	if si != "" {
		if !utils.IsURL(si) {
			warnings = append(warnings, models.Warning{
				WarningCode:    "MALFORMED_PARAM",
				WarningMessage: "Param 'si' is not a valid URL",
			})
		}
	}

	addParam("utm_source", params.DynamicLinkInfo.AnalyticsInfo.MarketingParameters.UtmSource)
	addParam("utm_medium", params.DynamicLinkInfo.AnalyticsInfo.MarketingParameters.UtmMedium)
	addParam("utm_campaign", params.DynamicLinkInfo.AnalyticsInfo.MarketingParameters.UtmCampaign)
	addParam("utm_term", params.DynamicLinkInfo.AnalyticsInfo.MarketingParameters.UtmTerm)
	addParam("utm_content", params.DynamicLinkInfo.AnalyticsInfo.MarketingParameters.UtmContent)
	pt := params.DynamicLinkInfo.AnalyticsInfo.ItunesConnectAnalytics.Pt
	addParam("pt", pt)

	if isi == "" {
		if at := params.DynamicLinkInfo.AnalyticsInfo.ItunesConnectAnalytics.At; at != "" {
			warnings = append(warnings, models.Warning{
				WarningCode:    "UNRECOGNIZED_PARAM",
				WarningMessage: "Param 'at' is not needed, since 'isi' is not specified.",
			})
		}
		if ct := params.DynamicLinkInfo.AnalyticsInfo.ItunesConnectAnalytics.Ct; ct != "" {
			warnings = append(warnings, models.Warning{
				WarningCode:    "UNRECOGNIZED_PARAM",
				WarningMessage: "Param 'ct' is not needed, since 'isi' is not specified.",
			})
		}
		if mt := params.DynamicLinkInfo.AnalyticsInfo.ItunesConnectAnalytics.Mt; mt != "" {
			warnings = append(warnings, models.Warning{
				WarningCode:    "UNRECOGNIZED_PARAM",
				WarningMessage: "Param 'mt' is not needed, since 'isi' is not specified.",
			})
		}
		if pt := params.DynamicLinkInfo.AnalyticsInfo.ItunesConnectAnalytics.Pt; pt != "" {
			warnings = append(warnings, models.Warning{
				WarningCode:    "UNRECOGNIZED_PARAM",
				WarningMessage: "Param 'pt' is not needed, since 'isi' is not specified.",
			})
		}
	}

	if pt == "" {
		if at := params.DynamicLinkInfo.AnalyticsInfo.ItunesConnectAnalytics.At; at != "" {
			warnings = append(warnings, models.Warning{
				WarningCode:    "UNRECOGNIZED_PARAM",
				WarningMessage: "Param 'at' is not needed, since 'pt' is not specified.",
			})
		}
		if ct := params.DynamicLinkInfo.AnalyticsInfo.ItunesConnectAnalytics.Ct; ct != "" {
			warnings = append(warnings, models.Warning{
				WarningCode:    "UNRECOGNIZED_PARAM",
				WarningMessage: "Param 'ct' is not needed, since 'pt' is not specified.",
			})
		}
		if mt := params.DynamicLinkInfo.AnalyticsInfo.ItunesConnectAnalytics.Mt; mt != "" {
			warnings = append(warnings, models.Warning{
				WarningCode:    "UNRECOGNIZED_PARAM",
				WarningMessage: "Param 'mt' is not needed, since 'pt' is not specified.",
			})
		}
	}

	addParam("at", params.DynamicLinkInfo.AnalyticsInfo.ItunesConnectAnalytics.At)
	addParam("ct", params.DynamicLinkInfo.AnalyticsInfo.ItunesConnectAnalytics.Ct)
	addParam("mt", params.DynamicLinkInfo.AnalyticsInfo.ItunesConnectAnalytics.Mt)

	shortPath := params.Suffix.Option == "SHORT"
	response, err := s.createOrGetShortLink(ctx, host, queryParams, shortPath)
	if err != nil {
		return nil, err
	}

	response.Warnings = warnings
	return response, nil
}

func (s *linkService) ParseLongDynamicLink(longDynamicLink string) (models.CreateDynamicLinkRequest, error) {
	var req models.CreateDynamicLinkRequest

	log.Debug().
		Str("long_link", longDynamicLink).
		Msg("Parsing long dynamic link")

	u, err := url.Parse(longDynamicLink)
	if err != nil {
		return req, apperrors.ErrInvalidURLFormat
	}

	if u.Host == "" {
		return req, apperrors.ErrHostInvalid
	}

	req.DynamicLinkInfo.Host = u.Host

	params := u.Query()

	req.DynamicLinkInfo.Link = params.Get("link")

	log.Debug().
		Str("link", req.DynamicLinkInfo.Link).
		Msg("Parsed link")

	if apn := params.Get("apn"); apn != "" {
		req.DynamicLinkInfo.AndroidParameters.AndroidPackageName = apn
	}
	if afl := params.Get("afl"); afl != "" {
		req.DynamicLinkInfo.AndroidParameters.AndroidFallbackLink = afl
	}
	if apv := params.Get("amv"); apv != "" {
		req.DynamicLinkInfo.AndroidParameters.AndroidMinPackageVersionCode = apv
	}

	if isi := params.Get("isi"); isi != "" {
		req.DynamicLinkInfo.IosParameters.IosAppStoreId = isi
	}
	if ifl := params.Get("ifl"); ifl != "" {
		req.DynamicLinkInfo.IosParameters.IosFallbackLink = ifl
	}
	if iflIpad := params.Get("ipfl"); iflIpad != "" {
		req.DynamicLinkInfo.IosParameters.IosIpadFallbackLink = iflIpad
	}

	if ofl := params.Get("ofl"); ofl != "" {
		req.DynamicLinkInfo.OtherPlatformParameters.FallbackURL = ofl
	}

	if utmSource := params.Get("utm_source"); utmSource != "" {
		req.DynamicLinkInfo.AnalyticsInfo.MarketingParameters.UtmSource = utmSource
	}
	if utmMedium := params.Get("utm_medium"); utmMedium != "" {
		req.DynamicLinkInfo.AnalyticsInfo.MarketingParameters.UtmMedium = utmMedium
	}
	if utmCampaign := params.Get("utm_campaign"); utmCampaign != "" {
		req.DynamicLinkInfo.AnalyticsInfo.MarketingParameters.UtmCampaign = utmCampaign
	}
	if utmTerm := params.Get("utm_term"); utmTerm != "" {
		req.DynamicLinkInfo.AnalyticsInfo.MarketingParameters.UtmTerm = utmTerm
	}
	if utmContent := params.Get("utm_content"); utmContent != "" {
		req.DynamicLinkInfo.AnalyticsInfo.MarketingParameters.UtmContent = utmContent
	}
	if at := params.Get("at"); at != "" {
		req.DynamicLinkInfo.AnalyticsInfo.ItunesConnectAnalytics.At = at
	}
	if ct := params.Get("ct"); ct != "" {
		req.DynamicLinkInfo.AnalyticsInfo.ItunesConnectAnalytics.Ct = ct
	}
	if mt := params.Get("mt"); mt != "" {
		req.DynamicLinkInfo.AnalyticsInfo.ItunesConnectAnalytics.Mt = mt
	}
	if pt := params.Get("pt"); pt != "" {
		req.DynamicLinkInfo.AnalyticsInfo.ItunesConnectAnalytics.Pt = pt
	}

	if socialTitle := params.Get("st"); socialTitle != "" {
		req.DynamicLinkInfo.SocialMetaTagInfo.SocialTitle = socialTitle
	}
	if socialDescription := params.Get("sd"); socialDescription != "" {
		req.DynamicLinkInfo.SocialMetaTagInfo.SocialDescription = socialDescription
	}
	if socialImageLink := params.Get("si"); socialImageLink != "" {
		req.DynamicLinkInfo.SocialMetaTagInfo.SocialImageLink = socialImageLink
	}

	if pathOption := params.Get("path"); pathOption != "" {
		req.Suffix.Option = pathOption
	}

	log.Debug().
		Str("req", fmt.Sprintf("%+v", req)).
		Msg("Parsed long dynamic link")

	return req, nil
}

func (s *linkService) createOrGetShortLink(
	ctx context.Context,
	host string,
	queryParams url.Values,
	shortPath bool,
) (*models.ShortLinkResponse, error) {
	rawQS := queryParams.Encode()
	if shortPath {
		if path, err := s.findExistingShortLink(ctx, host, rawQS); err == nil {
			full := fmt.Sprintf("%s://%s/%s", s.cfg.URLScheme, host, path)
			log.Debug().
				Str("path", path).
				Str("query_params", rawQS).
				Msg("Re‑using existing short link")
			return &models.ShortLinkResponse{ShortLink: full, Warnings: []models.Warning{}}, nil

		} else if err != sql.ErrNoRows {
			log.Error().
				Err(err).
				Msg("Error querying for existing short link")
			return nil, err
		}
	}

	length := s.cfg.ShortPathLength
	if !shortPath {
		length = s.cfg.UnguessablePathLength
	}
	path := utils.GenerateDynamicLinkPath(length)

	if err := s.createShortLink(ctx, host, path, rawQS, !shortPath); err != nil {
		return nil, fmt.Errorf("failed to store link: %w", err)
	}

	full := fmt.Sprintf("%s://%s/%s", s.cfg.URLScheme, host, path)
	log.Debug().
		Str("path", path).
		Str("query_params", rawQS).
		Msg("New link stored in database")

	return &models.ShortLinkResponse{ShortLink: full, Warnings: []models.Warning{}}, nil
}

func (s *linkService) findExistingShortLink(
	ctx context.Context,
	host, rawQS string,
) (string, error) {
	return s.repo.FindExistingShortLink(ctx, host, rawQS)
}

func (s *linkService) createShortLink(
	ctx context.Context,
	host, path, rawQS string,
	unguessable bool,
) error {
	return s.repo.CreateShortLink(ctx, host, path, rawQS, unguessable)
}

func (s *linkService) ResolveShortPath(ctx context.Context, rawURL string) (*models.LongLinkResponse, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, apperrors.ErrInvalidRequestedLink
	}

	host := u.Host
	pathParts := strings.Split(strings.Trim(u.Path, "/"), "/")
	if len(pathParts) != 1 {
		return nil, fmt.Errorf("unexpected path format: %w", apperrors.ErrInvalidPathFormat)
	}

	return s.getLongLinkFromHostAndPath(ctx, host, pathParts[0])
}

func (s *linkService) PrepareDynamicLinkRequest(input map[string]any) (models.CreateDynamicLinkRequest, error) {
	var req models.CreateDynamicLinkRequest

	if longLink, ok := input["longDynamicLink"].(string); ok && longLink != "" {
		parsedReq, err := s.ParseLongDynamicLink(longLink)
		if err != nil {
			return models.CreateDynamicLinkRequest{}, err
		}
		req = parsedReq
	} else {
		reqBytes, err := json.Marshal(input)
		if err != nil {
			return models.CreateDynamicLinkRequest{}, apperrors.ErrInvalidFormat
		}
		if err := json.Unmarshal(reqBytes, &req); err != nil {
			return models.CreateDynamicLinkRequest{}, apperrors.ErrInvalidFormat
		}
	}

	if req.DynamicLinkInfo.Host == "" {
		return models.CreateDynamicLinkRequest{}, apperrors.ErrMissingHost
	}
	if req.DynamicLinkInfo.Link == "" {
		return models.CreateDynamicLinkRequest{}, apperrors.ErrMissingLink
	}
	if err := utils.ValidateURLScheme(req.DynamicLinkInfo.Link); err != nil {
		return models.CreateDynamicLinkRequest{}, err
	}

	return req, nil
}

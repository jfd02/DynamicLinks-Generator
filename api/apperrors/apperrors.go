package apperrors

import "errors"

var (
	ErrInvalidURLFormat  = errors.New("invalid URL format")
	ErrHostInvalid       = errors.New("host is invalid")
	ErrInvalidAppStoreID = errors.New("app store id should contain numbers only")

	ErrDomainLinkNotAllowed = errors.New("domain link not in allow list")
	ErrInvalidPathFormat    = errors.New("path must contain exactly one segment")
	ErrInvalidRequestedLink = errors.New("invalid requested link")

	ErrInvalidFormat = errors.New("invalid request format")
	ErrMissingHost   = errors.New("missing host")
	ErrMissingLink   = errors.New("missing link")

	ErrLinkNotFound = errors.New("link not found")
)

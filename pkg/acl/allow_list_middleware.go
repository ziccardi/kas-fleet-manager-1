package acl

import (
	"fmt"
	"net/http"

	"gitlab.cee.redhat.com/service/managed-services-api/pkg/auth"
	"gitlab.cee.redhat.com/service/managed-services-api/pkg/errors"
	"gitlab.cee.redhat.com/service/managed-services-api/pkg/services"
	"gitlab.cee.redhat.com/service/managed-services-api/pkg/shared"
)

type AllowListMiddleware struct {
	configService services.ConfigService
}

func NewAllowListMiddleware(configService services.ConfigService) *AllowListMiddleware {
	middleware := AllowListMiddleware{
		configService: configService,
	}

	return &middleware
}

// Middleware handler to authorize users based on the provided ACL configuration
func (middleware *AllowListMiddleware) Authorize(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if middleware.configService.IsAllowListEnabled() {
			context := r.Context()
			username, orgId := auth.GetUsernameFromContext(context), auth.GetOrgIdFromContext(context)
			org, _ := middleware.configService.GetOrganisationById(orgId)
			var userIsAllowed, userIsAllowedAsServiceAccount bool

			// check if user is allowed within the organisation
			if org.IsUserAllowed(username) {
				userIsAllowed = true
				userIsAllowedAsServiceAccount = false
			} else {
				// check if user is allowed a service account
				_, found := middleware.configService.GetServiceAccountByUsername(username)
				userIsAllowed = found
				userIsAllowedAsServiceAccount = found
			}

			if !userIsAllowed {
				shared.HandleError(r.Context(), w, errors.ErrorForbidden, fmt.Sprintf("User '%s' is not authorized to access the service.", username))
				return
			}

			context = auth.SetUserIsAllowedAsServiceAccountContext(context, userIsAllowedAsServiceAccount)
			*r = *r.WithContext(context)
		}

		next.ServeHTTP(w, r)
	})
}

package middleware

import (
	apphttp "service-core/internal/common/http"
)

type Middleware func(apphttp.AppHandler) apphttp.AppHandler

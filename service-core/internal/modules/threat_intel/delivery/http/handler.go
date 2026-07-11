package http

import (
	"errors"
	"net/http"

	apperrors "service-core/internal/common/errors"
	apphttp "service-core/internal/common/http"
	"service-core/internal/modules/threat_intel/domain"
	"service-core/internal/modules/threat_intel/usecase"

	"github.com/go-chi/chi/v5"
)

type threatIntelHandler struct {
	analyzeIP *usecase.AnalyzeIPUsecase
	getGeoIP  *usecase.GetGeoIPUsecase
}

func NewThreatIntelHandler(
	analyzeIP *usecase.AnalyzeIPUsecase,
	getGeoIP *usecase.GetGeoIPUsecase,
) *threatIntelHandler {
	return &threatIntelHandler{
		analyzeIP: analyzeIP,
		getGeoIP:  getGeoIP,
	}
}

func (h *threatIntelHandler) AnalyzeIP(w http.ResponseWriter, r *http.Request) error {
	ip := chi.URLParam(r, "ip")
	if ip == "" {
		return apperrors.NewBadRequest("IP address is required")
	}

	apiKey := r.Header.Get("X-VT-Key")

	report, err := h.analyzeIP.Execute(r.Context(), ip, apiKey)
	if err != nil {
		if errors.Is(err, domain.ErrAPIKeyRequired) {
			return apperrors.NewBadRequest("VirusTotal API key is required but not configured")
		}
		if errors.Is(err, domain.ErrInvalidIP) {
			return apperrors.NewBadRequest("invalid IP address")
		}
		return err
	}

	apphttp.WriteJSON(w, http.StatusOK, report.RawReport)
	return nil
}

func (h *threatIntelHandler) GetGeolocation(w http.ResponseWriter, r *http.Request) error {
	ip := chi.URLParam(r, "ip")
	if ip == "" {
		return apperrors.NewBadRequest("IP address is required")
	}

	report, err := h.getGeoIP.Execute(r.Context(), ip)
	if err != nil {
		return err
	}

	apphttp.WriteJSON(w, http.StatusOK, report.RawReport)
	return nil
}

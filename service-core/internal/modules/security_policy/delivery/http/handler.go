package http

import (
	"net/http"

	apperrors "service-core/internal/common/errors"
	apphttp "service-core/internal/common/http"
	"service-core/internal/modules/security_policy/domain"
	"service-core/internal/modules/security_policy/usecase"

	"github.com/go-chi/chi/v5"
)

type securityPolicyHandler struct {
	listRules      *usecase.ListRulesUsecase
	createRule     *usecase.CreateRuleUsecase
	toggleRule     *usecase.ToggleRuleUsecase
	updateRule     *usecase.UpdateRuleUsecase
	deleteRule     *usecase.DeleteRuleUsecase
	getIPConfig    *usecase.GetIPConfigUsecase
	updateIPAction *usecase.UpdateIPActionUsecase
	getFilters     *usecase.GetFiltersUsecase
	updateFilter   *usecase.UpdateFilterUsecase
}

func NewSecurityPolicyHandler(
	listRules *usecase.ListRulesUsecase,
	createRule *usecase.CreateRuleUsecase,
	toggleRule *usecase.ToggleRuleUsecase,
	updateRule *usecase.UpdateRuleUsecase,
	deleteRule *usecase.DeleteRuleUsecase,
	getIPConfig *usecase.GetIPConfigUsecase,
	updateIPAction *usecase.UpdateIPActionUsecase,
	getFilters *usecase.GetFiltersUsecase,
	updateFilter *usecase.UpdateFilterUsecase,
) *securityPolicyHandler {
	return &securityPolicyHandler{
		listRules:      listRules,
		createRule:     createRule,
		toggleRule:     toggleRule,
		updateRule:     updateRule,
		deleteRule:     deleteRule,
		getIPConfig:    getIPConfig,
		updateIPAction: updateIPAction,
		getFilters:     getFilters,
		updateFilter:   updateFilter,
	}
}

// --- WAF Rules ---

func (h *securityPolicyHandler) ListRules(w http.ResponseWriter, r *http.Request) error {
	rules, err := h.listRules.Execute(r.Context())
	if err != nil {
		return err
	}

	result := make([]ruleResponse, 0, len(rules))
	for _, rule := range rules {
		result = append(result, toRuleResponse(rule))
	}

	apphttp.WriteJSON(w, http.StatusOK, map[string]any{"rules": result})
	return nil
}

func (h *securityPolicyHandler) CreateRule(w http.ResponseWriter, r *http.Request) error {
	var req createRuleRequest
	if err := apphttp.DecodeJSON(r, &req); err != nil {
		return apperrors.NewBadRequest("invalid request body")
	}

	if req.Description == "" {
		return apperrors.NewBadRequest("description is required")
	}
	if req.Pattern == "" {
		return apperrors.NewBadRequest("pattern is required")
	}

	rule, err := h.createRule.Execute(r.Context(), usecase.CreateRuleInput{
		Description: req.Description,
		Pattern:     req.Pattern,
		Tags:        req.Tags,
		Impact:      req.Impact,
	})
	if err != nil {
		if err == domain.ErrInvalidPattern {
			return apperrors.NewBadRequest("pattern is not a valid regular expression")
		}
		return err
	}

	apphttp.WriteJSON(w, http.StatusCreated, toRuleResponse(*rule))
	return nil
}

func (h *securityPolicyHandler) UpdateRule(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")
	if id == "" {
		return apperrors.NewBadRequest("rule id is required")
	}

	var req updateRuleRequest
	if err := apphttp.DecodeJSON(r, &req); err != nil {
		return apperrors.NewBadRequest("invalid request body")
	}

	rule, err := h.updateRule.Execute(r.Context(), usecase.UpdateRuleInput{
		ID:          id,
		Description: req.Description,
		Pattern:     req.Pattern,
		Tags:        req.Tags,
		Impact:      req.Impact,
		Enabled:     req.Enabled,
	})
	if err != nil {
		if err == domain.ErrInvalidPattern {
			return apperrors.NewBadRequest("pattern is not a valid regular expression")
		}
		if err.Error() == "update waf rule: rule not found" {
			return apperrors.NewNotFound("waf rule not found")
		}
		return err
	}

	apphttp.WriteJSON(w, http.StatusOK, toRuleResponse(*rule))
	return nil
}

func (h *securityPolicyHandler) DeleteRule(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")
	if id == "" {
		return apperrors.NewBadRequest("rule id is required")
	}

	if err := h.deleteRule.Execute(r.Context(), id); err != nil {
		if err == domain.ErrRuleNotFound {
			return apperrors.NewNotFound("waf rule not found")
		}
		return err
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}

// --- IP Access Control ---

func (h *securityPolicyHandler) ListIPConfig(w http.ResponseWriter, r *http.Request) error {
	config, err := h.getIPConfig.Execute(r.Context())
	if err != nil {
		return err
	}

	result := make([]ipEntryResponse, 0, len(config.Records))
	for _, rec := range config.Records {
		result = append(result, ipEntryResponse{
			IP:     rec.IP,
			Status: string(rec.Status),
			Reason: rec.Reason,
		})
	}

	apphttp.WriteJSON(w, http.StatusOK, map[string]any{"entries": result})
	return nil
}

func (h *securityPolicyHandler) UpdateIPAction(w http.ResponseWriter, r *http.Request) error {
	var req updateIPActionRequest
	if err := apphttp.DecodeJSON(r, &req); err != nil {
		return apperrors.NewBadRequest("invalid request body")
	}

	if req.IP == "" {
		return apperrors.NewBadRequest("ip is required")
	}
	if req.Action == "" {
		return apperrors.NewBadRequest("action is required")
	}

	if err := h.updateIPAction.Execute(r.Context(), usecase.UpdateIPActionInput{
		IP:     req.IP,
		Action: req.Action,
		Reason: req.Reason,
	}); err != nil {
		return err
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}

// --- Filters ---

func (h *securityPolicyHandler) GetFilters(w http.ResponseWriter, r *http.Request) error {
	config, err := h.getFilters.Execute(r.Context())
	if err != nil {
		return err
	}

	apphttp.WriteJSON(w, http.StatusOK, filterConfigResponse{
		Keywords:        config.Keywords,
		WhitelistedURLs: config.WhitelistedURLs,
	})
	return nil
}

func (h *securityPolicyHandler) UpdateFilter(w http.ResponseWriter, r *http.Request) error {
	var req updateFilterRequest
	if err := apphttp.DecodeJSON(r, &req); err != nil {
		return apperrors.NewBadRequest("invalid request body")
	}

	if req.Value == "" {
		return apperrors.NewBadRequest("value is required")
	}

	if err := h.updateFilter.Execute(r.Context(), usecase.UpdateFilterInput{
		Type:   req.Type,
		Action: req.Action,
		Value:  req.Value,
	}); err != nil {
		return err
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}

// --- Private helpers ---

func toRuleResponse(rule domain.WAFRule) ruleResponse {
	return ruleResponse{
		ID:          rule.ID,
		Description: rule.Description,
		Pattern:     rule.Pattern,
		Tags:        rule.Tags,
		Impact:      rule.Impact,
		Enabled:     rule.Enabled,
		CreatedAt:   rule.CreatedAt,
		UpdatedAt:   rule.UpdatedAt,
	}
}

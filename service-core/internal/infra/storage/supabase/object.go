package supabase

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path"
	"service-core/internal/infra/storage"
	"strings"
)

func (p *SupabaseProvider) Upload(input storage.UploadInput) (*storage.ObjectResponse, error) {
	key := p.normalizeObjectKey(input.Key)
	if key == "" {
		return nil, fmt.Errorf("storage key is required")
	}

	body, err := io.ReadAll(input.File)
	if err != nil {
		return nil, fmt.Errorf("read upload body: %w", err)
	}

	req, err := http.NewRequest(
		http.MethodPost,
		p.ObjectURL,
		bytes.NewReader(body),
	)
	if err != nil {
		return nil, fmt.Errorf("build upload request: %w", err)
	}

	token := p.SupabaseConfig.ServiceRoleKey
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("apikey", token)
	req.Header.Set("x-upsert", "true")
	if input.ContentType != "" {
		req.Header.Set("Content-Type", input.ContentType)
	}

	resp, err := p.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("upload object to supabase: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
		return nil, p.newResponseError("upload object to supabase", resp)
	}

	return &storage.ObjectResponse{
		Key:         key,
		ContentType: input.ContentType,
	}, nil
}

func (p *SupabaseProvider) Delete(key string) error {
	key = p.normalizeObjectKey(key)
	if key == "" {
		return fmt.Errorf("storage key is required")
	}

	req, err := http.NewRequest(
		http.MethodDelete,
		p.ObjectURL,
		nil,
	)
	if err != nil {
		return fmt.Errorf("build delete request: %w", err)
	}

	token := p.SupabaseConfig.ServiceRoleKey
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("apikey", token)

	resp, err := p.Client.Do(req)
	if err != nil {
		return fmt.Errorf("delete object from supabase: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil
	}
	if resp.StatusCode >= http.StatusBadRequest {
		return p.newResponseError("delete object from supabase", resp)
	}

	return nil
}

func (p *SupabaseProvider) Exists(key string) (bool, error) {
	key = p.normalizeObjectKey(key)
	if key == "" {
		return false, fmt.Errorf("storage key is required")
	}

	req, err := http.NewRequest(
		http.MethodGet,
		p.ObjectURL,
		nil,
	)
	if err != nil {
		return false, fmt.Errorf("build exists request: %w", err)
	}

	token := p.SupabaseConfig.ServiceRoleKey
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("apikey", token)

	resp, err := p.Client.Do(req)
	if err != nil {
		return false, fmt.Errorf("check object in supabase: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return false, nil
	}
	if resp.StatusCode >= http.StatusBadRequest {
		return false, p.newResponseError("check object in supabase", resp)
	}

	return true, nil
}

func (p *SupabaseProvider) PublicURL(key string) string {
	return fmt.Sprintf(
		"%s/storage/v1/object/public/%s/%s",
		strings.TrimRight(p.SupabaseConfig.ProjectURL, "/"),
		p.StorageConfig.BucketName,
		key,
	)
}

func (p *SupabaseProvider) SignedURL(key string) (string, error) {
	payload, err := json.Marshal(map[string]int64{
		"expiresIn": int64(p.StorageConfig.SignedURLExpiry.Seconds()),
	})
	if err != nil {
		return "", fmt.Errorf("encode signed url request: %w", err)
	}

	req, err := http.NewRequest(
		http.MethodPost,
		p.ObjectURL,
		bytes.NewReader(payload),
	)
	if err != nil {
		return "", fmt.Errorf("build signed url request: %w", err)
	}

	token := p.SupabaseConfig.ServiceRoleKey
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("apikey", token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.Client.Do(req)
	if err != nil {
		return "", fmt.Errorf("generate signed url from supabase: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
		return "", p.newResponseError("generate signed url from supabase", resp)
	}

	var data supabaseSignedURLResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", fmt.Errorf("decode signed url response: %w", err)
	}
	if strings.TrimSpace(data.SignedURL) == "" {
		return "", fmt.Errorf("supabase signed url response is empty")
	}

	return strings.TrimPrefix(data.SignedURL, "/"), nil
}

func (p *SupabaseProvider) normalizeObjectKey(key string) string {
	key = strings.TrimSpace(key)
	key = strings.ReplaceAll(key, "\\", "/")
	key = strings.TrimPrefix(key, "/")

	cleaned := path.Clean(key)
	if cleaned == "." || cleaned == "" {
		return ""
	}

	if strings.HasPrefix(cleaned, "../") || cleaned == ".." {
		return ""
	}

	return cleaned
}

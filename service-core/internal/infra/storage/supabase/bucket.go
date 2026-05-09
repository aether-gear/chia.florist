package supabase

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

func (p *SupabaseProvider) EnsureBucket() error {
	log.Println("validate existing bucket...")

	req, err := http.NewRequest(
		http.MethodGet,
		p.BucketURL,
		nil,
	)
	if err != nil {
		return fmt.Errorf("build bucket validation request: %w", err)
	}

	token := p.SupabaseConfig.ServiceRoleKey
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("apikey", token)

	resp, err := p.Client.Do(req)
	if err != nil {
		return fmt.Errorf("validate supabase bucket: %w", err)
	}
	defer resp.Body.Close()

	var buckets []struct {
		ID string `json:"id"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&buckets); err != nil {
		return fmt.Errorf("decode bucket response: %w", err)
	}

	for _, bucket := range buckets {
		if bucket.ID == p.StorageConfig.BucketName {
			log.Println("bucket already exists")
			return nil
		}
	}
	log.Println("bucket not found, creating...")

	if err := p.createBucket(); err != nil {
		return fmt.Errorf("failed to create bucket: %w", err)
	}

	return nil
}

func (p *SupabaseProvider) createBucket() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	payload := map[string]any{
		"id":     p.StorageConfig.BucketName,
		"name":   p.StorageConfig.BucketName,
		"public": true,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal bucket payload: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		p.BucketURL,
		bytes.NewReader(body),
	)
	if err != nil {
		return fmt.Errorf("build create bucket request: %w", err)
	}

	token := p.SupabaseConfig.ServiceRoleKey
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("apikey", token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.Client.Do(req)
	if err != nil {
		return fmt.Errorf("create bucket failed: %w", err)
	}
	defer resp.Body.Close()

	return nil
}

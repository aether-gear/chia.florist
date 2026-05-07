package storage

type ProviderName string

const (
	ProviderLocal    ProviderName = "local"
	ProviderSupabase ProviderName = "supabase"
)

type Object struct {
	Key         string
	URL         string
	ContentType string
	Size        int64
}

func (p ProviderName) IsValid() bool {
	switch p {
	case ProviderLocal, ProviderSupabase:
		return true
	default:
		return false
	}
}

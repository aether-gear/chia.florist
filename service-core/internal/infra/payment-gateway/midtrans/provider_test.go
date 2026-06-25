package midtrans

import (
	"reflect"
	"testing"

	paymentgateway "service-core/internal/infra/payment-gateway"
)

func TestParseAmount(t *testing.T) {
	tests := []struct {
		name    string
		raw     string
		want    int64
		wantErr bool
	}{
		{
			name:    "empty string",
			raw:     "",
			want:    0,
			wantErr: false,
		},
		{
			name:    "whole number",
			raw:     "150000",
			want:    150000,
			wantErr: false,
		},
		{
			name:    "with decimal",
			raw:     "150000.00",
			want:    150000,
			wantErr: false,
		},
		{
			name:    "invalid format",
			raw:     "invalid",
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseAmount(tt.raw)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseAmount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("parseAmount() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMapNotificationStatus(t *testing.T) {
	tests := []struct {
		name        string
		txStatus    string
		fraudStatus string
		want        paymentgateway.NotificationStatus
	}{
		{
			name:        "capture challenge",
			txStatus:    "capture",
			fraudStatus: "challenge",
			want:        paymentgateway.NotificationStatusChallenge,
		},
		{
			name:        "capture accept",
			txStatus:    "capture",
			fraudStatus: "accept",
			want:        paymentgateway.NotificationStatusSettlement,
		},
		{
			name:        "settlement",
			txStatus:    "settlement",
			fraudStatus: "",
			want:        paymentgateway.NotificationStatusSettlement,
		},
		{
			name:        "pending",
			txStatus:    "pending",
			fraudStatus: "",
			want:        paymentgateway.NotificationStatusPending,
		},
		{
			name:        "deny",
			txStatus:    "deny",
			fraudStatus: "",
			want:        paymentgateway.NotificationStatusDeny,
		},
		{
			name:        "expire",
			txStatus:    "expire",
			fraudStatus: "",
			want:        paymentgateway.NotificationStatusExpire,
		},
		{
			name:        "cancel",
			txStatus:    "cancel",
			fraudStatus: "",
			want:        paymentgateway.NotificationStatusCancel,
		},
		{
			name:        "refund",
			txStatus:    "refund",
			fraudStatus: "",
			want:        paymentgateway.NotificationStatusRefund,
		},
		{
			name:        "partial_refund",
			txStatus:    "partial_refund",
			fraudStatus: "",
			want:        paymentgateway.NotificationStatusRefund,
		},
		{
			name:        "unknown",
			txStatus:    "unknown_status",
			fraudStatus: "",
			want:        paymentgateway.NotificationStatus("unknown_status"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mapNotificationStatus(tt.txStatus, tt.fraudStatus); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("mapNotificationStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}

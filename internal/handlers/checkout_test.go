package handlers

import (
	"testing"

	"github.com/andreashoj/order-system/internal/domain"
	"github.com/andreashoj/order-system/internal/pubsub"
	"github.com/andreashoj/order-system/internal/services"
)

// So far added interface to event_handler (to make it easier to mock here)

func TestHandleCheckout(t *testing.T) {
	type args struct {
		eventHandler    pubsub.EventHandler
		shoppingService *services.ShoppingService
		user            *domain.User
	}

	var mockEventHandler pubsub.EventHandler

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy-path",
			args: args{
				eventHandler: mockEventHandler,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := HandleCheckout(tt.args.eventHandler, tt.args.shoppingService, tt.args.user); (err != nil) != tt.wantErr {
				t.Errorf("HandleCheckout() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

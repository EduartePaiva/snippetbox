package mysql_test

import (
	"reflect"
	"testing"
	"time"

	"guthub.com/eduartepaiva/snippetbox/pkg/models"
	"guthub.com/eduartepaiva/snippetbox/pkg/models/mysql"
)

func TestUserModelGet(t *testing.T) {
	if testing.Short() {
		t.Skip("mysql: skipping integration test")
	}

	tests := []struct {
		name      string
		userID    int
		wantUser  *models.User
		wantError error
	}{
		{
			name:   "Valid ID",
			userID: 1,
			wantUser: &models.User{
				ID:        1,
				Name:      "Alice Jones",
				Email:     "alice@example.com",
				CreatedAt: time.Date(2026, 6, 6, 21, 21, 48, 0, time.UTC),
			},
			wantError: nil,
		},
		{
			name:      "Zero ID",
			userID:    0,
			wantUser:  nil,
			wantError: models.ErrNoRecord,
		},
		{
			name:      "Non-existent ID",
			userID:    2,
			wantUser:  nil,
			wantError: models.ErrNoRecord,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, teardown := newTestDB(t)
			defer teardown()

			m := &mysql.UserModel{DB: db}
			user, err := m.Get(tt.userID)
			if err != tt.wantError {
				t.Errorf("want %v; got %s", tt.wantError, err)
			}

			if !reflect.DeepEqual(tt.wantUser, user) {
				t.Errorf("\nwant %+v;\ngot %+v", tt.wantUser, user)
			}

		})
	}
}

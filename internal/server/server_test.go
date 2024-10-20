package server

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

type MockRedisClient struct {
	data map[string]string
}

func NewMockRedisClient() *MockRedisClient {
	return &MockRedisClient{
		data: make(map[string]string),
	}
}

func (m *MockRedisClient) SetPendingUser(ctx context.Context, username string, userID string, expiration time.Duration) error {
	m.data["user:"+username] = userID
	return nil
}

func (m *MockRedisClient) GetPendingUser(ctx context.Context, username string) (string, error) {
	if userID, exists := m.data["user:"+username]; exists {
		return userID, nil
	}

	return "", redis.Nil
}

//
//func TestSignUpHandler(t *testing.T) {
//	cfg, err := config.LoadConfig()
//	if err != nil {
//		t.Fatalf("Failed to load config: %v", err)
//	}
//
//	uc := usecase.New(db.New(), cacher.New(context.Background(), config.Config{}), &telegram.Bot{}, "123456")
//	h := handler.New()
//	mockRedisClient := NewMockRedisClient()
//	bot := &telegram.Bot{}
//	server := New(bot)
//
//	ts := httptest.NewServer(server.h)
//	defer ts.Close()
//
//	payload := map[string]string{
//		"username": "testuser",
//		"user_id":  "123456",
//	}
//	jsonPayload, err := json.Marshal(payload)
//	if err != nil {
//		t.Fatalf("Failed to marshal payload: %v", err)
//	}
//
//	resp, err := http.Post(ts.URL+"/sign-up", "application/json", bytes.NewBuffer(jsonPayload))
//	if err != nil {
//		t.Fatalf("Failed to send POST request: %v", err)
//	}
//	defer resp.Body.Close()
//
//	if resp.StatusCode != http.StatusOK {
//		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
//	}
//
//	expectedMessage := "Please start the bot by sending /start to complete your registration: @YourBotUsername"
//	var responseMessage string
//	if err := json.NewDecoder(resp.Body).Decode(&responseMessage); err != nil {
//		t.Fatalf("Failed to decode response body: %v", err)
//	}
//
//	if responseMessage != expectedMessage {
//		t.Errorf("Expected response message %q, got %q", expectedMessage, responseMessage)
//	}
//
//	// Check if the user is registered in the mock Redis
//	if userID, exists := mockRedisClient.data["user:testuser"]; !exists || userID != "123456" {
//		t.Errorf("Expected user ID %q, got %q", "123456", userID)
//	}
//}

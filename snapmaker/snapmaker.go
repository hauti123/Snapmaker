package snapmaker

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type Snapmaker struct {
	ipAdress         string
	port             int
	token            string
	ctx              context.Context
	cancel           context.CancelFunc
	waitChanMutex    sync.Mutex
	waitForNewStatus chan string

	triggerStatusRetrieval chan time.Time
	statusLoopInterval     time.Duration

	httpClient HttpClient
}

func NewSnapmaker(ipAddress string, apiToken string) *Snapmaker {
	ctx, cancel := context.WithCancel(context.Background())
	return &Snapmaker{
		ipAdress:               ipAddress,
		port:                   snapmakerApiPort,
		token:                  apiToken,
		ctx:                    ctx,
		cancel:                 cancel,
		statusLoopInterval:     statusLoopInterval,
		httpClient:             &http.Client{},
		triggerStatusRetrieval: make(chan time.Time),
	}
}

func (sm *Snapmaker) WithHttpClient(httpClient HttpClient) *Snapmaker {
	sm.httpClient = httpClient
	return sm
}

func (sm *Snapmaker) GetIpAdress() string {
	return sm.ipAdress
}

func (sm *Snapmaker) GetPort() int {
	return sm.port
}

func (sm *Snapmaker) GetApiToken() string {
	return sm.token
}

func (sm Snapmaker) buildApiUrl(api string) string {
	return fmt.Sprintf("http://%s:%d/api/v1/%s", sm.ipAdress, sm.port, api)
}

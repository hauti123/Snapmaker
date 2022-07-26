package snapmaker

import (
	"time"
)

const snapmakerDiscoveryPort = 20054
const snapmakerDiscoveryPayload = "discover"
const snapmakerApiPort = 8080
const statusLoopTicker = 2 * time.Second

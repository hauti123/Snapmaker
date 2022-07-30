package snapmaker

import (
	"time"
)

const snapmakerDiscoveryPort = 20054
const snapmakerDiscoveryPayload = "discover"
const snapmakerApiPort = 8080
const statusLoopInterval = 2 * time.Second

const apiConnect = "connect"
const apiPrinterStatus = "status"
const apiGcodeUpload = "upload"
const apiEnclosureStatus = "enclosure"

const multipartBoundary = "----------------------------268923783128719097072428"

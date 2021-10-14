package config

import (
	"github.com/spf13/viper"
)

// nolint
// NewMockConfig creates a new mock instance.
func NewMockConfig() *viper.Viper {
	mock := viper.New()
	mock.Set("SERVICE_NAME", "FC Latency Manage")
	mock.Set("FILECOIN_NODE_URL", "https://node.glif.io/space07/lotus/rpc/v0")
	mock.Set("FILECOIN_BLOCKS_OFFSET", 10)
	mock.Set("DB_CONNECTION", "data/database.db")
	mock.Set("RIPE_API_KEY", "changeme")
	mock.Set("RIPE_ONE_OFF", true)
	mock.Set("RIPE_PING_INTERVAL", 60)
	mock.Set("RIPE_PING_RUNNING_TIME", 300)
	mock.Set("RIPE_PACKETS", 1)
	mock.Set("RIPE_REQUESTED_PROBES", 2)
	mock.Set("RIPE_LOCATION_RANGE_INIT", 0.01)
	mock.Set("RIPE_LOCATION_RANGE_MAX", 1000)
	mock.Set("RIPE_PROBES_PER_AIRPORT", 3)
	mock.Set("NEAREST_AIRPORTS", 50)
	mock.Set("CONSTANT_AIRPORTS", "constants/airport-codes.json")
	mock.Set("SQL_DEBUG", false)
	mock.Set("CRON_SCHEDULE_CREATE_MESURES", "0 0 0 * * *")
	mock.Set("CRON_SCHEDULE_IMPORT_MESURES", "0 */15 * * * *")
	mock.Set("WEBHOOK_NOTIFY_URLS", "http://webcode.me,https://postman-echo.com/post")
	return mock
}

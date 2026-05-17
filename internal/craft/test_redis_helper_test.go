package craft

import (
	"fmt"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
)

type testRedisServer struct {
	server *miniredis.Miniredis
}

func setupTestRedis(tb testing.TB) *testRedisServer {
	tb.Helper()

	server, err := miniredis.Run()
	if err != nil {
		tb.Fatalf("could not start miniredis: %v", err)
	}

	tb.Setenv("FC_REDIS_URI", fmt.Sprintf("redis://%s", server.Addr()))
	tb.Cleanup(server.Close)

	return &testRedisServer{server: server}
}

func (r *testRedisServer) FlushAll(tb testing.TB) {
	tb.Helper()
	r.server.FlushAll()
}

func (r *testRedisServer) SetString(tb testing.TB, key string, value string, ttl time.Duration) {
	tb.Helper()
	if err := r.server.Set(key, value); err != nil {
		tb.Fatalf("could not seed redis key %q: %v", key, err)
	}
	if ttl > 0 {
		r.server.SetTTL(key, ttl)
	}
}

package cache

import (
	"time"

	"github.com/patrickmn/go-cache"
)

var PendingLoginCache = cache.New(5*time.Minute, 10*time.Minute)
var SessionCache = cache.New(5*time.Minute, 10*time.Minute)

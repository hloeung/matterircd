package slackclient

import (
	"testing"
	"time"

	"github.com/slack-go/slack"
	"github.com/stretchr/testify/assert"
)

func TestStaleFilteredPresence(t *testing.T) {
	now := time.Now().Unix()

	// Away entry with unverified timestamp omits LastActivity.
	unverified := presenceEntry{
		presence: slack.UserPresence{
			Presence:     presenceAway,
			LastActivity: slack.JSONTime(now - 3600),
		},
	}
	assert.Zero(t, staleFilteredPresence(unverified, now).LastActivity)

	// Away entry verified > 5 minutes ago omits LastActivity.
	stale := presenceEntry{
		presence: slack.UserPresence{
			Presence:     presenceAway,
			LastActivity: slack.JSONTime(now - 3600),
		},
		verifiedAt: now - int64(presenceStaleTTL.Seconds()) - 1,
	}
	assert.Zero(t, staleFilteredPresence(stale, now).LastActivity)

	// Away entry verified within 5 minutes preserves LastActivity.
	fresh := presenceEntry{
		presence: slack.UserPresence{
			Presence:     presenceAway,
			LastActivity: slack.JSONTime(now - 3600),
		},
		verifiedAt: now - 60,
	}
	assert.Equal(t, slack.JSONTime(now-3600), staleFilteredPresence(fresh, now).LastActivity)

	// Active entry preserves LastActivity.
	active := presenceEntry{
		presence: slack.UserPresence{
			Presence:     presenceActive,
			Online:       true,
			LastActivity: slack.JSONTime(now - 10),
		},
	}
	assert.Equal(t, slack.JSONTime(now-10), staleFilteredPresence(active, now).LastActivity)
}

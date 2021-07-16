package message

import (
	"time"

	"github.com/target/goalert/notification"
)

// Throttle represents the throttled messages for a queue.
type Throttle struct {
	cfg      ThrottleConfig
	ignoreID bool
	now      time.Time

	first    map[ThrottleItem]time.Time
	count    map[ThrottleItem]int
	cooldown map[notification.Dest]bool
}

// ThrottleItem represents the messages being throttled.
type ThrottleItem struct {
	Dest      notification.Dest
	BucketDur time.Duration
}

// ThrottleConfig provides ThrottleRules for a given message.
type ThrottleConfig interface {
	Rules(Message) []ThrottleRule
	MaxDuration() time.Duration
}

func maxThrottleDuration(cfgs ...ThrottleConfig) time.Duration {
	var max time.Duration
	for _, cfg := range cfgs {
		dur := cfg.MaxDuration()
		if dur > max {
			max = dur
		}
	}
	return max
}

// NewThrottle creates a new Throttle used to manage outgoing messages in a queue.
func NewThrottle(cfg ThrottleConfig, now time.Time, ignoreID bool) *Throttle {
	return &Throttle{
		cfg:      cfg,
		now:      now,
		ignoreID: ignoreID,

		first:    make(map[ThrottleItem]time.Time),
		count:    make(map[ThrottleItem]int),
		cooldown: make(map[notification.Dest]bool),
	}
}

// Record keeps track of the outgoing messages being throttled in a queue.
func (tr *Throttle) Record(msg Message) {
	if tr.ignoreID {
		msg.Dest.ID = ""
	}
	msg.Dest.Value = ""

	since := tr.now.Sub(msg.SentAt)
	rules := tr.cfg.Rules(msg)
	for i, rule := range rules {
		if since >= rule.Per {
			continue
		}
		key := ThrottleItem{Dest: msg.Dest, BucketDur: rule.Per}
		tr.count[key]++
		count := tr.count[key]
		if tr.first[key].IsZero() || msg.SentAt.Before(tr.first[key]) {
			tr.first[key] = msg.SentAt
		}

		if count >= rule.Count {
			tr.cooldown[msg.Dest] = true
			continue
		}

		if !rule.Smooth {
			continue
		}

		// flat rate
		var prevRule ThrottleRule
		if i > 0 {
			prevRule = rules[i-1]
		}

		if count < prevRule.Count || count == 0 {
			// allow prev rule in entirety
			continue
		}

		// spread remainder evenly
		count -= prevRule.Count
		elapsed := tr.now.Sub(tr.first[key]) - prevRule.Per
		per := rule.Per - prevRule.Per

		if count > int(elapsed*time.Duration(rule.Count-prevRule.Count)/per) {
			tr.cooldown[msg.Dest] = true
		}
	}
}

// InCooldown returns true or false depending on the cooldown state of a throttled message.
func (tr *Throttle) InCooldown(msg Message) bool {
	if tr.ignoreID {
		msg.Dest.ID = ""
	}
	msg.Dest.Value = ""

	return tr.cooldown[msg.Dest]
}

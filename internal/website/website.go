package website

import (
	"errors"
	"github.com/utkuozdemir/website-watcher/internal/notification"
	"time"
)

const (
	defaultDelay             = 30 * time.Second
	defaultNotificationDelay = 5 * time.Minute
	defaultSuccessThreshold  = 2
)

type Website interface {
	Name() string
	URL() string
	Condition() func(pageSource *string, status int) bool
	Delay() time.Duration
	NotificationDelay() time.Duration
	SuccessThreshold() int
	ConsecutiveSuccessCount() int
	IncrementConsecutiveSuccessCount()
	ResetConsecutiveSuccessCount()
	LastNotificationTime() time.Time
	UpdateLastNotificationTime(t time.Time)
	NotificationHandler() notification.Handler
}

type website struct {
	name                    string
	url                     string
	condition               func(pageSource *string, status int) bool
	delay                   time.Duration
	notificationDelay       time.Duration
	successThreshold        int
	lastNotificationTime    time.Time
	notificationHandler     notification.Handler
	consecutiveSuccessCount int
}

func (w *website) Name() string {
	return w.name
}

func (w *website) URL() string {
	return w.url
}

func (w *website) Condition() func(pageSource *string, status int) bool {
	return w.condition
}

func (w *website) Delay() time.Duration {
	return w.delay
}

func (w *website) NotificationDelay() time.Duration {
	return w.delay
}

func (w *website) SuccessThreshold() int {
	return w.successThreshold
}

func (w *website) LastNotificationTime() time.Time {
	return w.lastNotificationTime
}

func (w *website) UpdateLastNotificationTime(t time.Time) {
	w.lastNotificationTime = t
}

func (w *website) NotificationHandler() notification.Handler {
	return w.notificationHandler
}

func (w *website) ConsecutiveSuccessCount() int {
	return w.consecutiveSuccessCount
}

func (w *website) IncrementConsecutiveSuccessCount() {
	w.consecutiveSuccessCount++
}

func (w *website) ResetConsecutiveSuccessCount() {
	w.consecutiveSuccessCount = 0
}

func New(name string, url string, condition func(pageSource *string, status int) bool,
	delay time.Duration, notificationDelay time.Duration, successThreshold int,
	notificationHandler notification.Handler) (Website, error) {
	if name == "" {
		return nil, errors.New("name is required")
	}

	if url == "" {
		return nil, errors.New("url is required")
	}

	if condition == nil {
		return nil, errors.New("condition is required")
	}

	if delay == 0 {
		return nil, errors.New("delay is required")
	}

	if notificationDelay == 0 {
		return nil, errors.New("notificationDelay is required")
	}

	if successThreshold == 0 {
		return nil, errors.New("successThreshold is required")
	}

	return &website{
		name:                name,
		url:                 url,
		condition:           condition,
		delay:               delay,
		notificationDelay:   notificationDelay,
		successThreshold:    successThreshold,
		notificationHandler: notificationHandler,
	}, nil
}

func NewWithDefaults(name string, url string,
	condition func(pageSource *string, status int) bool,
	notificationHandler notification.Handler) (Website, error) {
	return New(name, url, condition, defaultDelay, defaultNotificationDelay,
		defaultSuccessThreshold, notificationHandler)
}

func MustNewWithDefaults(name string, url string,
	condition func(pageSource *string, status int) bool,
	notificationHandler notification.Handler) Website {
	w, err := NewWithDefaults(name, url, condition, notificationHandler)
	if err != nil {
		panic("Error on creating website: " + err.Error())
	}
	return w
}

package observability

import (
	"FeedCraft/internal/dao"
	"FeedCraft/internal/model"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

const (
	ResourceTypeRecipe = "recipe"
	ResourceTypeTopic  = "topic"

	TriggerUserRequest      = "user_request"
	TriggerPreheating       = "preheating"
	TriggerTopicAggregation = "topic_aggregation"
	TriggerSystem           = "system"

	ErrorKindTimeout                = "timeout"
	ErrorKindNetwork                = "network"
	ErrorKindHTTPStatus             = "http_status"
	ErrorKindParse                  = "parse"
	ErrorKindEmptyFeed              = "empty_feed"
	ErrorKindCraft                  = "craft"
	ErrorKindDependency             = "dependency"
	ErrorKindConfig                 = "config"
	ErrorKindUpstreamPartialFailure = "upstream_partial_failure"
	ErrorKindUnknown                = "unknown"
)

const (
	DegradedFailureThreshold = 3
	PausedFailureThreshold   = 5
)

type ExecutionEvent struct {
	ResourceType string
	ResourceID   string
	ResourceName string
	Trigger      string
	Status       string
	ErrorKind    string
	Message      string
	Details      map[string]any
	RequestID    string
	Duration     time.Duration
	OccurredAt   time.Time
}

type Service struct {
	db      *gorm.DB
	events  chan ExecutionEvent
	done    chan struct{}
	wg      sync.WaitGroup
	closed  bool
	closeMu sync.Mutex
}

var globalService *Service

func Init(db *gorm.DB) {
	if db == nil {
		return
	}
	globalService = newService(db)
}

func Shutdown() {
	if globalService != nil {
		svc := globalService
		globalService = nil
		svc.Close()
	}
}

func Global() *Service {
	return globalService
}

func newService(db *gorm.DB) *Service {
	s := &Service{
		db:     db,
		events: make(chan ExecutionEvent, 256),
		done:   make(chan struct{}),
	}
	s.wg.Add(1)
	go s.run()
	return s
}

func (s *Service) Close() {
	s.closeMu.Lock()
	if s.closed {
		s.closeMu.Unlock()
		return
	}
	s.closed = true
	close(s.done)
	s.closeMu.Unlock()
	s.wg.Wait()
}

func Report(event ExecutionEvent) {
	if globalService == nil {
		return
	}
	globalService.Report(event)
}

func (s *Service) Report(event ExecutionEvent) {
	if event.OccurredAt.IsZero() {
		event.OccurredAt = time.Now()
	}
	select {
	case <-s.done:
		return
	case s.events <- event:
	default:
		logrus.Warnf("observability event dropped for %s/%s status=%s", event.ResourceType, event.ResourceID, event.Status)
	}
}

func (s *Service) run() {
	defer s.wg.Done()
	for {
		select {
		case event := <-s.events:
			if err := s.persistEvent(event); err != nil {
				logrus.Errorf("persist observability event failed for %s/%s: %v", event.ResourceType, event.ResourceID, err)
			}
		case <-s.done:
			for {
				select {
				case event := <-s.events:
					if err := s.persistEvent(event); err != nil {
						logrus.Errorf("persist observability event failed for %s/%s: %v", event.ResourceType, event.ResourceID, err)
					}
				default:
					return
				}
			}
		}
	}
}

func (s *Service) persistEvent(event ExecutionEvent) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		detailsJSON := ""
		if len(event.Details) > 0 {
			raw, err := json.Marshal(event.Details)
			if err != nil {
				return err
			}
			detailsJSON = string(raw)
		}

		if err := dao.CreateExecutionLog(tx, &dao.ExecutionLog{
			ResourceType: event.ResourceType,
			ResourceID:   event.ResourceID,
			ResourceName: event.ResourceName,
			Trigger:      event.Trigger,
			Status:       event.Status,
			ErrorKind:    event.ErrorKind,
			Message:      event.Message,
			DetailsJSON:  detailsJSON,
			RequestID:    event.RequestID,
			DurationMS:   event.Duration.Milliseconds(),
			BaseModelWithoutPK: dao.BaseModelWithoutPK{
				CreatedAt: event.OccurredAt,
				UpdatedAt: event.OccurredAt,
			},
		}); err != nil {
			return err
		}

		health, err := dao.GetResourceHealth(tx, event.ResourceType, event.ResourceID)
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}
		if err == gorm.ErrRecordNotFound {
			health = &dao.ResourceHealth{
				ResourceType:  event.ResourceType,
				ResourceID:    event.ResourceID,
				ResourceName:  event.ResourceName,
				CurrentStatus: dao.ResourceStatusHealthy,
			}
		}

		previousStatus := health.CurrentStatus
		applyEventToHealth(health, event)

		if err := dao.UpsertResourceHealth(tx, health); err != nil {
			return err
		}

		if notification := buildNotification(previousStatus, health, event); notification != nil {
			if err := dao.CreateSystemNotification(tx, notification); err != nil {
				if strings.Contains(strings.ToLower(err.Error()), "unique") {
					return nil
				}
				return err
			}
		}

		return nil
	})
}

func applyEventToHealth(health *dao.ResourceHealth, event ExecutionEvent) {
	health.ResourceName = event.ResourceName
	now := event.OccurredAt
	health.UpdatedAt = now
	if health.CreatedAt.IsZero() {
		health.CreatedAt = now
	}

	switch event.Status {
	case dao.ExecutionStatusSuccess:
		health.CurrentStatus = dao.ResourceStatusHealthy
		health.ConsecutiveFailures = 0
		health.LastSuccessAt = &now
		health.LastErrorKind = ""
		health.LastErrorMessage = ""
		health.PausedAt = nil
		health.PausedReason = ""
	case dao.ExecutionStatusPartialSuccess:
		health.CurrentStatus = dao.ResourceStatusDegraded
		health.ConsecutiveFailures = 0
		health.LastSuccessAt = &now
		health.LastErrorKind = event.ErrorKind
		health.LastErrorMessage = event.Message
	case dao.ExecutionStatusFailure:
		health.ConsecutiveFailures++
		health.LastFailureAt = &now
		health.LastErrorKind = event.ErrorKind
		health.LastErrorMessage = event.Message
		if health.ConsecutiveFailures >= PausedFailureThreshold {
			health.CurrentStatus = dao.ResourceStatusPaused
			health.PausedAt = &now
			health.PausedReason = event.Message
		} else if health.ConsecutiveFailures >= DegradedFailureThreshold {
			health.CurrentStatus = dao.ResourceStatusDegraded
		}
	case dao.ExecutionStatusPausedSkip:
		health.CurrentStatus = dao.ResourceStatusPaused
		if health.PausedAt == nil {
			health.PausedAt = &now
		}
		if event.Message != "" {
			health.PausedReason = event.Message
		}
	}
}

func buildNotification(previousStatus string, health *dao.ResourceHealth, event ExecutionEvent) *dao.SystemNotification {
	if health == nil {
		return nil
	}
	eventType := ""
	switch {
	case event.Status == dao.ExecutionStatusPartialSuccess && previousStatus != dao.ResourceStatusDegraded:
		eventType = "resource_degraded"
	case event.Status == dao.ExecutionStatusFailure && health.CurrentStatus == dao.ResourceStatusDegraded && previousStatus != dao.ResourceStatusDegraded:
		eventType = "resource_degraded"
	case event.Status == dao.ExecutionStatusFailure && health.CurrentStatus == dao.ResourceStatusPaused && previousStatus != dao.ResourceStatusPaused:
		eventType = "resource_paused"
	case event.Trigger == TriggerSystem && event.Status == dao.ExecutionStatusSuccess:
		eventType = "resource_resumed"
	default:
		return nil
	}

	resourceName := health.ResourceName
	if resourceName == "" {
		resourceName = health.ResourceID
	}
	title := fmt.Sprintf("[%s] %s %s", strings.ToUpper(health.ResourceType), resourceName, strings.ReplaceAll(eventType, "_", " "))
	content := fmt.Sprintf("Resource %s/%s changed to %s. %s", health.ResourceType, health.ResourceID, health.CurrentStatus, event.Message)
	return &dao.SystemNotification{
		ResourceType: health.ResourceType,
		ResourceID:   health.ResourceID,
		EventType:    eventType,
		Title:        title,
		Content:      content,
		StatusAfter:  health.CurrentStatus,
		DedupeKey:    fmt.Sprintf("%s:%s:%s:%d", health.ResourceType, health.ResourceID, eventType, event.OccurredAt.Unix()),
	}
}

func ClassifyError(err error) string {
	if err == nil {
		return ""
	}
	message := strings.ToLower(err.Error())
	switch {
	case strings.Contains(message, "timeout"), strings.Contains(message, "deadline exceeded"):
		return ErrorKindTimeout
	case strings.Contains(message, "connection refused"), strings.Contains(message, "no such host"), strings.Contains(message, "network"):
		return ErrorKindNetwork
	case strings.Contains(message, "status"):
		return ErrorKindHTTPStatus
	case strings.Contains(message, "parse"), strings.Contains(message, "selector"), strings.Contains(message, "json"):
		return ErrorKindParse
	case strings.Contains(message, "craft"), strings.Contains(message, "process feed"):
		return ErrorKindCraft
	case strings.Contains(message, "config"), strings.Contains(message, "invalid source"):
		return ErrorKindConfig
	case strings.Contains(message, "dependency"):
		return ErrorKindDependency
	case strings.Contains(message, "empty feed"), strings.Contains(message, "no items"):
		return ErrorKindEmptyFeed
	default:
		return ErrorKindUnknown
	}
}

func ShouldSkipRecipe(resourceID string) bool {
	if globalService == nil {
		return false
	}
	health, err := dao.GetResourceHealth(globalService.db, ResourceTypeRecipe, resourceID)
	if err != nil {
		return false
	}
	return health.CurrentStatus == dao.ResourceStatusPaused
}

func ResumeResource(resourceType string, resourceID string) error {
	if globalService == nil {
		return nil
	}
	return globalService.db.Transaction(func(tx *gorm.DB) error {
		health, err := dao.GetResourceHealth(tx, resourceType, resourceID)
		if err != nil {
			return err
		}
		now := time.Now()
		health.CurrentStatus = dao.ResourceStatusHealthy
		health.ConsecutiveFailures = 0
		health.LastErrorKind = ""
		health.LastErrorMessage = ""
		health.PausedAt = nil
		health.PausedReason = ""
		health.UpdatedAt = now
		if err := dao.UpsertResourceHealth(tx, health); err != nil {
			return err
		}
		svc := globalService
		if svc != nil {
			svc.Report(ExecutionEvent{
				ResourceType: resourceType,
				ResourceID:   resourceID,
				ResourceName: health.ResourceName,
				Trigger:      TriggerSystem,
				Status:       dao.ExecutionStatusSuccess,
				Message:      "resource resumed by admin",
				OccurredAt:   now,
			})
		}
		return nil
	})
}

func BuildNotificationFeed(ctx context.Context) (*model.CraftFeed, error) {
	_ = ctx
	if globalService == nil {
		return &model.CraftFeed{
			Id:          "feedcraft://system/notifications",
			Title:       "FeedCraft System Notifications",
			Description: "Built-in FeedCraft system alerts and state transitions",
			Link:        "/system/notifications/rss",
			Updated:     time.Now(),
			Created:     time.Now(),
		}, nil
	}
	items, err := dao.ListSystemNotifications(globalService.db, globalService.db.Model(&dao.SystemNotification{}).Limit(100))
	if err != nil {
		return nil, err
	}

	feed := &model.CraftFeed{
		Id:          "feedcraft://system/notifications",
		Title:       "FeedCraft System Notifications",
		Description: "Built-in FeedCraft system alerts and state transitions",
		Link:        "/system/notifications/rss",
		Updated:     time.Now(),
		Created:     time.Now(),
	}
	for _, item := range items {
		feed.Articles = append(feed.Articles, &model.CraftArticle{
			Id:          fmt.Sprintf("system-notification:%d", item.ID),
			Title:       item.Title,
			Link:        "/system/notifications/rss",
			Description: item.Content,
			Content:     item.Content,
			Created:     item.CreatedAt,
			Updated:     item.UpdatedAt,
			AuthorName:  "FeedCraft",
		})
	}
	return feed, nil
}

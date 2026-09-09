package task

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

func NewRedisCachingRepository(rcl *redis.Client, repo Repository, ttl time.Duration) redisCaching {
	return redisCaching{
		rcl:  rcl,
		repo: repo,
	}
}

type redisCaching struct {
	rcl  *redis.Client
	repo Repository
	ttl  time.Duration
}

func (r redisCaching) key(id int64) string {
	return fmt.Sprintf("task:%d", id)
}

// it should be atomic and cache must be sync with database
// each cache manipulation must commit with database transacation
// but i'm too lazy to do it now and keep it simple :)
func (r redisCaching) Total(ctx context.Context) (int64, error) {
	return r.repo.Total(ctx)
}

func (r redisCaching) Get(ctx context.Context, id int64) (Task, error) {
	b, err := r.rcl.Get(ctx, r.key(id)).Bytes()
	if err != nil {
		// redis is down or cache missed because of ttl
		if !errors.Is(err, redis.Nil) {
			log.Printf("redis error: %v\n", err)
		}
		t, err := r.repo.Get(ctx, id)
		_ = r.set(ctx, t)
		return t, err
	}
	var t Task
	if err = json.Unmarshal(b, &t); err != nil {
		return Task{}, err
	}
	return t, nil
}

func (r redisCaching) set(ctx context.Context, t Task) error {
	b, err := json.Marshal(t)
	if err != nil {
		return err
	}
	return r.rcl.Set(ctx, r.key(t.ID), b, r.ttl).Err()
}

func (r redisCaching) Create(ctx context.Context, t Task) (int64, error) {
	id, err := r.repo.Create(ctx, t)
	t.ID = id
	_ = r.set(ctx, t)
	return id, err
}

func (r redisCaching) ChangeStatus(ctx context.Context, id int64, newStatus Status) error {
	err := r.repo.ChangeStatus(ctx, id, newStatus)
	if err != nil {
		return err
	}
	//replace data
	_ = r.rcl.Del(ctx, r.key(id)).Err()
	t, _ := r.Get(ctx, id)
	_ = r.set(ctx, t)
	return nil
}

func (r redisCaching) ChangeAssignee(ctx context.Context, id int64, assigneeID int) error {
	err := r.repo.ChangeAssignee(ctx, id, assigneeID)
	if err != nil {
		return err
	}
	//replace data
	_ = r.rcl.Del(ctx, r.key(id)).Err()
	t, _ := r.Get(ctx, id)
	_ = r.set(ctx, t)
	return nil
}

func (r redisCaching) ChangePriority(ctx context.Context, id int64, newPriority Priority) error {
	err := r.repo.ChangePriority(ctx, id, newPriority)
	if err != nil {
		return err
	}
	//replace data
	_ = r.rcl.Del(ctx, r.key(id)).Err()
	t, _ := r.Get(ctx, id)
	_ = r.set(ctx, t)
	return nil
}

func (r redisCaching) List(ctx context.Context, req ListRequest) ([]Task, int64, error) {
	return r.repo.List(ctx, req)
}

func (r redisCaching) Delete(ctx context.Context, id int64) error {
	err := r.Delete(ctx, id)
	if err != nil {
		return err
	}
	_ = r.rcl.Del(ctx, r.key(id)).Err()
	return nil
}

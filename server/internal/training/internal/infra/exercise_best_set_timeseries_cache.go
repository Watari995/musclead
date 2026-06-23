package traininginfra

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	trainingdomain "github.com/Watari995/musclead/internal/training/internal/domain"
	"github.com/Watari995/musclead/internal/valueobject"
	"github.com/redis/go-redis/v9"
)

// Redis実装。weight_timeseries_cache.go と同じ ZSet + Hash 2層構造を採用。
//
// # なぜ weight は ZSet + Hash 2層構造なのか
//
// weight は個別レコードの削除・更新があるため、
// 「特定の1件だけ消す（ZRem + HDel）」操作を効率よく行う必要がある。
// そのため ZSet をインデックス、Hash を実データとして分離している。
//
// # この実装はシンプルな1RTT構造でも良かった
//
// exercise best-set timeseries は個別レコードの更新・削除を行わず、
// Evict で全削除するだけのため、ZSet のメンバーに JSON を直接埋め込む
// 1RTT 構造（ZRangeByScore だけで完結）でも十分だった。
// weight との一貫性を優先して2層構造にしている。

type redisExerciseBestSetTimeseriesCache struct {
	client *redis.Client
	ttl    time.Duration
}

func NewRedisExerciseBestSetTimeseriesCache(client *redis.Client, ttl time.Duration) trainingdomain.ExerciseBestSetTimeseriesCache {
	return &redisExerciseBestSetTimeseriesCache{client: client, ttl: ttl}
}

func (c *redisExerciseBestSetTimeseriesCache) idxKey(userID valueobject.UserID, exerciseID valueobject.ExerciseID) string {
	return fmt.Sprintf("exercises:best-set-timeseries:%s:%s:idx", userID.Value(), exerciseID.Value())
}

func (c *redisExerciseBestSetTimeseriesCache) dataKey(userID valueobject.UserID, exerciseID valueobject.ExerciseID) string {
	return fmt.Sprintf("exercises:best-set-timeseries:%s:%s:data", userID.Value(), exerciseID.Value())
}

type bestSetCacheRecord struct {
	WeightKg    string    `json:"weight_kg"`
	Reps        int32     `json:"reps"`
	PerformedAt time.Time `json:"performed_at"`
	TrainingID  string    `json:"training_id"`
	ExerciseID  string    `json:"exercise_id"`
}

func (c *redisExerciseBestSetTimeseriesCache) FindByPeriod(ctx context.Context, userID valueobject.UserID, exerciseID valueobject.ExerciseID, from, to time.Time) ([]*trainingdomain.BestSetView, bool, error) {
	// 1. ZRangeArgs で期間内の trainingID 一覧を取得
	ids, err := c.client.ZRangeArgs(ctx, redis.ZRangeArgs{
		Key:     c.idxKey(userID, exerciseID),
		Start:   from.Unix(),
		Stop:    to.Unix(),
		ByScore: true,
	}).Result()
	if err != nil {
		return nil, false, err
	}
	if len(ids) == 0 {
		return nil, false, nil
	}
	// 2. HMGet で JSON を一括取得
	vals, err := c.client.HMGet(ctx, c.dataKey(userID, exerciseID), ids...).Result()
	if err != nil {
		return nil, false, err
	}
	// 3. JSON を BestSetView に変換して返す
	bestSets := make([]*trainingdomain.BestSetView, 0, len(ids))
	for _, v := range vals {
		s, ok := v.(string)
		if !ok {
			return nil, false, nil
		}
		bestSet, err := decodeBestSetCacheRecord(s)
		if err != nil {
			return nil, false, err
		}
		bestSets = append(bestSets, bestSet)
	}
	if len(bestSets) == 0 {
		return nil, false, nil
	}
	return bestSets, true, nil
}

func (c *redisExerciseBestSetTimeseriesCache) Save(ctx context.Context, userID valueobject.UserID, bestSet *trainingdomain.BestSetView) error {
	jsonStr, err := encodeBestSetCacheRecord(bestSet)
	if err != nil {
		return err
	}
	idx := c.idxKey(userID, bestSet.ExerciseID)
	data := c.dataKey(userID, bestSet.ExerciseID)

	// MULTI/EXEC でZADDとHSETをatomicに実行する
	_, err = c.client.TxPipelined(ctx, func(p redis.Pipeliner) error {
		// id を SORTのためのスコア付きで保存
		p.ZAdd(ctx, idx, redis.Z{Score: float64(bestSet.PerformedAt.Unix()), Member: bestSet.TrainingID.Value()})
		// 実データをハッシュで保存
		p.HSet(ctx, data, bestSet.TrainingID.Value(), jsonStr)
		p.Expire(ctx, idx, c.ttl)
		p.Expire(ctx, data, c.ttl)
		return nil
	})
	return err
}

func (c *redisExerciseBestSetTimeseriesCache) Evict(ctx context.Context, userID valueobject.UserID, exerciseID valueobject.ExerciseID) error {
	idx := c.idxKey(userID, exerciseID)
	data := c.dataKey(userID, exerciseID)
	_, err := c.client.TxPipelined(ctx, func(p redis.Pipeliner) error {
		p.Del(ctx, idx)
		p.Del(ctx, data)
		return nil
	})
	return err
}

// ---helper functions---
func encodeBestSetCacheRecord(b *trainingdomain.BestSetView) (string, error) {
	record := bestSetCacheRecord{
		WeightKg:    b.WeightKg.String(),
		Reps:        int32(b.Reps.Value()),
		PerformedAt: b.PerformedAt,
		TrainingID:  b.TrainingID.Value(),
		ExerciseID:  b.ExerciseID.Value(),
	}
	d, err := json.Marshal(record)
	if err != nil {
		return "", err
	}
	return string(d), nil
}

func decodeBestSetCacheRecord(data string) (*trainingdomain.BestSetView, error) {
	var record bestSetCacheRecord
	err := json.Unmarshal([]byte(data), &record)
	if err != nil {
		return nil, err
	}
	weightKg, err := valueobject.NewNonNegativeDecimalFromString(record.WeightKg)
	if err != nil {
		return nil, err
	}
	reps, err := valueobject.NewNonNegativeInt(int(record.Reps))
	if err != nil {
		return nil, err
	}
	trainingID, err := valueobject.NewPrimaryIDFromString[valueobject.TrainingID](record.TrainingID)
	if err != nil {
		return nil, err
	}
	exerciseID, err := valueobject.NewPrimaryIDFromString[valueobject.ExerciseID](record.ExerciseID)
	if err != nil {
		return nil, err
	}
	return &trainingdomain.BestSetView{
		WeightKg:    *weightKg,
		Reps:        *reps,
		PerformedAt: record.PerformedAt,
		TrainingID:  *trainingID,
		ExerciseID:  *exerciseID,
	}, nil
}

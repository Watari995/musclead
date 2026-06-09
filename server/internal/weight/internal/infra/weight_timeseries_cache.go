package weightinfra

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Watari995/musclead/internal/valueobject"
	weightdomain "github.com/Watari995/musclead/internal/weight/internal/domain"
	"github.com/redis/go-redis/v9"
)

type RedisWeightTimeseriesCache struct {
	client *redis.Client
	ttl    time.Duration
}

func NewRedisWeightTimeseriesCache(client *redis.Client, ttl time.Duration) *RedisWeightTimeseriesCache {
	return &RedisWeightTimeseriesCache{client: client, ttl: ttl}
}

func (c *RedisWeightTimeseriesCache) dataKey(userID valueobject.UserID) string {
	return fmt.Sprintf("weights:timeseries:%s:data", userID.Value())
}

func (c *RedisWeightTimeseriesCache) idxKey(userID valueobject.UserID) string {
	return fmt.Sprintf("weights:timeseries:%s:idx", userID.Value())
}

type weightCacheRecord struct {
	ID                string    `json:"id"`
	UserID            string    `json:"user_id"`
	WeightKg          string    `json:"weight_kg"`
	BodyFatPercentage *string   `json:"body_fat_percentage,omitempty"`
	SkeletalMuscleKg  *string   `json:"skeletal_muscle_kg,omitempty"`
	MeasuredAt        time.Time `json:"measured_at"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

func (c *RedisWeightTimeseriesCache) FindByPeriod(ctx context.Context, userID valueobject.UserID, from, to time.Time) ([]*weightdomain.Weight, bool, error) {
	// ids (期間で絞り込んだweight_idの一覧) => idsをもとに、hashからid->jsonを一括で取得する (合計2RTT)
	ids, err := c.client.ZRangeArgs(ctx, redis.ZRangeArgs{
		Key:     c.idxKey(userID),
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
	// hashからid->jsonを一括で取得する
	vals, err := c.client.HMGet(ctx, c.dataKey(userID), ids...).Result()
	if err != nil {
		return nil, false, err
	}
	// jsonをentityに復元する
	weights := make([]*weightdomain.Weight, 0, len(ids))
	for _, v := range vals {
		s, ok := v.(string)
		if !ok {
			// HMGetの戻り値は[]interface{}なので、nilの時はcache missとして扱いcallerにDBを参照するようにする。
			return nil, false, nil
		}
		weight, err := decodeWeightCacheRecord(s)
		if err != nil {
			return nil, false, err
		}
		weights = append(weights, weight)

	}
	return weights, true, nil
}

func (c *RedisWeightTimeseriesCache) Save(ctx context.Context, weight *weightdomain.Weight) error {
	jsonStr, err := encodeWeightCacheRecord(weight)
	if err != nil {
		return err
	}
	idx := c.idxKey(weight.UserID())
	data := c.dataKey(weight.UserID())

	// MULTI/EXEC でZADDとHSETをatomicに実行する
	_, err = c.client.TxPipelined(ctx, func(p redis.Pipeliner) error {
		p.ZAdd(ctx, idx, redis.Z{Score: float64(weight.MeasuredAt().Unix()), Member: weight.ID().Value()})
		p.HSet(ctx, data, weight.ID().Value(), jsonStr)
		// TTLは毎回touchする、アクセス権があるときは失効しない
		p.Expire(ctx, idx, c.ttl)
		p.Expire(ctx, data, c.ttl)
		return nil
	})
	return err
}

func (c *RedisWeightTimeseriesCache) Delete(ctx context.Context, userID valueobject.UserID, weightID valueobject.WeightID) error {
	idx := c.idxKey(userID)
	data := c.dataKey(userID)
	_, err := c.client.TxPipelined(ctx, func(p redis.Pipeliner) error {
		p.ZRem(ctx, idx, weightID.Value())
		p.HDel(ctx, data, weightID.Value())
		return nil
	})
	return err
}

// ----helper functions----
func encodeWeightCacheRecord(weight *weightdomain.Weight) (string, error) {
	record := weightCacheRecord{
		ID:         weight.ID().Value(),
		UserID:     weight.UserID().Value(),
		WeightKg:   weight.WeightKg().String(),
		MeasuredAt: weight.MeasuredAt(),
		CreatedAt:  weight.CreatedAt(),
		UpdatedAt:  weight.UpdatedAt(),
	}
	if weight.BodyFatPercentage() != nil {
		s := weight.BodyFatPercentage().Value().String()
		record.BodyFatPercentage = &s
	}
	if weight.SkeletalMuscleKg() != nil {
		s := weight.SkeletalMuscleKg().Value().String()
		record.SkeletalMuscleKg = &s
	}
	b, err := json.Marshal(record)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func decodeWeightCacheRecord(data string) (*weightdomain.Weight, error) {
	var record weightCacheRecord
	err := json.Unmarshal([]byte(data), &record)
	if err != nil {
		return nil, err
	}
	weightID, err := valueobject.NewPrimaryIDFromString[valueobject.WeightID](record.ID)
	if err != nil {
		return nil, err
	}
	userID, err := valueobject.NewPrimaryIDFromString[valueobject.UserID](record.UserID)
	if err != nil {
		return nil, err
	}
	weightKg, err := valueobject.NewWeightKgFromString(record.WeightKg)
	if err != nil {
		return nil, err
	}
	var bodyFatPercentage *valueobject.Percentage
	if record.BodyFatPercentage != nil {
		bodyFatPercentage, err = valueobject.NewPercentageFromString(*record.BodyFatPercentage)
		if err != nil {
			return nil, err
		}
	}
	var skeletalMuscleKg *valueobject.WeightKg
	if record.SkeletalMuscleKg != nil {
		skeletalMuscleKg, err = valueobject.NewWeightKgFromString(*record.SkeletalMuscleKg)
		if err != nil {
			return nil, err
		}
	}
	return weightdomain.NewWeight(
		*weightID,
		*userID,
		*weightKg,
		bodyFatPercentage,
		skeletalMuscleKg,
		record.MeasuredAt,
		record.CreatedAt,
		record.UpdatedAt,
	), nil
}

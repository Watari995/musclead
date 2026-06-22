package healthsyncusecase

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	healthsyncdomain "github.com/Watari995/musclead/internal/healthsync/internal/domain"
	"github.com/Watari995/musclead/internal/myerror"
	"github.com/Watari995/musclead/internal/valueobject"
	weightpublicfunctions "github.com/Watari995/musclead/internal/weight/interface/publicfunctions"
)

type SyncWeights struct {
	tokenRepo       healthsyncdomain.TokenRepository
	metricsProvider healthsyncdomain.BodyMetricsProvider
	tokenRefresher  healthsyncdomain.TokenRefresher
	weightCommand   weightpublicfunctions.WeightCommand
	weightQuery     weightpublicfunctions.WeightQuery
}

func NewSyncWeights(
	tokenRepo healthsyncdomain.TokenRepository,
	metricsProvider healthsyncdomain.BodyMetricsProvider,
	tokenRefresher healthsyncdomain.TokenRefresher,
	weightCommand weightpublicfunctions.WeightCommand,
	weightQuery weightpublicfunctions.WeightQuery,
) *SyncWeights {
	return &SyncWeights{
		tokenRepo:       tokenRepo,
		metricsProvider: metricsProvider,
		tokenRefresher:  tokenRefresher,
		weightCommand:   weightCommand,
		weightQuery:     weightQuery,
	}
}

func (uc *SyncWeights) Execute(ctx context.Context) error {
	tokens, err := uc.tokenRepo.FindAllActive(ctx)
	if err != nil {
		return myerror.NewInternalError().Wrap(err)
	}

	for _, token := range tokens {
		if err := uc.syncToken(ctx, token); err != nil {
			slog.Error("healthsync: failed to sync user", "userID", token.UserID().Value(), "err", err)
		}
	}
	return nil
}

func (uc *SyncWeights) syncToken(ctx context.Context, token *healthsyncdomain.Token) error {
	if token.ExpiresAt().Before(time.Now()) {
		accessToken, refreshToken, expiresAt, err := uc.tokenRefresher.RefreshToken(ctx, token.RefreshToken())
		if err != nil {
			return fmt.Errorf("refresh token: %w", err)
		}
		token.UpdateTokens(accessToken, refreshToken, expiresAt)
		if err := uc.tokenRepo.Save(ctx, token); err != nil {
			return fmt.Errorf("save refreshed token: %w", err)
		}
	}

	var from time.Time
	if token.LastSyncedAt() != nil {
		from = *token.LastSyncedAt()
	} else {
		from = time.Now().AddDate(0, 0, -30)
	}

	metrics, err := uc.metricsProvider.FetchMetrics(ctx, token.AccessToken(), from, time.Now())
	if err != nil {
		return fmt.Errorf("fetch metrics: %w", err)
	}

	for _, m := range metrics {
		if err := uc.recordMetric(ctx, token, m); err != nil {
			slog.Error("healthsync: failed to record metric", "userID", token.UserID().Value(), "measuredAt", m.MeasuredAt, "err", err)
		}
	}

	token.SetLastSyncedAt(time.Now())
	if err := uc.tokenRepo.Save(ctx, token); err != nil {
		return fmt.Errorf("save last synced at: %w", err)
	}
	return nil
}

func (uc *SyncWeights) recordMetric(ctx context.Context, token *healthsyncdomain.Token, m healthsyncdomain.BodyMetrics) error {
	exists, err := uc.weightQuery.CheckIfExistsWeightByUserIDAndMeasuredAt(ctx, token.UserID(), m.MeasuredAt)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	weightKg, err := valueobject.NewWeightKgFromString(fmt.Sprintf("%g", m.Weight))
	if err != nil {
		return fmt.Errorf("parse weight: %w", err)
	}

	var bodyFatPercentage *valueobject.Percentage
	if m.BodyFatPercent != nil {
		bodyFatPercentage, err = valueobject.NewPercentageFromString(fmt.Sprintf("%g", *m.BodyFatPercent))
		if err != nil {
			return fmt.Errorf("parse body fat: %w", err)
		}
	}

	var skeletalMuscleKg *valueobject.WeightKg
	if m.SkeletalMuscleKg != nil {
		skeletalMuscleKg, err = valueobject.NewWeightKgFromString(fmt.Sprintf("%g", *m.SkeletalMuscleKg))
		if err != nil {
			return fmt.Errorf("parse skeletal muscle: %w", err)
		}
	}

	_, err = uc.weightCommand.Record(ctx, weightpublicfunctions.WeightRecordInput{
		UserID:            token.UserID(),
		WeightKg:          *weightKg,
		BodyFatPercentage: bodyFatPercentage,
		SkeletalMuscleKg:  skeletalMuscleKg,
		MeasuredAt:        m.MeasuredAt,
	})
	return err
}

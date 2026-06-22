package healthsync

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	healthsynchandler "github.com/Watari995/musclead/internal/healthsync/internal/handler"
	healthsyncusecase "github.com/Watari995/musclead/internal/healthsync/internal/usecase"
	"github.com/Watari995/musclead/internal/weight/interface/publicfunctions"
	"github.com/go-gorp/gorp/v3"

	healthsyncinfra "github.com/Watari995/musclead/internal/healthsync/internal/infra"
)

const syncInterval = 90 * time.Second

type Module struct {
	syncWeights *healthsyncusecase.SyncWeights
	Handler     http.Handler
}

func NewModule(
	dbmap *gorp.DbMap,
	httpClint *http.Client,
	clientID string,
	clientSecret string,
	weightCommand publicfunctions.WeightCommand,
	weightQuery publicfunctions.WeightQuery,
) *Module {
	// repo
	dbmap.AddTableWithName(healthsyncinfra.HealthPlanetTokenModel{}, "healthplanet_tokens").SetKeys(false, "id")
	tokenRepo := healthsyncinfra.NewTokenRepository(dbmap)
	client := healthsyncinfra.NewHealthPlanetClient(httpClint, clientID, clientSecret)

	// use-case
	connect := healthsyncusecase.NewConnectHealthPlanet(tokenRepo, client)
	sync := healthsyncusecase.NewSyncWeights(tokenRepo, client, client, weightCommand, weightQuery)

	healthsyncHandler := healthsynchandler.New(clientID, connect)

	return &Module{
		syncWeights: sync,
		Handler:     healthsyncHandler,
	}
}

// 90秒ごとに体重を同期する。
func (m *Module) RunSync(ctx context.Context) {
	ticker := time.NewTicker(syncInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := m.syncWeights.Execute(ctx); err != nil {
				slog.Error("healthsync: failed to sync weights", "err", err)
			}
		}
	}
}

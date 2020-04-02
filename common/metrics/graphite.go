package metrics

import (
	"github.com/marpaia/graphite-golang"
	"madnessBot/common/logger"
	"madnessBot/config"
)

type GraphiteMetric struct {
	server *graphite.Graphite
}

var graphiteInstance = &GraphiteMetric{server: nil}

func Graphite() *GraphiteMetric {
	return graphiteInstance
}

func Init() {
	graphiteInstance = &GraphiteMetric{server: nil}

	graphiteSrv, err := graphite.NewGraphite(
		config.Config.Graphite.Host,
		config.Config.Graphite.Port,
	)

	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to initialize graphite")
		graphiteInstance = &GraphiteMetric{server: nil}
	}

	graphiteInstance = &GraphiteMetric{server: graphiteSrv}
}

func (g *GraphiteMetric) Send(metric graphite.Metric) {
	if g.server == nil {
		logger.Log.Debug().Interface("metric", metric).Msg("Graphite is disabled")
		return
	}

	if err := g.server.SendMetric(metric); err != nil {
		logger.Log.Error().Err(err).Msg("Failed to send metric")
	}
}

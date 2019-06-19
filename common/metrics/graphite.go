package metrics

import (
	"github.com/Defman21/madnessBot/common"
	"github.com/marpaia/graphite-golang"
	"os"
	"strconv"
)

type GraphiteMetric struct {
	server *graphite.Graphite
}

var graphiteInstance = &GraphiteMetric{server: nil}

func Graphite() *GraphiteMetric {
	return graphiteInstance
}

func Init() {
	port, err := strconv.Atoi(os.Getenv("GRAPHITE_PORT"))

	if err != nil {
		common.Log.Error().Err(err).Msg("Invalid GRAPHITE_PORT")
		graphiteInstance = &GraphiteMetric{server: nil}
	}

	graphiteSrv, err := graphite.NewGraphite(os.Getenv("GRAPHITE_HOST"), port)

	if err != nil {
		common.Log.Error().Err(err).Msg("Failed to initialize graphite")
		graphiteInstance = &GraphiteMetric{server: nil}
	}

	graphiteInstance = &GraphiteMetric{server: graphiteSrv}
}

func (g *GraphiteMetric) Send(metric graphite.Metric) {
	if g.server == nil {
		common.Log.Debug().Interface("metric", metric).Msg("Graphite is disabled")
		return
	}

	if err := g.server.SendMetric(metric); err != nil {
		common.Log.Error().Err(err).Msg("Failed to send metric")
	}
}

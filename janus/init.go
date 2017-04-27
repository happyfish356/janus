package main

import (
	"strings"
	"time"

	"config"
	log "github.com/Sirupsen/logrus"
	"github.com/bshuster-repo/logrus-logstash-hook"
	stats "github.com/hellofresh/stats-go"
	opentracing "github.com/opentracing/opentracing-go"
	tracerfactory "opentracing"
	"store"
)

var (
	globalConfig *config.Specification
	statsClient  stats.Client
	storage      store.Store
)

func init() {
	c, err := config.Load(configFile)
	if nil != err {
		log.WithError(err).Panic("Could not parse the environment configurations")
	}

	globalConfig = c
}

// initializes the basic configuration for the log wrapper
func init() {
	level, err := log.ParseLevel(strings.ToLower(globalConfig.LogLevel))
	if err != nil {
		log.WithError(err).Error("Error getting log level")
	}

	log.SetLevel(level)
	log.SetFormatter(&logrus_logstash.LogstashFormatter{
		Type:            "Janus",
		TimestampFormat: time.RFC3339Nano,
	})
}

// initializes distributed tracing
func init() {
	log.Debug("initializing Open Tracing")
	tracer, err := tracerfactory.Build(globalConfig.Tracing)
	if err != nil {
		log.WithError(err).Panic("Could not build a tracer for open tracing")
	}

	opentracing.InitGlobalTracer(tracer)
}

func init() {
	sectionsTestsMap, err := stats.ParseSectionsTestsMap(globalConfig.Stats.IDs)
	if err != nil {
		log.WithError(err).WithField("config", globalConfig.Stats.IDs).
			Error("Failed to parse stats second level IDs from env")
		sectionsTestsMap = map[stats.PathSection]stats.SectionTestDefinition{}
	}
	log.WithField("config", globalConfig.Stats.IDs).
		WithField("map", sectionsTestsMap.String()).
		Debug("Setting stats second level IDs")

	statsClient, err = stats.NewClient(globalConfig.Stats.DSN, globalConfig.Stats.Prefix)
	if err != nil {
		log.WithError(err).Panic("Error initializing statsd client")
	}

	statsClient.SetHTTPMetricCallback(stats.NewHasIDAtSecondLevelCallback(sectionsTestsMap))
}

// initializes the storage and managers
func init() {
	s, err := store.Build(globalConfig.Storage.DSN)
	if nil != err {
		log.Panic(err)
	}

	storage = s
}

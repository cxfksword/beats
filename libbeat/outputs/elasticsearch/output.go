package elasticsearch

import (
	"bytes"
	"crypto/tls"
	"errors"
	"io/ioutil"
	"net/url"
	"strings"
	"time"

	"github.com/cxfksword/beats/libbeat/common"
	"github.com/cxfksword/beats/libbeat/logp"
	"github.com/cxfksword/beats/libbeat/outputs"
	"github.com/cxfksword/beats/libbeat/outputs/mode"
)

type elasticsearchOutput struct {
	index string
	mode  mode.ConnectionMode
	topology
}

func init() {
	outputs.RegisterOutputPlugin("elasticsearch", New)
}

var (
	debug = logp.MakeDebug("elasticsearch")
)

var (
	// ErrNotConnected indicates failure due to client having no valid connection
	ErrNotConnected = errors.New("not connected")

	// ErrJSONEncodeFailed indicates encoding failures
	ErrJSONEncodeFailed = errors.New("json encode failed")

	// ErrResponseRead indicates error parsing Elasticsearch response
	ErrResponseRead = errors.New("bulk item status parse failed.")
)

// NewOutput instantiates a new output plugin instance publishing to elasticsearch.
func New(cfg *common.Config, topologyExpire int) (outputs.Outputer, error) {
	if !cfg.HasField("bulk_max_size") {
		cfg.SetInt("bulk_max_size", 0, defaultBulkSize)
	}

	output := &elasticsearchOutput{}
	err := output.init(cfg, topologyExpire)
	if err != nil {
		return nil, err
	}
	return output, nil
}

func (out *elasticsearchOutput) init(
	cfg *common.Config,
	topologyExpire int,
) error {
	config := defaultConfig
	if err := cfg.Unpack(&config); err != nil {
		return err
	}

	tlsConfig, err := outputs.LoadTLSConfig(config.TLS)
	if err != nil {
		return err
	}

	clients, err := mode.MakeClients(cfg, makeClientFactory(tlsConfig, &config))
	if err != nil {
		return err
	}

	maxRetries := config.MaxRetries
	maxAttempts := maxRetries + 1 // maximum number of send attempts (-1 = infinite)
	if maxRetries < 0 {
		maxAttempts = 0
	}

	var waitRetry = time.Duration(1) * time.Second
	var maxWaitRetry = time.Duration(60) * time.Second

	out.clients = clients
	loadBalance := config.LoadBalance
	m, err := mode.NewConnectionMode(clients, !loadBalance,
		maxAttempts, waitRetry, config.Timeout, maxWaitRetry)
	if err != nil {
		return err
	}

	loadTemplate(config.Template, clients)

	if config.SaveTopology {
		err := out.EnableTTL()
		if err != nil {
			logp.Err("Fail to set _ttl mapping: %s", err)
			// keep trying in the background
			go func() {
				for {
					err := out.EnableTTL()
					if err == nil {
						break
					}
					logp.Err("Fail to set _ttl mapping: %s", err)
					time.Sleep(5 * time.Second)
				}
			}()
		}
	}

	out.TopologyExpire = 15000
	if topologyExpire != 0 {
		out.TopologyExpire = topologyExpire * 1000 // millisec
	}

	out.mode = m
	out.index = config.Index

	return nil
}

// loadTemplate checks if the index mapping template should be loaded
// In case template loading is enabled, template is written to index
func loadTemplate(config Template, clients []mode.ProtocolClient) {
	// Check if template should be loaded
	// Not being able to load the template will output an error but will not stop execution
	if config.Name != "" && len(clients) > 0 {

		// Always takes the first client
		esClient := clients[0].(*Client)

		logp.Info("Loading template enabled. Trying to load template: %v", config.Path)

		exists := esClient.CheckTemplate(config.Name)

		// Check if template already exist or should be overwritten
		if !exists || config.Overwrite {

			if config.Overwrite {
				logp.Info("Existing template will be overwritten, as overwrite is enabled.")
			}

			// Load template from file
			content, err := ioutil.ReadFile(config.Path)
			if err != nil {
				logp.Err("Could not load template from file path: %s; Error: %s", config.Path, err)
			} else {
				reader := bytes.NewReader(content)
				err = esClient.LoadTemplate(config.Name, reader)

				if err != nil {
					logp.Err("Could not load template: %v", err)
				}
			}
		} else {
			logp.Info("Template already exists and will not be overwritten.")
		}

	}
}

func makeClientFactory(
	tls *tls.Config,
	config *elasticsearchConfig,
) func(string) (mode.ProtocolClient, error) {
	return func(host string) (mode.ProtocolClient, error) {
		esURL, err := getURL(config.Protocol, config.Path, host)
		if err != nil {
			logp.Err("Invalid host param set: %s, Error: %v", host, err)
			return nil, err
		}

		var proxyURL *url.URL
		if config.ProxyURL != "" {
			proxyURL, err = url.Parse(config.ProxyURL)
			if err != nil || !strings.HasPrefix(proxyURL.Scheme, "http") {
				// Proxy was bogus. Try prepending "http://" to it and
				// see if that parses correctly. If not, we fall
				// through and complain about the original one.
				proxyURL, err = url.Parse("http://" + config.ProxyURL)
				if err != nil {
					return nil, err
				}
			}

			logp.Info("Using proxy URL: %s", proxyURL)
		}

		params := config.Params
		if len(params) == 0 {
			params = nil
		}
		client := NewClient(
			esURL, config.Index, proxyURL, tls,
			config.Username, config.Password,
			params)
		return client, nil
	}
}

func (out *elasticsearchOutput) Close() error {
	return out.mode.Close()
}

func (out *elasticsearchOutput) PublishEvent(
	signaler outputs.Signaler,
	opts outputs.Options,
	event common.MapStr,
) error {
	return out.mode.PublishEvent(signaler, opts, event)
}

func (out *elasticsearchOutput) BulkPublish(
	trans outputs.Signaler,
	opts outputs.Options,
	events []common.MapStr,
) error {
	return out.mode.PublishEvents(trans, opts, events)
}

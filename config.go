// Copyright 2016-2017 冯立强 fenglq@tingyun.com.  All rights reserved.

package tingyun

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"time"

	"github.com/TingYunAPM/go/utils/cache_config"
	"github.com/TingYunAPM/go/utils/logger"
	"github.com/TingYunAPM/go/utils/service"
)

const (
	configLocalStringNbsHost            = 1
	configLocalStringNbsLicenseKey      = 2
	configLocalStringNbsAppName         = 3
	configLocalStringNbsLevel           = log.ConfigStringNBSLevel
	configLocalStringNbsLogFileName     = log.ConfigStringNBSLogFileName
	configLocalStringMax                = 8
	configLocalBoolAgentEnable          = 1
	configLocalBoolSSL                  = 2
	configLocalBoolAudit                = log.ConfigBoolNBSAudit
	configLocalBoolMax                  = 8
	configLocalIntegerNbsPort           = 1
	configLocalIntegerNbsSaveCount      = 2
	configLocalIntegerNbsMaxLogSize     = log.ConfigIntegerNBSMaxLogSize
	configLocalIntegerNbsMaxLogCount    = log.ConfigIntegerNBSMaxLogCount
	configLocalIntegerNbsActionCacheMax = 5
	configLocalIntegerMax               = 8

	configServerStringAppSessionKey     = 1
	configServerStringTingyunIdSecret   = 2
	configServerStringApplicationId     = 3
	configServerStringMax               = 8
	configServerBoolEnabled             = 1
	configServerBoolMax                 = 8
	configServerIntegerDataSentInterval = 1
	configServerIntegerApdex_t          = 2
	configServerIntegerMax              = 8

	configServerConfigStringActionTracerRecordSQL      = 1
	configServerConfigStringRumScript                  = 2
	configServerConfigStringExternalUrlParamsCaptured  = 3
	configServerConfigStringWebActionURIParamsCaptured = 4
	configServerConfigStringInstrumentationCustom      = 5
	configServerConfigStringQuantile                   = 6
	configServerConfigStringMax                        = 8

	configServerConfigBoolAgentEnabled                   = 1
	configServerConfigBoolAutoActionNaming               = 2
	configServerConfigBoolCaptureParams                  = 3
	configServerConfigBoolErrorCollectorEnabled          = 4
	configServerConfigBoolErrorCollectorRecordDBErrors   = 5
	configServerConfigBoolActionTracerEnabled            = 6
	configServerConfigBoolActionTracerSlowSQL            = 7
	configServerConfigBoolActionTracerExplainEnabled     = 8
	configServerConfigBoolTransactionTracerEnabled       = 9
	configServerConfigBoolActionTracerNbsua              = 10
	configServerConfigBoolRumEnabled                     = 11
	configServerConfigBoolIgnoreStaticResources          = 12
	configServerConfigBoolActionTracerRemoveTrailingPath = 13
	configServerConfigBoolHotspotEnabled                 = 14
	configServerConfigBoolRumMixEnabled                  = 15
	configServerConfigBoolTransactionTracerThrift        = 16
	configServerConfigBoolMQEnabled                      = 17
	configServerConfigBoolResourceEnabled                = 18
	configServerConfigBoolLogTracking                    = 19
	configServerConfigBoolMax                            = 24

	configServerConfigIntegerActionTracerActionThreshold     = 1
	configServerConfigIntegerActionTracerSlowSQLThreshold    = 2
	configServerConfigIntegerActionTracerExplainThreshold    = 3
	configServerConfigIntegerActionTracerStacktraceThreshold = 4
	configServerConfigIntegerRumSampleRatio                  = 5
	configServerConfigIntegerResourceLow                     = 6
	configServerConfigIntegerResourceHigh                    = 7
	configServerConfigIntegerResourceSafe                    = 8
	configServerConfigIntegerMax                             = 10
)

var localStringKeyMap = map[string]int{
	"nbs.host":          configLocalStringNbsHost,
	"nbs.license_key":   configLocalStringNbsLicenseKey,
	"nbs.app_name":      configLocalStringNbsAppName,
	"nbs.level":         configLocalStringNbsLevel,
	"nbs.log_file_name": configLocalStringNbsLogFileName,
}
var localBoolKeyMap = map[string]int{
	"nbs.agent_enabled": configLocalBoolAgentEnable,
	"nbs.ssl":           configLocalBoolSSL,
	"nbs.audit":         configLocalBoolAudit,
}

var localIntegerKeyMap = map[string]int{
	"nbs.port":             configLocalIntegerNbsPort,
	"nbs.savecount":        configLocalIntegerNbsSaveCount,
	"nbs.max_log_size":     configLocalIntegerNbsMaxLogSize,
	"nbs.max_log_count":    configLocalIntegerNbsMaxLogCount,
	"nbs.action_cache_max": configLocalIntegerNbsActionCacheMax,
}

var serverStringKeyMap = map[string]int{
	"appSessionKey":   configServerStringAppSessionKey,
	"tingyunIdSecret": configServerStringTingyunIdSecret,
	"applicationId":   configServerStringApplicationId,
}

var serverBoolKeyMap = map[string]int{
	"enabled": configServerBoolEnabled,
}

var serverIntegerKeyMap = map[string]int{
	"dataSentInterval": configServerIntegerDataSentInterval,
	"apdex_t":          configServerIntegerApdex_t,
}

var serverConfigStringKeyMap = map[string]int{
	"nbs.action_tracer.record_sql":       configServerConfigStringActionTracerRecordSQL,
	"nbs.rum.script":                     configServerConfigStringRumScript,
	"nbs.external_url_params_captured":   configServerConfigStringExternalUrlParamsCaptured,
	"nbs.web_action_uri_params_captured": configServerConfigStringWebActionURIParamsCaptured,
	"nbs.instrumentation_custom":         configServerConfigStringInstrumentationCustom,
	"nbs.quantile":                       configServerConfigStringQuantile,
}

var serverConfigBoolKeyMap = map[string]int{
	"nbs.agent_enabled":                      configServerConfigBoolAgentEnabled,
	"nbs.auto_action_naming":                 configServerConfigBoolAutoActionNaming,
	"nbs.capture_params":                     configServerConfigBoolCaptureParams,
	"nbs.error_collector.enabled":            configServerConfigBoolErrorCollectorEnabled,
	"nbs.error_collector.record_db_errors":   configServerConfigBoolErrorCollectorRecordDBErrors,
	"nbs.action_tracer.enabled":              configServerConfigBoolActionTracerEnabled,
	"nbs.action_tracer.slow_sql":             configServerConfigBoolActionTracerSlowSQL,
	"nbs.action_tracer.explain_enabled":      configServerConfigBoolActionTracerExplainEnabled,
	"nbs.transaction_tracer.enabled":         configServerConfigBoolTransactionTracerEnabled,
	"nbs.action_tracer.nbsua":                configServerConfigBoolActionTracerNbsua,
	"nbs.rum.enabled":                        configServerConfigBoolRumEnabled,
	"nbs.ignore_static_resources":            configServerConfigBoolIgnoreStaticResources,
	"nbs.action_tracer.remove_trailing_path": configServerConfigBoolActionTracerRemoveTrailingPath,
	"nbs.hotspot.enabled":                    configServerConfigBoolHotspotEnabled,
	"nbs.rum.mix_enabled":                    configServerConfigBoolRumMixEnabled,
	"nbs.transaction_tracer.thrift":          configServerConfigBoolTransactionTracerThrift,
	"nbs.mq.enabled":                         configServerConfigBoolMQEnabled,
	"nbs.resource.enabled":                   configServerConfigBoolResourceEnabled,
	"nbs.log_tracking":                       configServerConfigBoolLogTracking,
}

var serverConfigIntegerKeyMap = map[string]int{
	"nbs.action_tracer.action_threshold":      configServerConfigIntegerActionTracerActionThreshold,
	"nbs.action_tracer.slow_sql_threshold":    configServerConfigIntegerActionTracerSlowSQLThreshold,
	"nbs.action_tracer.explain_threshold":     configServerConfigIntegerActionTracerExplainThreshold,
	"nbs.action_tracer.stack_trace_threshold": configServerConfigIntegerActionTracerStacktraceThreshold,
	"nbs.rum.sample_ratio":                    configServerConfigIntegerRumSampleRatio,
	"nbs.resource.low":                        configServerConfigIntegerResourceLow,
	"nbs.resource.high":                       configServerConfigIntegerResourceHigh,
	"nbs.resource.safe":                       configServerConfigIntegerResourceSafe,
}

type configKeyMaps struct {
	strings  map[string]int
	bools    map[string]int
	integers map[string]int
}

var local_key_maps = configKeyMaps{localStringKeyMap, localBoolKeyMap, localIntegerKeyMap}
var server_key_maps = configKeyMaps{serverStringKeyMap, serverBoolKeyMap, serverIntegerKeyMap}

type configurations struct {
	local       cache_config.Configuration
	server      cache_config.Configuration
	server_ext  cache_config.Configuration
	svc         service.Service
	apdexs      apdex_action_map
	started     bool
	login_error bool
	login_count int64
	reported    bool
}

func parse_config(filename string, c *cache_config.Configuration) error {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	jsonData := map[string]interface{}{}
	if err = json.Unmarshal(bytes, &jsonData); err != nil {
		return err
	}
	for k, v := range jsonData {
		c.Update(localStringKeyMap, localBoolKeyMap, localIntegerKeyMap, k, v)
	}
	c.Commit()
	return nil
}
func (c *configurations) Init(configfile string) error {
	c.local.Init(configLocalStringMax, configLocalBoolMax, configLocalIntegerMax)
	c.server.Init(configServerStringMax, configServerBoolMax, configServerIntegerMax)
	c.server_ext.Init(configServerConfigStringMax, configServerConfigBoolMax, configServerConfigIntegerMax)
	c.apdexs.Init()
	err := parse_config(configfile, &c.local)
	c.started = err == nil
	c.reported = false
	c.login_count = 0
	c.login_error = false
	if c.started {
		c.svc.Start(func(running func() bool) {
			lastTime := time.Now()
			lastModify, err := modifyTime(configfile)
			if err != nil {
				lastModify = lastTime
			}
			for running() {
				time.Sleep(100 * time.Millisecond)
				if now := time.Now(); now.Sub(lastTime) < 60*time.Second {
					continue
				}
				if modTime, err := modifyTime(configfile); err == nil {
					if modTime.Equal(lastModify) {
						continue
					}
					if parse_config(configfile, &c.local) == nil {
						lastModify = modTime
					}
				}
			}

		})
	}
	return err
}
func (c *configurations) NeverLogin() bool { return c.login_count == 0 }
func (c *configurations) HasLogin() bool   { return c.login_count > 0 && !c.login_error }
func (c *configurations) Release() {
	if c.started {
		c.started = false
		c.svc.Stop()
	}
}
func (c *configurations) UpdateServerConfig(result map[string]interface{}) bool {
	for {
		if err, _ := jsonReadString(result, "appSessionKey"); err != nil {
			break
		} else if err, _ = jsonToString(result, "applicationId"); err != nil {
			break
		} else if err, _ = jsonReadBool(result, "enabled"); err != nil {
			break
		} else if err, _ = jsonReadInt(result, "dataSentInterval"); err != nil {
			break
		} else if err, apdex_t := jsonReadInt(result, "apdex_t"); err != nil {
			break
		} else {
			for k, v := range result {
				c.server.Update(serverStringKeyMap, serverBoolKeyMap, serverIntegerKeyMap, k, v)
			}

			if err, config := jsonReadObjects(result, "config"); err == nil {
				for k, v := range config {
					c.server_ext.Update(serverConfigStringKeyMap, serverConfigBoolKeyMap, serverConfigIntegerKeyMap, k, v)
				}
			}
			if err, config := jsonReadObjects(result, "agreementConfig"); err == nil {
				for k, v := range config {
					c.server_ext.Update(serverConfigStringKeyMap, serverConfigBoolKeyMap, serverConfigIntegerKeyMap, k, v)
				}
			}
			c.apdexs.apdex_t = apdex_t
			if err, actionApdex := jsonReadObjects(result, "actionApdex"); err == nil {
				for k, v := range actionApdex {
					if err, val := readInt(v); err == nil {
						c.apdexs.Update(k, val)
					}
				}
			}

			c.server.Commit()
			c.server_ext.Commit()
			c.apdexs.Commit()
			c.login_count += 1
		}
		return true
	}
	return false
}
func (c *configurations) UpdateConfig(onFirst func()) {
	if c.login_count > 0 {
		if !c.reported {
			c.reported = true
			onFirst()
		}
	}
}
func config_value(config *cache_config.Configuration, key string, maps *configKeyMaps) (interface{}, bool) {
	if v, found := maps.strings[key]; found {
		return config.CStrings.Find(v)
	} else if v, found := maps.bools[key]; found {
		return config.CBools.Find(v)
	} else if v, found := maps.integers[key]; found {
		return config.CIntegers.Find(v)
	}
	return nil, false
}
func (c *configurations) Value(name string) (interface{}, bool) {
	if v, found := config_value(&c.server, name, &server_key_maps); found {
		return v, found
	}
	return config_value(&c.local, name, &local_key_maps)
}

func modifyTime(filename string) (time.Time, error) {
	info, err := os.Stat(filename)
	if err != nil {
		return time.Time{}, err
	}
	return info.ModTime(), nil
}

//func configInt(c *configBase.ConfigBase, key string, defaultValue int) int {
//	if t, found := c.Value(key); found {
//		switch r := t.(type) {
//		case float64:
//			return int(r)
//		case float32:
//			return int(r)
//		case int32:
//			return int(r)
//		case int64:
//			return int(r)
//		case uint32:
//			return int(r)
//		case uint64:
//			return int(r)
//		}
//	}
//	return defaultValue
//}

type apdex_action_map struct {
	current int
	apdex_t int
	arrays  [4]map[string]int
}

func (s *apdex_action_map) Init() *apdex_action_map {
	s.current = 3
	s.apdex_t = 500
	for i := 0; i < 4; i++ {
		s.arrays[i] = make(map[string]int)
	}
	return s
}
func (s *apdex_action_map) Read(key string) int {
	if r, ok := s.arrays[s.current][key]; ok {
		return int(r)
	}
	return s.apdex_t
}
func (s *apdex_action_map) Update(key string, value int) {
	s.arrays[(s.current+1)%4][key] = value
}
func (s *apdex_action_map) Commit() {
	s.current = (s.current + 1) % 4
}

package exporter

import "encoding/json"

type SlackExporterPayloadSettings struct {
	ExporterPayloadSettings
	Channel string `json:"channel"`
}

type JiraExporterPayloadSettings struct {
	ExporterPayloadSettings
	Board string `json:"board"`
}

type ExporterPayloadSettings struct {
}

type ExporterPayload struct {
	ID        string          `json:"id"`
	ProductID string          `json:"product_id"`
	Type      string          `json:"type"`
	Settings  json.RawMessage `json:"settings"`
}

func NewExporterPayload(exporter Exporter) *ExporterPayload {
	var err error

	var settings json.RawMessage
	switch exporter.Type {
	case ExporterTypeSlack:
		_settings := exporter.Settings.(SlackExporterSettings) // nolint: errcheck
		settings, err = json.Marshal(SlackExporterPayloadSettings{
			Channel: _settings.Channel,
		})
		if err != nil {
			panic(err)
		}

	case ExporterTypeJira:
		_settings := exporter.Settings.(JiraExporterSettings) // nolint: errcheck
		settings, err = json.Marshal(JiraExporterPayloadSettings{
			Board: _settings.Board,
		})
		if err != nil {
			panic(err)
		}

	default:
		panic(exporter.Type)
	}

	return &ExporterPayload{
		ID:        exporter.ID,
		ProductID: exporter.ProductID,
		Type:      exporter.Type,
		Settings:  settings,
	}
}

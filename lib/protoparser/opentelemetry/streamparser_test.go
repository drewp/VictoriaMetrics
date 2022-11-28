package opentelemetry

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/VictoriaMetrics/VictoriaMetrics/lib/protoparser/opentelemetry/opentelemetrypb"
	"github.com/golang/protobuf/jsonpb"

	"github.com/VictoriaMetrics/VictoriaMetrics/lib/prompb"
)

func TestParseStream(t *testing.T) {
	prettifyLabel := func(label prompb.Label) string {
		return fmt.Sprintf("name=%q value=%q", label.Name, label.Value)
	}
	prettifySample := func(sample prompb.Sample) string {
		return fmt.Sprintf("sample=%f timestamp: %d", sample.Value, sample.Timestamp)
	}
	f := func(name string, samples []*opentelemetrypb.Metric, expectedTss []prompb.TimeSeries) {
		t.Run(name, func(t *testing.T) {
			pbRequest := opentelemetrypb.ExportMetricsServiceRequest{
				ResourceMetrics: []*opentelemetrypb.ResourceMetrics{generateOTLPSamples(samples...)},
			}
			m := jsonpb.Marshaler{
				EnumsAsInts: true,
			}
			data, err := m.MarshalToString(&pbRequest)
			if err != nil {
				t.Fatalf("cannot marshal data: %s", err)
			}
			err = ParseStream(bytes.NewBuffer([]byte(data)), true, false, func(tss []prompb.TimeSeries) error {
				if len(tss) != len(expectedTss) {
					t.Fatalf("not expected tss count, got: %d, want: %d", len(tss), len(expectedTss))
				}
				for i := 0; i < len(tss); i++ {
					ts := tss[i]
					tsWant := expectedTss[i]
					if len(ts.Labels) != len(tsWant.Labels) {
						t.Fatalf("idx: %d, not expected labels count, got: %d, want: %d", i, len(ts.Labels), len(tsWant.Labels))
					}
					for j, label := range ts.Labels {
						wantLabel := tsWant.Labels[j]
						if !reflect.DeepEqual(label, wantLabel) {
							t.Fatalf("idx: %d, label idx: %d, not equal label pairs, \ngot: \n%s, \nwant: \n%s", i, j, prettifyLabel(label), prettifyLabel(wantLabel))
						}
					}
					if len(ts.Samples) != len(tsWant.Samples) {
						t.Fatalf("idx: %d, not expected samples count, got: %d, want: %d", i, len(ts.Samples), len(tsWant.Samples))
					}
					for j, sample := range ts.Samples {
						wantSample := tsWant.Samples[j]
						if !reflect.DeepEqual(sample, wantSample) {
							t.Fatalf("idx: %d, label idx: %d, not equal sample pairs, \ngot: \n%s,\nwant: \n%s", i, j, prettifySample(sample), prettifySample(wantSample))
						}
					}
				}
				return nil
			})
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
		})
	}
	jobLabelValue := prompb.Label{Name: []byte(`job`), Value: []byte(`vm`)}
	leLabel := func(value string) prompb.Label {
		return prompb.Label{Name: boundLabel, Value: []byte(value)}
	}
	kvLabel := func(k, v string) prompb.Label {
		return prompb.Label{Name: []byte(k), Value: []byte(v)}
	}
	f("test all metric types",
		[]*opentelemetrypb.Metric{generateGauge("my-gauge"), generateHistogram("my-histogram"), generateSum("my-sum"), generateSummary("my-summary")},
		[]prompb.TimeSeries{
			newPromPBTs("my-gauge", 15000, 15.0, false, jobLabelValue, kvLabel("label1", "value1")),
			newPromPBTs("my-histogram_count", 30000, 15.0, false, jobLabelValue, kvLabel("label2", "value2")),
			newPromPBTs("my-histogram_sum", 30000, 30.0, false, jobLabelValue, kvLabel("label2", "value2")),
			newPromPBTs("my-histogram_bucket", 30000, 0.0, false, jobLabelValue, kvLabel("label2", "value2"), leLabel("0.1")),
			newPromPBTs("my-histogram_bucket", 30000, 5.0, false, jobLabelValue, kvLabel("label2", "value2"), leLabel("0.5")),
			newPromPBTs("my-histogram_bucket", 30000, 15.0, false, jobLabelValue, kvLabel("label2", "value2"), leLabel("1")),
			newPromPBTs("my-histogram_bucket", 30000, 15.0, false, jobLabelValue, kvLabel("label2", "value2"), leLabel("5")),
			newPromPBTs("my-histogram_bucket", 30000, 15.0, false, jobLabelValue, kvLabel("label2", "value2"), leLabel("+Inf")),
			newPromPBTs("my-sum", 150000, 15.5, false, jobLabelValue, kvLabel("label5", "value5")),
			newPromPBTs("my-summary_sum", 35000, 32.5, false, jobLabelValue, kvLabel("label6", "value6")),
			newPromPBTs("my-summary_count", 35000, 5.0, false, jobLabelValue, kvLabel("label6", "value6")),
			newPromPBTs("my-summary", 35000, 7.5, false, jobLabelValue, kvLabel("label6", "value6"), kvLabel("quantile", "0.1")),
			newPromPBTs("my-summary", 35000, 10.0, false, jobLabelValue, kvLabel("label6", "value6"), kvLabel("quantile", "0.5")),
			newPromPBTs("my-summary", 35000, 15.0, false, jobLabelValue, kvLabel("label6", "value6"), kvLabel("quantile", "1")),
		})
	f("test gauge",
		[]*opentelemetrypb.Metric{generateGauge("my-gauge")},
		[]prompb.TimeSeries{newPromPBTs("my-gauge", 15000, 15.0, false, jobLabelValue, kvLabel("label1", "value1"))})
}

func attributesFromKV(kvs ...[2]string) []opentelemetrypb.KeyValue {
	var r []opentelemetrypb.KeyValue
	for _, kv := range kvs {
		r = append(r, opentelemetrypb.KeyValue{
			Key:   kv[0],
			Value: opentelemetrypb.AnyValue{Value: &opentelemetrypb.AnyValue_StringValue{StringValue: kv[1]}},
		})
	}
	return r
}

func generateGauge(name string, points ...*opentelemetrypb.NumberDataPoint) *opentelemetrypb.Metric {
	defaultPoints := []*opentelemetrypb.NumberDataPoint{
		{

			Attributes:   attributesFromKV([2]string{"label1", "value1"}),
			Value:        &opentelemetrypb.NumberDataPoint_AsInt{AsInt: 15},
			TimeUnixNano: uint64(15 * time.Second),
		},
	}
	if len(points) == 0 {
		points = defaultPoints
	}
	return &opentelemetrypb.Metric{
		Name: name,
		Data: &opentelemetrypb.Metric_Gauge{
			Gauge: &opentelemetrypb.Gauge{
				DataPoints: points,
			},
		},
	}
}

func generateHistogram(name string, points ...*opentelemetrypb.HistogramDataPoint) *opentelemetrypb.Metric {
	defaultPoints := []*opentelemetrypb.HistogramDataPoint{
		{

			Attributes:     attributesFromKV([2]string{"label2", "value2"}),
			Count:          15,
			Sum_:           &opentelemetrypb.HistogramDataPoint_Sum{Sum: 30.0},
			ExplicitBounds: []float64{0.1, 0.5, 1.0, 5.0},
			BucketCounts:   []uint64{0, 5, 10, 0, 0},
			TimeUnixNano:   uint64(30 * time.Second),
		},
	}
	if len(points) == 0 {
		points = defaultPoints
	}
	return &opentelemetrypb.Metric{
		Name: name,
		Data: &opentelemetrypb.Metric_Histogram{
			Histogram: &opentelemetrypb.Histogram{
				AggregationTemporality: opentelemetrypb.AggregationTemporality_AGGREGATION_TEMPORALITY_CUMULATIVE,
				DataPoints:             points,
			},
		},
	}
}

func generateSum(name string, points ...*opentelemetrypb.NumberDataPoint) *opentelemetrypb.Metric {
	defaultPoints := []*opentelemetrypb.NumberDataPoint{
		{
			Attributes:   attributesFromKV([2]string{"label5", "value5"}),
			Value:        &opentelemetrypb.NumberDataPoint_AsDouble{AsDouble: 15.5},
			TimeUnixNano: uint64(150 * time.Second),
		},
	}
	if len(points) == 0 {
		points = defaultPoints
	}
	return &opentelemetrypb.Metric{
		Name: name,
		Data: &opentelemetrypb.Metric_Sum{
			Sum: &opentelemetrypb.Sum{
				AggregationTemporality: opentelemetrypb.AggregationTemporality_AGGREGATION_TEMPORALITY_CUMULATIVE,
				DataPoints:             points,
			},
		},
	}
}

func generateSummary(name string, points ...*opentelemetrypb.SummaryDataPoint) *opentelemetrypb.Metric {
	defaultPoints := []*opentelemetrypb.SummaryDataPoint{
		{
			Attributes:   attributesFromKV([2]string{"label6", "value6"}),
			TimeUnixNano: uint64(35 * time.Second),
			Sum:          32.5,
			Count:        5,
			QuantileValues: []*opentelemetrypb.SummaryDataPoint_ValueAtQuantile{
				{
					Quantile: 0.1,
					Value:    7.5,
				},
				{
					Quantile: 0.5,
					Value:    10.0,
				},
				{
					Quantile: 1.0,
					Value:    15.0,
				},
			},
		},
	}
	if len(points) == 0 {
		points = defaultPoints
	}
	return &opentelemetrypb.Metric{
		Name: name,
		Data: &opentelemetrypb.Metric_Summary{
			Summary: &opentelemetrypb.Summary{
				DataPoints: points,
			},
		},
	}
}

func generateOTLPSamples(srcs ...*opentelemetrypb.Metric) *opentelemetrypb.ResourceMetrics {
	otlpMetrics := &opentelemetrypb.ResourceMetrics{
		Resource: opentelemetrypb.Resource{Attributes: attributesFromKV([2]string{"job", "vm"})},
	}
	otlpMetrics.ScopeMetrics = []*opentelemetrypb.ScopeMetrics{
		{
			Metrics: append([]*opentelemetrypb.Metric{}, srcs...),
		},
	}
	return otlpMetrics
}

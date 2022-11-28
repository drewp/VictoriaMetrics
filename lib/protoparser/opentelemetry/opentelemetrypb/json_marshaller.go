package opentelemetrypb

import (
	"encoding/base64"
	"fmt"
	"strconv"

	jsoniter "github.com/json-iterator/go"
)

func (m *ExportMetricsServiceRequest) UnmarshalJSONExportMetricsServiceRequest(buf []byte) error {
	iter := jsoniter.ConfigFastest.BorrowIterator(buf)
	defer jsoniter.ConfigFastest.ReturnIterator(iter)
	iter.ReadObjectCB(func(iter *jsoniter.Iterator, f string) bool {
		switch f {
		case "resource_metrics", "resourceMetrics":
			iter.ReadArrayCB(func(iterator *jsoniter.Iterator) bool {
				m.ResourceMetrics = append(m.ResourceMetrics, readResourceMetrics(iter))
				return true
			})
		default:
			iter.Skip()
		}
		return true
	})
	for _, rm := range m.ResourceMetrics {
		if len(rm.ScopeMetrics) == 0 {
			rm.ScopeMetrics = rm.DeprecatedScopeMetrics
		}
		rm.DeprecatedScopeMetrics = nil
	}
	return iter.Error
}

func readResourceMetrics(iter *jsoniter.Iterator) *ResourceMetrics {
	rs := &ResourceMetrics{}
	iter.ReadObjectCB(func(iter *jsoniter.Iterator, f string) bool {
		switch f {
		case "resource":
			ReadResource(iter, &rs.Resource)
		case "scopeMetrics", "scope_metrics":
			iter.ReadArrayCB(func(iter *jsoniter.Iterator) bool {
				rs.ScopeMetrics = append(rs.ScopeMetrics,
					readScopeMetrics(iter))
				return true
			})
		case "schemaUrl", "schema_url":
			rs.SchemaUrl = iter.ReadString()
		default:
			iter.Skip()
		}
		return true
	})
	return rs
}

func readScopeMetrics(iter *jsoniter.Iterator) *ScopeMetrics {
	ils := &ScopeMetrics{}
	iter.ReadObjectCB(func(iter *jsoniter.Iterator, f string) bool {
		switch f {
		case "scope":
			ReadScope(iter, &ils.Scope)
		case "metrics":
			iter.ReadArrayCB(func(iter *jsoniter.Iterator) bool {
				ils.Metrics = append(ils.Metrics, readMetric(iter))
				return true
			})
		case "schemaUrl", "schema_url":
			ils.SchemaUrl = iter.ReadString()
		default:
			iter.Skip()
		}
		return true
	})
	return ils
}

func readMetric(iter *jsoniter.Iterator) *Metric {
	sp := &Metric{}
	iter.ReadObjectCB(func(iter *jsoniter.Iterator, f string) bool {
		switch f {
		case "name":
			sp.Name = iter.ReadString()
		case "description":
			sp.Description = iter.ReadString()
		case "unit":
			sp.Unit = iter.ReadString()
		case "sum":
			sp.Data = readSumMetric(iter)
		case "gauge":
			sp.Data = readGaugeMetric(iter)
		case "histogram":
			sp.Data = readHistogramMetric(iter)
		case "exponential_histogram", "exponentialHistogram":
			sp.Data = readExponentialHistogramMetric(iter)
		case "summary":
			sp.Data = readSummaryMetric(iter)
		default:
			iter.Skip()
		}
		return true
	})
	return sp
}

func readSumMetric(iter *jsoniter.Iterator) *Metric_Sum {
	data := &Metric_Sum{
		Sum: &Sum{},
	}
	iter.ReadObjectCB(func(iter *jsoniter.Iterator, f string) bool {
		switch f {
		case "aggregation_temporality", "aggregationTemporality":
			data.Sum.AggregationTemporality = readAggregationTemporality(iter)
		case "is_monotonic", "isMonotonic":
			data.Sum.IsMonotonic = iter.ReadBool()
		case "data_points", "dataPoints":
			var dataPoints []*NumberDataPoint
			iter.ReadArrayCB(func(iter *jsoniter.Iterator) bool {
				dataPoints = append(dataPoints, readNumberDataPoint(iter))
				return true
			})
			data.Sum.DataPoints = dataPoints
		default:
			iter.Skip()
		}
		return true
	})
	return data
}

func readGaugeMetric(iter *jsoniter.Iterator) *Metric_Gauge {
	data := &Metric_Gauge{
		Gauge: &Gauge{},
	}
	iter.ReadObjectCB(func(iter *jsoniter.Iterator, f string) bool {
		switch f {
		case "data_points", "dataPoints":
			var dataPoints []*NumberDataPoint
			iter.ReadArrayCB(func(iter *jsoniter.Iterator) bool {
				dataPoints = append(dataPoints, readNumberDataPoint(iter))
				return true
			})
			data.Gauge.DataPoints = dataPoints
		default:
			iter.Skip()
		}
		return true
	})
	return data
}

func readHistogramMetric(iter *jsoniter.Iterator) *Metric_Histogram {
	data := &Metric_Histogram{
		Histogram: &Histogram{},
	}
	iter.ReadObjectCB(func(iter *jsoniter.Iterator, f string) bool {
		switch f {
		case "data_points", "dataPoints":
			var dataPoints []*HistogramDataPoint
			iter.ReadArrayCB(func(iter *jsoniter.Iterator) bool {
				dataPoints = append(dataPoints, readHistogramDataPoint(iter))
				return true
			})
			data.Histogram.DataPoints = dataPoints
		case "aggregation_temporality", "aggregationTemporality":
			data.Histogram.AggregationTemporality = readAggregationTemporality(iter)
		default:
			iter.Skip()
		}
		return true
	})
	return data
}

func readExponentialHistogramMetric(iter *jsoniter.Iterator) *Metric_ExponentialHistogram {
	data := &Metric_ExponentialHistogram{
		ExponentialHistogram: &ExponentialHistogram{},
	}
	iter.ReadObjectCB(func(iter *jsoniter.Iterator, f string) bool {
		switch f {
		case "data_points", "dataPoints":
			iter.ReadArrayCB(func(iter *jsoniter.Iterator) bool {
				data.ExponentialHistogram.DataPoints = append(data.ExponentialHistogram.DataPoints,
					readExponentialHistogramDataPoint(iter))
				return true
			})
		case "aggregation_temporality", "aggregationTemporality":
			data.ExponentialHistogram.AggregationTemporality = readAggregationTemporality(iter)
		default:
			iter.Skip()
		}
		return true
	})
	return data
}

func readSummaryMetric(iter *jsoniter.Iterator) *Metric_Summary {
	data := &Metric_Summary{
		Summary: &Summary{},
	}
	iter.ReadObjectCB(func(iter *jsoniter.Iterator, f string) bool {
		switch f {
		case "data_points", "dataPoints":
			iter.ReadArrayCB(func(iter *jsoniter.Iterator) bool {
				data.Summary.DataPoints = append(data.Summary.DataPoints,
					readSummaryDataPoint(iter))
				return true
			})
		default:
			iter.Skip()
		}
		return true
	})
	return data
}

func readExemplar(iter *jsoniter.Iterator) Exemplar {
	exemplar := Exemplar{}
	iter.ReadObjectCB(func(iter *jsoniter.Iterator, f string) bool {
		switch f {
		case "filtered_attributes", "filteredAttributes":
			iter.ReadArrayCB(func(iter *jsoniter.Iterator) bool {
				exemplar.FilteredAttributes = append(exemplar.FilteredAttributes, ReadAttribute(iter))
				return true
			})
		case "timeUnixNano", "time_unix_nano":
			exemplar.TimeUnixNano = ReadUint64(iter)
		case "as_int", "asInt":
			exemplar.Value = &Exemplar_AsInt{
				AsInt: ReadInt64(iter),
			}
		case "as_double", "asDouble":
			exemplar.Value = &Exemplar_AsDouble{
				AsDouble: ReadFloat64(iter),
			}
		case "traceId", "trace_id":
			if err := exemplar.TraceId.UnmarshalJSON([]byte(iter.ReadString())); err != nil {
				iter.ReportError("exemplar.traceId", fmt.Sprintf("parse trace_id:%v", err))
			}
		case "spanId", "span_id":
			if err := exemplar.SpanId.UnmarshalJSON([]byte(iter.ReadString())); err != nil {
				iter.ReportError("exemplar.spanId", fmt.Sprintf("parse span_id:%v", err))
			}
		default:
			iter.Skip()
		}
		return true
	})
	return exemplar
}

func readNumberDataPoint(iter *jsoniter.Iterator) *NumberDataPoint {
	point := &NumberDataPoint{}
	iter.ReadObjectCB(func(iter *jsoniter.Iterator, f string) bool {
		switch f {
		case "timeUnixNano", "time_unix_nano":
			point.TimeUnixNano = ReadUint64(iter)
		case "start_time_unix_nano", "startTimeUnixNano":
			point.StartTimeUnixNano = ReadUint64(iter)
		case "as_int", "asInt":
			point.Value = &NumberDataPoint_AsInt{
				AsInt: ReadInt64(iter),
			}
		case "as_double", "asDouble":
			point.Value = &NumberDataPoint_AsDouble{
				AsDouble: ReadFloat64(iter),
			}
		case "attributes":
			iter.ReadArrayCB(func(iter *jsoniter.Iterator) bool {
				point.Attributes = append(point.Attributes, ReadAttribute(iter))
				return true
			})
		case "exemplars":
			iter.ReadArrayCB(func(iter *jsoniter.Iterator) bool {
				point.Exemplars = append(point.Exemplars, readExemplar(iter))
				return true
			})
		case "flags":
			point.Flags = ReadUint32(iter)
		default:
			iter.Skip()
		}
		return true
	})
	return point
}

func readHistogramDataPoint(iter *jsoniter.Iterator) *HistogramDataPoint {
	point := &HistogramDataPoint{}
	iter.ReadObjectCB(func(iter *jsoniter.Iterator, f string) bool {
		switch f {
		case "timeUnixNano", "time_unix_nano":
			point.TimeUnixNano = ReadUint64(iter)
		case "start_time_unix_nano", "startTimeUnixNano":
			point.StartTimeUnixNano = ReadUint64(iter)
		case "attributes":
			iter.ReadArrayCB(func(iter *jsoniter.Iterator) bool {
				point.Attributes = append(point.Attributes, ReadAttribute(iter))
				return true
			})
		case "count":
			point.Count = ReadUint64(iter)
		case "sum":
			point.Sum_ = &HistogramDataPoint_Sum{Sum: ReadFloat64(iter)}
		case "bucket_counts", "bucketCounts":
			iter.ReadArrayCB(func(iter *jsoniter.Iterator) bool {
				point.BucketCounts = append(point.BucketCounts, ReadUint64(iter))
				return true
			})
		case "explicit_bounds", "explicitBounds":
			iter.ReadArrayCB(func(iter *jsoniter.Iterator) bool {
				point.ExplicitBounds = append(point.ExplicitBounds, ReadFloat64(iter))
				return true
			})
		case "exemplars":
			iter.ReadArrayCB(func(iter *jsoniter.Iterator) bool {
				point.Exemplars = append(point.Exemplars, readExemplar(iter))
				return true
			})
		case "flags":
			point.Flags = ReadUint32(iter)
		case "max":
			point.Max_ = &HistogramDataPoint_Max{
				Max: ReadFloat64(iter),
			}
		case "min":
			point.Min_ = &HistogramDataPoint_Min{
				Min: ReadFloat64(iter),
			}
		default:
			iter.Skip()
		}
		return true
	})
	return point
}

func readExponentialHistogramDataPoint(iter *jsoniter.Iterator) *ExponentialHistogramDataPoint {
	point := &ExponentialHistogramDataPoint{}
	iter.ReadObjectCB(func(iter *jsoniter.Iterator, f string) bool {
		switch f {
		case "timeUnixNano", "time_unix_nano":
			point.TimeUnixNano = ReadUint64(iter)
		case "start_time_unix_nano", "startTimeUnixNano":
			point.StartTimeUnixNano = ReadUint64(iter)
		case "attributes":
			iter.ReadArrayCB(func(iter *jsoniter.Iterator) bool {
				point.Attributes = append(point.Attributes, ReadAttribute(iter))
				return true
			})
		case "count":
			point.Count = ReadUint64(iter)
		case "sum":
			point.Sum_ = &ExponentialHistogramDataPoint_Sum{
				Sum: ReadFloat64(iter),
			}
		case "scale":
			point.Scale = iter.ReadInt32()
		case "zero_count", "zeroCount":
			point.ZeroCount = ReadUint64(iter)
		case "positive":
			point.Positive = readExponentialHistogramBuckets(iter)
		case "negative":
			point.Negative = readExponentialHistogramBuckets(iter)
		case "exemplars":
			iter.ReadArrayCB(func(iter *jsoniter.Iterator) bool {
				point.Exemplars = append(point.Exemplars, readExemplar(iter))
				return true
			})
		case "flags":
			point.Flags = ReadUint32(iter)
		case "max":
			point.Max_ = &ExponentialHistogramDataPoint_Max{
				Max: ReadFloat64(iter),
			}
		case "min":
			point.Min_ = &ExponentialHistogramDataPoint_Min{
				Min: ReadFloat64(iter),
			}
		default:
			iter.Skip()
		}
		return true
	})
	return point
}

func readSummaryDataPoint(iter *jsoniter.Iterator) *SummaryDataPoint {
	point := &SummaryDataPoint{}
	iter.ReadObjectCB(func(iter *jsoniter.Iterator, f string) bool {
		switch f {
		case "timeUnixNano", "time_unix_nano":
			point.TimeUnixNano = ReadUint64(iter)
		case "start_time_unix_nano", "startTimeUnixNano":
			point.StartTimeUnixNano = ReadUint64(iter)
		case "attributes":
			iter.ReadArrayCB(func(iter *jsoniter.Iterator) bool {
				point.Attributes = append(point.Attributes, ReadAttribute(iter))
				return true
			})
		case "count":
			point.Count = ReadUint64(iter)
		case "sum":
			point.Sum = ReadFloat64(iter)
		case "quantile_values", "quantileValues":
			iter.ReadArrayCB(func(iter *jsoniter.Iterator) bool {
				point.QuantileValues = append(point.QuantileValues, readQuantileValue(iter))
				return true
			})
		case "flags":
			point.Flags = ReadUint32(iter)
		default:
			iter.Skip()
		}
		return true
	})
	return point
}

func readExponentialHistogramBuckets(iter *jsoniter.Iterator) ExponentialHistogramDataPoint_Buckets {
	buckets := ExponentialHistogramDataPoint_Buckets{}
	iter.ReadObjectCB(func(iter *jsoniter.Iterator, f string) bool {
		switch f {
		case "bucket_counts", "bucketCounts":
			iter.ReadArrayCB(func(iter *jsoniter.Iterator) bool {
				buckets.BucketCounts = append(buckets.BucketCounts, ReadUint64(iter))
				return true
			})
		case "offset":
			buckets.Offset = iter.ReadInt32()
		default:
			iter.Skip()
		}
		return true
	})
	return buckets
}

func readQuantileValue(iter *jsoniter.Iterator) *SummaryDataPoint_ValueAtQuantile {
	point := &SummaryDataPoint_ValueAtQuantile{}
	iter.ReadObjectCB(func(iter *jsoniter.Iterator, f string) bool {
		switch f {
		case "quantile":
			point.Quantile = ReadFloat64(iter)
		case "value":
			point.Value = ReadFloat64(iter)
		default:
			iter.Skip()
		}
		return true
	})
	return point
}

func readAggregationTemporality(iter *jsoniter.Iterator) AggregationTemporality {
	return AggregationTemporality(ReadEnumValue(iter, AggregationTemporality_value))
}

func ReadResource(iter *jsoniter.Iterator, resource *Resource) {
	iter.ReadObjectCB(func(iter *jsoniter.Iterator, f string) bool {
		switch f {
		case "attributes":
			iter.ReadArrayCB(func(iter *jsoniter.Iterator) bool {
				resource.Attributes = append(resource.Attributes, ReadAttribute(iter))
				return true
			})
		case "droppedAttributesCount", "dropped_attributes_count":
			resource.DroppedAttributesCount = ReadUint32(iter)
		default:
			iter.Skip()
		}
		return true
	})
}

// ReadInt32 unmarshalls JSON data into an int32. Accepts both numbers and strings decimal.
// See https://developers.google.com/protocol-buffers/docs/proto3#
func ReadInt32(iter *jsoniter.Iterator) int32 {
	switch iter.WhatIsNext() {
	case jsoniter.NumberValue:
		return iter.ReadInt32()
	case jsoniter.StringValue:
		val, err := strconv.ParseInt(iter.ReadString(), 10, 32)
		if err != nil {
			iter.ReportError("ReadInt32", err.Error())
			return 0
		}
		return int32(val)
	default:
		iter.ReportError("ReadInt32", "unsupported value type")
		return 0
	}
}

// ReadUint32 unmarshalls JSON data into an uint32. Accepts both numbers and strings decimal.
// See https://developers.google.com/protocol-buffers/docs/proto3#
func ReadUint32(iter *jsoniter.Iterator) uint32 {
	switch iter.WhatIsNext() {
	case jsoniter.NumberValue:
		return iter.ReadUint32()
	case jsoniter.StringValue:
		val, err := strconv.ParseUint(iter.ReadString(), 10, 32)
		if err != nil {
			iter.ReportError("ReadUint32", err.Error())
			return 0
		}
		return uint32(val)
	default:
		iter.ReportError("ReadUint32", "unsupported value type")
		return 0
	}
}

// ReadInt64 unmarshalls JSON data into an int64. Accepts both numbers and strings decimal.
// See https://developers.google.com/protocol-buffers/docs/proto3#
func ReadInt64(iter *jsoniter.Iterator) int64 {
	switch iter.WhatIsNext() {
	case jsoniter.NumberValue:
		return iter.ReadInt64()
	case jsoniter.StringValue:
		val, err := strconv.ParseInt(iter.ReadString(), 10, 64)
		if err != nil {
			iter.ReportError("ReadInt64", err.Error())
			return 0
		}
		return val
	default:
		iter.ReportError("ReadInt64", "unsupported value type")
		return 0
	}
}

// ReadUint64 unmarshalls JSON data into an uint64. Accepts both numbers and strings decimal.
// See https://developers.google.com/protocol-buffers/docs/proto3#
func ReadUint64(iter *jsoniter.Iterator) uint64 {
	switch iter.WhatIsNext() {
	case jsoniter.NumberValue:
		return iter.ReadUint64()
	case jsoniter.StringValue:
		val, err := strconv.ParseUint(iter.ReadString(), 10, 64)
		if err != nil {
			iter.ReportError("ReadUint64", err.Error())
			return 0
		}
		return val
	default:
		iter.ReportError("ReadUint64", "unsupported value type")
		return 0
	}
}

func ReadFloat64(iter *jsoniter.Iterator) float64 {
	switch iter.WhatIsNext() {
	case jsoniter.NumberValue:
		return iter.ReadFloat64()
	case jsoniter.StringValue:
		val, err := strconv.ParseFloat(iter.ReadString(), 64)
		if err != nil {
			iter.ReportError("ReadUint64", err.Error())
			return 0
		}
		return val
	default:
		iter.ReportError("ReadUint64", "unsupported value type")
		return 0
	}
}

func ReadScope(iter *jsoniter.Iterator, scope *InstrumentationScope) {
	iter.ReadObjectCB(func(iter *jsoniter.Iterator, f string) bool {
		switch f {
		case "name":
			scope.Name = iter.ReadString()
		case "version":
			scope.Version = iter.ReadString()
		case "attributes":
			iter.ReadArrayCB(func(iter *jsoniter.Iterator) bool {
				scope.Attributes = append(scope.Attributes, ReadAttribute(iter))
				return true
			})
		case "droppedAttributesCount", "dropped_attributes_count":
			scope.DroppedAttributesCount = ReadUint32(iter)
		default:
			iter.Skip()
		}
		return true
	})
}

// ReadAttribute Unmarshal JSON data and return otlpcommon.KeyValue
func ReadAttribute(iter *jsoniter.Iterator) KeyValue {
	kv := KeyValue{}
	iter.ReadObjectCB(func(iter *jsoniter.Iterator, f string) bool {
		switch f {
		case "key":
			kv.Key = iter.ReadString()
		case "value":
			ReadValue(iter, &kv.Value)
		default:
			iter.Skip()
		}
		return true
	})
	return kv
}

// ReadValue Unmarshal JSON data and return otlpcommon.AnyValue
func ReadValue(iter *jsoniter.Iterator, val *AnyValue) {
	iter.ReadObjectCB(func(iter *jsoniter.Iterator, f string) bool {
		switch f {
		case "stringValue", "string_value":
			val.Value = &AnyValue_StringValue{
				StringValue: iter.ReadString(),
			}

		case "boolValue", "bool_value":
			val.Value = &AnyValue_BoolValue{
				BoolValue: iter.ReadBool(),
			}
		case "intValue", "int_value":
			val.Value = &AnyValue_IntValue{
				IntValue: ReadInt64(iter),
			}
		case "doubleValue", "double_value":
			val.Value = &AnyValue_DoubleValue{
				DoubleValue: ReadFloat64(iter),
			}
		case "bytesValue", "bytes_value":
			v, err := base64.StdEncoding.DecodeString(iter.ReadString())
			if err != nil {
				iter.ReportError("bytesValue", fmt.Sprintf("base64 decode:%v", err))
				break
			}
			val.Value = &AnyValue_BytesValue{
				BytesValue: v,
			}
		case "arrayValue", "array_value":
			val.Value = &AnyValue_ArrayValue{
				ArrayValue: readArray(iter),
			}
		case "kvlistValue", "kvlist_value":
			val.Value = &AnyValue_KvlistValue{
				KvlistValue: readKvlistValue(iter),
			}
		default:
			iter.Skip()
		}
		return true
	})
}

func readArray(iter *jsoniter.Iterator) *ArrayValue {
	v := &ArrayValue{}
	iter.ReadObjectCB(func(iter *jsoniter.Iterator, f string) bool {
		switch f {
		case "values":
			iter.ReadArrayCB(func(iter *jsoniter.Iterator) bool {
				v.Values = append(v.Values, AnyValue{})
				ReadValue(iter, &v.Values[len(v.Values)-1])
				return true
			})
		default:
			iter.Skip()
		}
		return true
	})
	return v
}

func readKvlistValue(iter *jsoniter.Iterator) *KeyValueList {
	v := &KeyValueList{}
	iter.ReadObjectCB(func(iter *jsoniter.Iterator, f string) bool {
		switch f {
		case "values":
			iter.ReadArrayCB(func(iter *jsoniter.Iterator) bool {
				v.Values = append(v.Values, ReadAttribute(iter))
				return true
			})
		default:
			iter.Skip()
		}
		return true
	})
	return v
}

// ReadEnumValue returns the enum integer value representation. Accepts both enum names and enum integer values.
// See https://developers.google.com/protocol-buffers/docs/proto3#json.
func ReadEnumValue(iter *jsoniter.Iterator, valueMap map[string]int32) int32 {
	switch iter.WhatIsNext() {
	case jsoniter.NumberValue:
		return iter.ReadInt32()
	case jsoniter.StringValue:
		val, ok := valueMap[iter.ReadString()]
		// Same behavior with official protbuf JSON decoder,
		// see https://github.com/open-telemetry/opentelemetry-proto-go/pull/81
		if !ok {
			iter.ReportError("ReadEnumValue", "unknown string value")
			return 0
		}
		return val
	default:
		iter.ReportError("ReadEnumValue", "unsupported value type")
		return 0
	}
}

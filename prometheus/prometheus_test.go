package prometheus

import (
	"fmt"
	"github.com/pixie79/data-utils/utils"
	"testing"
)

// TestSplitTags tests the SplitTags function
func TestSplitTags(t *testing.T) {
	data := `a="b",c="d",e="f"`
	result := SplitTags(data)

	utils.Logger.Info(fmt.Sprintf("%+v", result))
	if len(result) != 5 {
		t.Errorf("expected 5 tags, got %d", len(result))
	}

	if result[0].Name != "a" || result[0].Value != "b" {
		t.Errorf("expected 'a:b', got '%s:%s'", result[0].Name, result[0].Value)
	}

	if result[1].Name != "c" || result[1].Value != "d" {
		t.Errorf("expected 'c:d', got '%s:%s'", result[1].Name, result[1].Value)
	}

	if result[2].Name != "e" || result[2].Value != "f" {
		t.Errorf("expected 'e:f', got '%s:%s'", result[2].Name, result[2].Value)
	}
}

// TestBuildMetrics tests the BuildMetrics function
func TestBuildMetrics(t *testing.T) {
	payload := []string{
		`# HELP redpanda_cloud_storage_segments_pending_deletion Total number of segments pending deletion from the cloud for the topic`,
		`# TYPE redpanda_cloud_storage_segments_pending_deletion gauge`,
		`redpanda_cloud_storage_segments_pending_deletion{redpanda_namespace="kafka",redpanda_topic="test",instance="10.0.0.240:9644"} 0`,
		`redpanda_cloud_storage_segments_pending_deletion{redpanda_namespace="kafka",redpanda_topic="s3-connector-10x-dlq-no",instance="10.0.0.240:9644"} 1`,
		`redpanda_cloud_storage_segments_pending_deletion{redpanda_namespace="kafka",redpanda_topic="s3-connector-10x-dlq-ab",instance="10.0.0.240:9644"} 5.4`,
	}

	result := BuildMetrics(payload)

	if len(result) != 3 {
		t.Errorf("expected 3 metrics, got %d", len(result))
	}

	if result[0].Metric != `redpanda.redpanda_cloud_storage_segments_pending_deletion` {
		t.Errorf("expected metric 'redpanda.redpanda_cloud_storage_segments_pending_deletion', got '%s'", result[0].Metric)
	}

	if result[0].Tags[0] != `redpanda_namespace=kafka` {
		t.Errorf("expected tag 'redpanda_namespace=kafka', got '%s'", result[0].Tags[0])
	}

	if result[0].Tags[1] != `redpanda_topic=test` {
		t.Errorf(`expected tag 'redpanda_topic="test"', got '%s'`, result[0].Tags[1])
	}

	if result[0].Tags[2] != `instance=10.0.0.240:9644` {
		t.Errorf(`expected tag 'instance="10.0.0.240:9644', got '%s'`, result[0].Tags[2])
	}

	if result[1].Metric != `redpanda.redpanda_cloud_storage_segments_pending_deletion` {
		t.Errorf(`expected metric 'redpanda.redpanda_cloud_storage_segments_pending_deletion', got '%s'`, result[1].Metric)
	}

	if result[1].Tags[0] != `redpanda_namespace=kafka` {
		t.Errorf(`expected tag 'redpanda_namespace=kafka', got '%s'`, result[1].Tags[0])
	}

	if result[2].Metric != `redpanda.redpanda_cloud_storage_segments_pending_deletion` {
		t.Errorf(`expected metric 'redpanda.redpanda_cloud_storage_segments_pending_deletion', got '%s'`, result[2].Metric)
	}

	if result[2].Tags[0] != `redpanda_namespace=kafka` {
		t.Errorf(`expected tag 'redpanda_namespace=kafka', got '%s'`, result[2].Tags[0])
	}
}

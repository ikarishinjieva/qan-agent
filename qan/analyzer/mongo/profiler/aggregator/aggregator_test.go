package aggregator

import (
	"testing"
	"time"

	"github.com/percona/go-mysql/event"
	"github.com/percona/percona-toolkit/src/go/mongolib/proto"
	pc "github.com/percona/pmm/proto/config"
	"github.com/percona/pmm/proto/qan"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAggregator_Add(t *testing.T) {
	t.Parallel()

	timeStart, err := time.Parse("2006-01-02 15:04:05", "2017-07-02 07:55:00")
	require.NoError(t, err)
	timeEnd, err := time.Parse("2006-01-02 15:04:05", "2017-07-02 07:56:00")
	require.NoError(t, err)

	config := pc.QAN{
		UUID:     "abc",
		Interval: 60, // 60s
	}

	aggregator := New(timeStart, config)

	{
		doc := proto.SystemProfile{
			Ts: timeStart,
		}
		report, err := aggregator.Add(doc)
		require.NoError(t, err)
		assert.Nil(t, report)
	}

	{
		doc := proto.SystemProfile{
			Ts: timeEnd,
		}
		expected := qan.Report{
			UUID:    config.UUID,
			StartTs: timeStart,
			EndTs:   timeEnd,
			Global: &event.Class{
				TotalQueries:  1,
				UniqueQueries: 1,
				Metrics: &event.Metrics{
					TimeMetrics: map[string]*event.TimeStats{
						"Query_time": {},
					},
					NumberMetrics: map[string]*event.NumberStats{
						"Bytes_sent":    {},
						"Rows_examined": {},
						"Rows_sent":     {},
					},
					BoolMetrics: map[string]*event.BoolStats{},
				},
			},
			Class: []*event.Class{
				{
					Id:            "d41d8cd98f00b204e9800998ecf8427e",
					TotalQueries:  1,
					UniqueQueries: 1,
					Metrics: &event.Metrics{
						TimeMetrics: map[string]*event.TimeStats{
							"Query_time": {},
						},
						NumberMetrics: map[string]*event.NumberStats{
							"Bytes_sent":    {},
							"Rows_examined": {},
							"Rows_sent":     {},
						},
						BoolMetrics: map[string]*event.BoolStats{},
					},
					Example: &event.Example{},
				},
			},
		}
		report, err := aggregator.Add(doc)
		require.NoError(t, err)
		assert.Equal(t, expected, *report)
	}
}

// TestAggregator_Add_EmptyInterval verifies that no report is returned if there were no samples in interval #PMM-927
func TestAggregator_Add_EmptyInterval(t *testing.T) {
	t.Parallel()

	timeStart, err := time.Parse("2006-01-02 15:04:05", "2017-07-02 07:55:00")
	require.NoError(t, err)
	timeEnd, err := time.Parse("2006-01-02 15:04:05", "2017-07-02 07:56:00")
	require.NoError(t, err)

	config := pc.QAN{
		UUID:     "abc",
		Interval: 60, // 60s
	}

	aggregator := New(timeStart, config)

	// finish interval immediately
	{
		doc := proto.SystemProfile{
			Ts: timeEnd,
		}
		report, err := aggregator.Add(doc)
		require.NoError(t, err)

		// no report should be returned
		assert.Nil(t, report)
	}
}

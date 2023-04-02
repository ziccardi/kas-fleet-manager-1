package integration

import (
	"fmt"
	kafkatest "github.com/bf2fc6cc711aee1a0c2a/kas-fleet-manager/internal/kafka/test"
	"github.com/bf2fc6cc711aee1a0c2a/kas-fleet-manager/pkg/shared/utils/arrays"
	sync2 "github.com/bf2fc6cc711aee1a0c2a/kas-fleet-manager/pkg/shared/utils/sync"
	"github.com/bf2fc6cc711aee1a0c2a/kas-fleet-manager/test/mocks"
	"github.com/onsi/gomega"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

type testRoutineParams struct {
	wg                     *sync.WaitGroup
	distributedLockManager sync2.DistributedLockMgr
	delay                  time.Duration
	lockID                 string
	current                *int32
	dest                   *int32
	g                      *gomega.WithT
}

func getNext(p *testRoutineParams) {
	p.g.Expect(p.distributedLockManager.Lock(p.lockID)).To(gomega.Succeed())

	*p.dest = *p.current + 1
	time.Sleep(p.delay)
	atomic.AddInt32(p.current, 1)

	p.g.Expect(p.distributedLockManager.Unlock(p.lockID)).To(gomega.Succeed())

	p.wg.Done()
}

func TestDistributedLockMgr_ConcurrentAccess(t *testing.T) {
	ocmServer := mocks.NewMockConfigurableServerBuilder().Build()
	defer ocmServer.Close()

	h, _, teardown := kafkatest.NewKafkaHelper(t, ocmServer)
	defer teardown()

	type routineTestConfig struct {
		lockID      string
		outputValue int32
		delay       time.Duration
	}

	mostFrequent := func(arr []routineTestConfig) (int32, int32) {
		m := map[int32]int32{}
		var maxCnt int32
		var mostFrequentValue int32
		var freq int32
		for _, r := range arr {
			m[r.outputValue]++
			if m[r.outputValue] > maxCnt {
				maxCnt = m[r.outputValue]
				freq = m[r.outputValue]
				mostFrequentValue = r.outputValue
			}
		}
		return mostFrequentValue, freq
	}

	createRoutineConfig := func(lockIDGen func() string, delay func() time.Duration, count int) []routineTestConfig {
		ret := make([]routineTestConfig, count)

		for i := 0; i < count; i++ {
			ret[i] = routineTestConfig{
				lockID: lockIDGen(),
				delay:  delay(),
			}
		}

		return ret
	}

	counter := 0

	tests := []struct {
		name                   string
		reset                  func()
		startWith              int32
		routines               []routineTestConfig
		minimumResultFrequency int
		maximumResultFrequency int
	}{
		{
			name:                   "All thread same lock",
			startWith:              0,
			minimumResultFrequency: 1, // each thread should have been executed
			maximumResultFrequency: 1, // no thread should have the same result of another thread
			routines:               createRoutineConfig(func() string { return "lock1" }, func() time.Duration { return 5 * time.Millisecond }, 100),
		},
		{
			name:                   "All thread different lock",
			startWith:              0,
			minimumResultFrequency: 2,  // at least two thread should have executed at the same time
			maximumResultFrequency: -1, // up to all thread could have been executed at the same time (depends on CPU load)
			routines:               createRoutineConfig(func() string { counter++; return fmt.Sprintf("lock%d", counter) }, func() time.Duration { return 500 * time.Millisecond }, 100),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			g := gomega.NewWithT(t)
			dbConn := h.DBFactory().New()
			distLockMgr := sync2.NewDistributedLockMgr(dbConn)

			routineCount := int32(len(tt.routines))

			var wg sync.WaitGroup
			wg.Add(int(routineCount))

			for i := range tt.routines {
				routine := &tt.routines[i]
				p := testRoutineParams{
					wg:                     &wg,
					distributedLockManager: distLockMgr,
					delay:                  routine.delay,
					lockID:                 routine.lockID,
					current:                &tt.startWith,
					dest:                   &routine.outputValue,
					g:                      g,
				}

				go getNext(&p)
			}

			wg.Wait()

			mostFrequentValue, frequency := mostFrequent(tt.routines)

			g.Expect(tt.startWith).To(gomega.Equal(routineCount))
			if tt.maximumResultFrequency != -1 {
				g.Expect(frequency).To(gomega.BeNumerically("<=", tt.maximumResultFrequency), "frequency: %d, %d", mostFrequentValue, frequency)
			}

			if tt.minimumResultFrequency != -1 {
				g.Expect(frequency).To(gomega.BeNumerically(">=", tt.minimumResultFrequency))
			}

			g.Expect(arrays.Contains(arrays.Map(tt.routines, func(x routineTestConfig) int32 { return x.outputValue }), 0)).To(gomega.BeFalse())
		})
	}
}

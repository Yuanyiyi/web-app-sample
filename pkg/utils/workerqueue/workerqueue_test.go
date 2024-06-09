package workerqueue

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/wait"
)

func TestWorkerQueueRun(t *testing.T) {
	t.Parallel()

	received := make(chan string)
	defer close(received)

	syncHandler := func(ctx context.Context, name string) error {
		if name == "test" {
			assert.Equal(t, "test", name)
			fmt.Println(name)
			received <- name
			return nil
		}
		fmt.Println("error: ", name)
		return errors.New("test")
	}

	wq := NewWorkerQueue(syncHandler, logrus.WithField("source", "test"), "testKey", "test")
	stop := make(chan struct{})
	defer close(stop)

	go wq.Run(context.Background(), 1)

	// no change, should be no value
	select {
	case <-received:
		assert.Fail(t, "should not have received value")
	case <-time.After(1 * time.Second):
	}

	wq.Enqueue("test")
	wq.Enqueue("test1")

	select {
	case <-received:
	case <-time.After(5 * time.Second):
		assert.Fail(t, "should have received value")
	}
}

func TestWorkerQueueHealthy(t *testing.T) {
	t.Parallel()

	done := make(chan struct{})
	handler := func(context.Context, string) error {
		<-done
		return nil
	}
	wq := NewWorkerQueue(handler, logrus.WithField("source", "test"), "testKey", "test")
	wq.Enqueue("test")

	ctx, cancel := context.WithCancel(context.Background())
	go wq.Run(ctx, 1)

	// Yield to the scheduler to ensure the worker queue goroutine can run.
	err := wait.Poll(100*time.Millisecond, 3*time.Second, func() (done bool, err error) {
		if (wq.RunCount() == 1) && wq.Healthy() == nil {
			return true, nil
		}

		return false, nil
	})
	assert.Nil(t, err)

	close(done) // Ensure the handler no longer blocks.
	cancel()    // Stop the worker queue.

	// Yield to the scheduler again to ensure the worker queue goroutine can
	// finish.
	err = wait.Poll(100*time.Millisecond, 3*time.Second, func() (done bool, err error) {
		if (wq.RunCount() == 0) && wq.Healthy() != nil {
			return true, nil
		}

		return false, nil
	})
	assert.Nil(t, err)
}

func TestWorkerQueueEnqueueAfter(t *testing.T) {
	t.Parallel()

	updated := make(chan bool)
	syncHandler := func(ctx context.Context, s string) error {
		updated <- true
		return nil
	}
	wq := NewWorkerQueue(syncHandler, logrus.WithField("source", "test"), "testKey", "test")
	stop := make(chan struct{})
	defer close(stop)

	go wq.Run(context.Background(), 1)

	wq.EnqueueAfter("test", 2*time.Second)

	select {
	case <-updated:
		assert.FailNow(t, "should not be a result in queue yet")
	case <-time.After(time.Second):
	}

	select {
	case <-updated:
	case <-time.After(2 * time.Second):
		assert.Fail(t, "should have got a queue'd message by now")
	}
}

func TestDebugError(t *testing.T) {
	err := errors.New("not a debug error")
	assert.False(t, isDebugError(err))

	err = NewDebugError(err)
	assert.True(t, isDebugError(err))
	assert.EqualError(t, err, "not a debug error")

	err = NewDebugError(nil)
	assert.True(t, isDebugError(err))
	assert.EqualError(t, err, "<nil>")
}

package workerqueue

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/time/rate"

	k8serror "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/util/workqueue"

	"github.com/web-app-sample/pkg/utils/runtime"
)

const (
	workFx = time.Millisecond * 1000
)

// debugError is a marker type for errors that that should only be logged at a Debug level.
// Useful if you want a Handler to be retried, but not logged at an Error level.
type debugError struct {
	err error
}

// NewDebugError returns a debugError wrapper around an error.
func NewDebugError(err error) error {
	return &debugError{err: err}
}

// Error returns the error string
func (l *debugError) Error() string {
	if l.err == nil {
		return "<nil>"
	}
	return l.err.Error()
}

// isDebugError returns if the error is a debug error or not
func isDebugError(err error) bool {
	cause := errors.Cause(err)
	_, ok := cause.(*debugError)
	return ok
}

// Handler is the handler for processing the work queue
// This is usually a syncronisation handler for a controller or related
type Handler func(context.Context, string) error

// WorkerQueue is an opinionated queue + worker for use
// with controllers and related and processing Kubernetes watched
// events and synchronising resources
//
//nolint:govet // ignore fieldalignment, singleton
type WorkerQueue struct {
	logger  *logrus.Entry
	keyName string
	queue   workqueue.RateLimitingInterface
	// SyncHandler is exported to make testing easier (hack)
	SyncHandler Handler

	mu      sync.Mutex
	workers int
	running int
}

// FastRateLimiter returns a rate limiter without exponential back-off, with specified maximum per-item retry delay.
func FastRateLimiter(maxDelay time.Duration) workqueue.RateLimiter {
	const numFastRetries = 5
	const fastDelay = 200 * time.Millisecond // first few retries up to 'numFastRetries' are fast

	return workqueue.NewItemFastSlowRateLimiter(fastDelay, maxDelay, numFastRetries)
}

func CustomizeRateLimiter() workqueue.RateLimiter {
	return workqueue.NewMaxOfRateLimiter(
		fastRateLimiter(),
		// 50 qps, 200 bucket size.  This is only for retry speed and its only the overall factor (not per item)
		&workqueue.BucketRateLimiter{Limiter: rate.NewLimiter(rate.Limit(500), 1000)},
	)
}

func DefaultRateLimiter() workqueue.RateLimiter {
	return workqueue.DefaultControllerRateLimiter()
}

// fastRateLimiter returns a fast rate limiter, without exponential back-off.
func fastRateLimiter() workqueue.RateLimiter {
	const numFastRetries = 5
	const fastDelay = 50 * time.Millisecond  // first few retries up to 'numFastRetries' are fast
	const slowDelay = 200 * time.Millisecond // subsequent retries are slow

	return workqueue.NewItemFastSlowRateLimiter(fastDelay, slowDelay, numFastRetries)
}

// NewWorkerQueue returns a new worker queue for a given name
func NewWorkerQueue(handler Handler, logger *logrus.Entry, keyName string, queueName string) *WorkerQueue {
	return NewWorkerQueueWithRateLimiter(handler, logger, keyName, queueName, workqueue.DefaultControllerRateLimiter())
}

// NewWorkerQueueWithRateLimiter returns a new worker queue for a given name and a custom rate limiter.
func NewWorkerQueueWithRateLimiter(handler Handler, logger *logrus.Entry, keyName string, queueName string, rateLimiter workqueue.RateLimiter) *WorkerQueue {
	return &WorkerQueue{
		keyName:     string(keyName),
		logger:      logger.WithField("queue", queueName),
		queue:       workqueue.NewNamedRateLimitingQueue(rateLimiter, queueName),
		SyncHandler: handler,
	}
}

// Enqueue puts the name of the runtime.Object in the
// queue to be processed. If you need to send through an
// explicit key, use an cache.ExplicitKey
func (wq *WorkerQueue) Enqueue(key string) {
	wq.logger.WithField(wq.keyName, key).Debug("Enqueuing")
	wq.queue.AddRateLimited(key)
}

// EnqueueImmediately performs Enqueue but without rate-limiting.
// This should be used to continue partially completed work after giving other
// items in the queue a chance of running.
func (wq *WorkerQueue) EnqueueImmediately(key string) {
	wq.logger.WithField(wq.keyName, key).Debug("Enqueuing immediately")
	wq.queue.Add(key)
}

// EnqueueAfter delays an enqueue operation by duration
func (wq *WorkerQueue) EnqueueAfter(key string, duration time.Duration) {
	wq.logger.WithField(wq.keyName, key).WithField("duration", duration).Debug("Enqueueing after duration")
	wq.queue.AddAfter(key, duration)
}

// runWorker is a long-running function that will continually call the
// processNextWorkItem function in order to read and process a message on the
// workqueue.
func (wq *WorkerQueue) runWorker(ctx context.Context) {
	for wq.processNextWorkItem(ctx) {
	}
}

// processNextWorkItem processes the next work item.
// pretty self explanatory :)
func (wq *WorkerQueue) processNextWorkItem(ctx context.Context) bool {
	obj, quit := wq.queue.Get()
	if quit {
		return false
	}
	defer wq.queue.Done(obj)

	wq.logger.WithField(wq.keyName, obj).Debug("Processing")

	var key string
	var ok bool
	if key, ok = obj.(string); !ok {
		runtime.HandleError(wq.logger.WithField(wq.keyName, obj), errors.Errorf("expected string in queue, but got %T", obj))
		// this is a bad entry, we don't want to reprocess
		wq.queue.Forget(obj)
		return true
	}

	if err := wq.SyncHandler(ctx, key); err != nil {
		// Conflicts are expected, so only show them in debug operations.
		// Also check is debugError for other expected errors.
		if k8serror.IsConflict(errors.Cause(err)) || isDebugError(err) {
			wq.logger.WithField(wq.keyName, obj).Debug(err)
		} else {
			runtime.HandleError(wq.logger.WithField(wq.keyName, obj), err)
		}

		// we don't forget here, because we want this to be retried via the queue
		wq.queue.AddRateLimited(obj)
		return true
	}

	wq.queue.Forget(obj)
	return true
}

// Run the WorkerQueue processing via the Handler. Will block until stop is closed.
// Runs a certain number workers to process the rate limited queue
func (wq *WorkerQueue) Run(ctx context.Context, workers int) {
	wq.setWorkerCount(workers)
	wq.logger.WithField("workers", workers).Info("Starting workers...")
	for i := 0; i < workers; i++ {
		go wq.run(ctx)
	}

	<-ctx.Done()
	wq.logger.Info("...shutting down workers")
	wq.queue.ShutDown()
}

func (wq *WorkerQueue) run(ctx context.Context) {
	wq.inc()
	defer wq.dec()
	wait.Until(func() { wq.runWorker(ctx) }, workFx, ctx.Done())
}

// Healthy reports whether all the worker goroutines are running.
func (wq *WorkerQueue) Healthy() error {
	wq.mu.Lock()
	defer wq.mu.Unlock()
	want := wq.workers
	got := wq.running

	if want != got {
		return fmt.Errorf("want %d worker goroutine(s), got %d", want, got)
	}
	return nil
}

// RunCount reports the number of running worker goroutines started by Run.
func (wq *WorkerQueue) RunCount() int {
	wq.mu.Lock()
	defer wq.mu.Unlock()
	return wq.running
}

func (wq *WorkerQueue) setWorkerCount(n int) {
	wq.mu.Lock()
	defer wq.mu.Unlock()
	wq.workers = n
}

func (wq *WorkerQueue) inc() {
	wq.mu.Lock()
	defer wq.mu.Unlock()
	wq.running++
}

func (wq *WorkerQueue) dec() {
	wq.mu.Lock()
	defer wq.mu.Unlock()
	wq.running--
}

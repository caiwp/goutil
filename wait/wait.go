package wait

import "time"

type ConditionFunc func() (done bool, err error)
type WaitFunc func(done <-chan struct{}) <-chan struct{}

func WaitFor(wait WaitFunc, fn ConditionFunc, done <-chan struct{}) error {
	ch := wait(done)
	for {
		_, open := <-ch
		ok, err := fn()
		if err != nil {
			return err
		}
		if ok {
			return nil
		}
		if !open {
			break
		}
	}
	return nil
}

func poller(interval, timeout time.Duration) WaitFunc {
	return WaitFunc(func(done <-chan struct{}) <-chan struct{} {
		ch := make(chan struct{})

		go func() {
			defer close(ch)

			tick := time.NewTicker(interval)
			defer tick.Stop()

			var after <-chan time.Time
			if timeout != 0 {
				timer := time.NewTimer(timeout)
				after = timer.C
				defer timer.Stop()
			}

			for {
				select {
				case <-tick.C:
					select {
					case ch <- struct{}{}:
					default:
					}
				case <-after:
					return
				case <-done:
					return
				}
			}
		}()

		return ch
	})
}

func Poll(interval, timeout time.Duration, condition ConditionFunc) error {
	return pollInterval(poller(interval, timeout), condition)
}

func pollInterval(wait WaitFunc, condition ConditionFunc) error {
	done := make(chan struct{})
	defer close(done)
	return WaitFor(wait, condition, done)
}

func PollImmediate(interval, timeout time.Duration, condition ConditionFunc) error {
	return pollImmediateInterval(poller(interval, timeout), condition)
}

func pollImmediateInterval(wait WaitFunc, condition ConditionFunc) error {
	done, err := condition()
	if err != nil {
		return err
	}
	if done {
		return nil
	}
	return pollInterval(wait, condition)
}

func PollInfinite(interval time.Duration, condition ConditionFunc) error {
	done := make(chan struct{})
	defer close(done)
	return PollUntil(interval, condition, done)
}

func PollImmediateInfinite(interval time.Duration, condition ConditionFunc) error {
	done, err := condition()
	if err != nil {
		return err
	}
	if done {
		return nil
	}
	return PollInfinite(interval, condition)
}

func PollUntil(interval time.Duration, condition ConditionFunc, stopCh <-chan struct{}) error {
	return WaitFor(poller(interval, 0), condition, stopCh)
}

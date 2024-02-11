package tx

import (
	"fmt"
	"sync"
)

// TODO: TEST!!!

type fn = func(map[any]any) error

type fnPair struct {
	logic  fn
	revert fn
}

type Tx struct {
	jobs []fnPair

	logicErrorsMu sync.Mutex
	logicErrors   []error

	revertErrorsMu sync.Mutex
	revertErrors   []error

	panicErrorsMu sync.Mutex
	panicErrors   []error
}

func NewTx() Tx {
	return Tx{
		jobs:           make([]fnPair, 0),
		logicErrorsMu:  sync.Mutex{},
		logicErrors:    make([]error, 0),
		revertErrorsMu: sync.Mutex{},
		revertErrors:   make([]error, 0),
	}
}

func (t *Tx) Add(logic, revert fn) {
	t.jobs = append(t.jobs, fnPair{logic, revert})
}

// Run
// bool in return says if the error is fatal
// design is pretty strange, but I haven't come up to the better one
func (t *Tx) Run() (error, bool) {
	wg := sync.WaitGroup{}
	wg.Add(len(t.jobs))

	logicWg := sync.WaitGroup{}
	logicWg.Add(len(t.jobs))

	for _, job := range t.jobs {
		job := job
		go func() {
			defer func() {
				if r := recover(); r != nil {
					t.panicErrorsMu.Lock()
					t.panicErrors = append(t.panicErrors, fmt.Errorf("panic: %v", r))
					t.panicErrorsMu.Unlock()
				}
				wg.Done()
			}()

			commonData := make(map[any]any)
			err := job.logic(commonData)
			if err != nil {
				t.logicErrorsMu.Lock()
				t.logicErrors = append(t.logicErrors, err)
				t.logicErrorsMu.Unlock()

				logicWg.Done()
				return
			}
			logicWg.Done()
			logicWg.Wait()

			t.logicErrorsMu.Lock()
			errorsCount := len(t.logicErrors)
			t.logicErrorsMu.Unlock()

			if errorsCount > 0 {
				err = job.revert(commonData)
				if err != nil {

					t.revertErrorsMu.Lock()
					t.revertErrors = append(t.revertErrors, err)
					t.revertErrorsMu.Unlock()

				}
			}
		}()
	}
	wg.Wait()

	if panicErrors := t.PanicErrors(); len(panicErrors) != 0 {
		return panicErrors[0], true
	}

	if revertErrors := t.RevertErrors(); len(revertErrors) != 0 {
		return revertErrors[0], true
	}

	if logicErrors := t.LogicErrors(); len(logicErrors) != 0 {
		return logicErrors[0], false
	}

	return nil, false
}

func (t *Tx) LogicErrors() []error {
	t.logicErrorsMu.Lock()
	res := make([]error, len(t.logicErrors))
	copy(res, t.logicErrors)
	t.logicErrorsMu.Unlock()
	return res
}

func (t *Tx) RevertErrors() []error {
	t.revertErrorsMu.Lock()
	res := make([]error, len(t.revertErrors))
	copy(res, t.revertErrors)
	t.revertErrorsMu.Unlock()
	return res
}

func (t *Tx) PanicErrors() []error {
	t.panicErrorsMu.Lock()
	res := make([]error, len(t.panicErrors))
	copy(res, t.panicErrors)
	t.panicErrorsMu.Unlock()
	return res
}

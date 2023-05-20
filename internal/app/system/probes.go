package system

import (
	"sync"
)

// ProbeFunc функция проверки работоспособности
type ProbeFunc func() error

var (
	readyzMx      sync.Mutex
	readyzProbes  []ProbeFunc
	healthzMx     sync.Mutex
	healthzProbes []ProbeFunc
)

// AddReadyzProbe добавление функции проверки достуности
func AddReadyzProbe(probe ProbeFunc) {
	readyzMx.Lock()
	defer readyzMx.Unlock()
	readyzProbes = append(readyzProbes, probe)
}

// AddHealthzProbe добавление функции проверки здоровья
func AddHealthzProbe(probe ProbeFunc) {
	healthzMx.Lock()
	defer healthzMx.Unlock()
	healthzProbes = append(healthzProbes, probe)
}

// Проверка валидности
func isSuccess(mx *sync.Mutex, probes []ProbeFunc) error {
	if mx != nil {
		mx.Lock()
		defer mx.Unlock()
	}
	for _, probe := range probes {
		err := probe()
		if err != nil {
			return err
		}
	}
	return nil
}

// Healthz проверка жизнеспособности сервиса
func Healthz() error {
	return isSuccess(&healthzMx, healthzProbes)
}

// Readyz проверка доступности сервиса
func Readyz() error {
	return isSuccess(&readyzMx, readyzProbes)
}

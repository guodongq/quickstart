package workerpool

func NewWorkerPool(size int32) *WorkerPool {
	return &WorkerPool{
		queue: make(chan struct{}, size),
	}
}

type WorkerPool struct {
	queue chan struct{}
}

func (w *WorkerPool) GetWorker() {
	w.queue <- struct{}{}
}

func (w *WorkerPool) ReleaseWorker() {
	<-w.queue
}

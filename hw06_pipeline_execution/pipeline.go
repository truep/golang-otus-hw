package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	for _, stage := range stages {
		if stage != nil {
			dataChannel := make(Bi)
			go worker(in, done, dataChannel)
			in = stage(dataChannel)
		}
	}
	return in
}

func worker(in In, done In, ch Bi) {
	defer func() {
		close(ch)
		<-in
	}()

	for {
		select {
		case v, ok := <-in:
			if !ok {
				return
			}
			select {
			case ch <- v:
			case <-done:
				return
			}
		case <-done:
			return
		}
	}
}

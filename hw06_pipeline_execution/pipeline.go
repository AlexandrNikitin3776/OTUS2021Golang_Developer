package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func wrapWithDone(in In, done In) Out {
	out := make(Bi)

	go func() {
		defer close(out)

		for {
			select {
			case <-done:
				return
			default:
			}

			select {
			case <-done:
				return
			case data, ok := <-in:
				if !ok {
					return
				}
				out <- data
			}
		}
	}()

	return out
}

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := wrapWithDone(in, done)
	for _, stage := range stages {
		out = stage(wrapWithDone(out, done))
	}
	return out
}

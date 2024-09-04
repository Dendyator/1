package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)
type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := in

	for _, s := range stages {
		out = run(out, done, s)
	}
	return out
}

func run(in In, done In, stage Stage) Out {
	inner := make(Bi)

	go func() {
		defer close(inner)

		for {
			select {
			case v, ok := <-in:
				if !ok {
					return
				}
				inner <- v
			case <-done:
				return
			}
		}

	}()

	return stage(inner)
}

package engage

import "context"

// kindActor blocks until finish it's job and pass down to next actor
func kindActor(ctx context.Context, input <-chan interface{}) {

	for {
		select {
		case <-ctx.Done():
		default:
			// fmt.Printf("%s\n", string(d.Body))
		}
	}
}

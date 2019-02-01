package engage

import (
	"context"
	"fmt"
)

// kindActor blocks until finish it's job and pass down to next actor
func kindActor(ctx context.Context, input <-chan interface{}) {

	for {
		select {
		case <-ctx.Done():
			return
		case d := <-input:
			// base on the input then dispatch to next actor
			if v, ok := d.([]byte); ok {
				fmt.Printf("%s\n", string(v))
			}
		}
	}
}

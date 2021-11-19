package pkg

import (
	"io"

	"github.com/latonaio/aion-core/pkg/log"
)

func SafeClose(closer io.Closer){
	if closer != nil {
		if err := closer.Close(); err != nil {
			log.Printf("failed to close: %v", err)
		}
	}
}
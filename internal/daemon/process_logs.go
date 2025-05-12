package daemon

import (
	"context"
	"sync"

	"github.com/scalarorg/scalar-healer/pkg/evm"
)

func (s *Service) ProcessMissingLogs(ctx context.Context) {
	var wg sync.WaitGroup
	for _, client := range s.EvmClients {
		wg.Add(1)
		go func(client *evm.EvmClient) {
			defer wg.Done()
			client.ProcessMissingLogs(ctx)
		}(client)
	}

	wg.Wait()
}

package tofnd

import (
	"context"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

const DEFAULT_TIMEOUT = 20 * time.Second

type Manager struct {
	clients   []*Client
	Threshold int
}

func NewManager(configPath string) *Manager {
	cfgs, err := ReadTofndClientConfig(configPath)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to read tofnd client configs")
	}
	clients := make([]*Client, len(cfgs))
	threshold := 0
	for i, cfg := range cfgs {
		clients[i], err = NewClient(&cfg, DEFAULT_TIMEOUT)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to create tofnd client")
		}
		threshold += cfg.Weight
	}

	log.Info().Msgf("Connected to %d tofnd clients", len(clients))

	threshold = threshold * 2 / 3

	return &Manager{
		clients:   clients,
		Threshold: threshold,
	}
}

type SignatureResult struct {
	Client *Client
	Sig    []byte
	Err    error
}

func (m *Manager) Sign(ctx context.Context, msg []byte) ([]SignatureResult, error) {
	var (
		mu          sync.Mutex
		results     []SignatureResult
		totalWeight int
		wg          sync.WaitGroup
		done        = make(chan struct{})
	)

	resultCh := make(chan SignatureResult, len(m.clients))

	for _, client := range m.clients {
		wg.Add(1)
		go func(client *Client) {
			defer wg.Done()

			sig, err := client.Sign(ctx, msg)
			select {
			case resultCh <- SignatureResult{
				Client: client,
				Sig:    sig.GetSignature(),
				Err:    err,
			}:
			case <-done:
			}
		}(client)
	}

	for {
		res := <-resultCh
		if res.Err == nil && len(res.Sig) >= 64 {
			mu.Lock()
			results = append(results, res)
			totalWeight += res.Client.Weight
			if totalWeight >= m.Threshold {
				close(done) // signal goroutines to stop sending
				mu.Unlock()
				break
			}
			mu.Unlock()
		} else {
			log.Error().Err(res.Err).Msg("Failed to sign message")
		}
	}

	go func() {
		wg.Wait()
		close(resultCh)
	}()

	return results, nil
}

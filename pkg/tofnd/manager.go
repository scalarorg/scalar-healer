package tofnd

import (
	"context"
	"encoding/json"
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

type Musig struct {
	Sig       Signature `json:"signature"`
	Weight    int       `json:"weight"`
	Threshold int       `json:"threshold"`
	PartyID   string    `json:"party_id"`
	KeyID     string    `json:"key_id"`
}

type Musigs []Musigs

type SignResult struct {
	Musig Musig
	Err   error
}

func (m *Musig) MarshalJSON() ([]byte, error) {
	return json.Marshal(*m)
}

func (r Musigs) MarshalJSON() ([]byte, error) {
	return json.Marshal(r)
}

func (m *Manager) Sign(ctx context.Context, msg []byte) ([]Musig, error) {
	var (
		mu          sync.Mutex
		results     []Musig
		totalWeight int
		wg          sync.WaitGroup
		done        = make(chan struct{})
	)

	resultCh := make(chan SignResult, len(m.clients))

	for _, client := range m.clients {
		wg.Add(1)
		go func(client *Client) {
			defer wg.Done()

			resp, err := client.Sign(ctx, msg)
			select {
			case resultCh <- SignResult{
				Musig: Musig{
					Sig:       resp.Sig,
					Weight:    client.Weight,
					Threshold: m.Threshold,
					PartyID:   client.PartyID,
					KeyID:     client.KeyID,
				},
				Err: err,
			}:
			case <-done:
			}
		}(client)
	}

	for {
		res := <-resultCh
		if res.Err == nil && len(res.Musig.Sig) >= 64 {
			mu.Lock()
			results = append(results, res.Musig)
			totalWeight += res.Musig.Weight
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

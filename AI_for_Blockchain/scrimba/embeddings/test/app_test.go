package main_test

import (
	"testing"

	"github.com/peartes/scrimba/embeddings/app"
	"github.com/peartes/scrimba/embeddings/config"
	"github.com/stretchr/testify/require"
	"github.com/supabase-community/supabase-go"
)

func TestApp(t *testing.T) {
	t.Parallel()
	supabaseClient, err := supabase.NewClient(config.GetSupaBaseUrl(), config.GetSupaBaseAPIKey(), nil)
	if err != nil {
		t.Errorf("supabase.NewClient() error = %v, want nil", err)
		return
	}

	t.Run("saves embedding into superbase", func(t *testing.T) {
		t.Parallel()

		err := app.RunApp()

		if err != nil {
			t.Errorf("RunApp() error = %v, want nil", err)
			return
		}

		// confirm that the embeddings were saved
		_, count, err := supabaseClient.From("documents").Select("id", "exact", false).Execute()
		if err != nil {
			t.Errorf("supabaseClient.From().Select() error = %v, want nil", err)
			return
		}
		require.Equal(t, len(config.GetPodcastMockup()), count)
	})
}

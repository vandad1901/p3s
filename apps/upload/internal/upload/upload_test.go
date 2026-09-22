package upload_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vandad1901/p3s/apps/upload/internal/app"
	"github.com/vandad1901/p3s/apps/upload/internal/config"
	"github.com/vandad1901/p3s/apps/upload/internal/upload"
	"github.com/vandad1901/p3s/packages/go/usercontext"
)

// purposefully passing nil to test whether validations short-circuit.
func nilService() *upload.Service {
	return upload.NewService(nil, nil, nil)
}

func TestGenerateURL_RequiresAuthenticatedUser(t *testing.T) {
	t.Parallel()

	svc := nilService()
	ctx := context.Background()

	_, _, err := svc.GenerateURL(ctx, "valid-key")
	require.ErrorIs(t, err, usercontext.ErrNotFound)
}

func TestGenerateURL_RejectsInvalidKeys(t *testing.T) {
	t.Parallel()

	svc := nilService()
	ctx := usercontext.CtxWithUser(context.Background(), 1)

	tests := []struct {
		name string
		key  string
	}{
		{"empty", ""},
		{"slash", "foo/bar"},
		{"path traversal", "../etc/passwd"},
		{"dot", "image.png"},
		{"space", "my image"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, _, err := svc.GenerateURL(ctx, tt.key)
			require.Error(t, err)
		})
	}
}

func TestGenerateURL_AcceptsValidKeys(t *testing.T) {
	t.Parallel()

	svc := app.MustBoot(config.LoadConfig()).UploadService
	ctx := usercontext.CtxWithUser(context.Background(), 1)

	tests := []string{"myimage", "MYIMAGE", "my-image_final", "aaabbb123", "a_b-c_1-2-3"}

	for _, key := range tests {
		t.Run(key, func(t *testing.T) {
			t.Parallel()

			_, _, err := svc.GenerateURL(ctx, key)
			require.NoError(t, err)
		})
	}
}

func TestGenerateURL_ScopesKeyToUser(t *testing.T) {
	t.Parallel()

	svc := app.MustBoot(config.LoadConfig()).UploadService
	ctx := usercontext.CtxWithUser(context.Background(), 42)

	_, formFields, err := svc.GenerateURL(ctx, "valid-key")
	if err != nil {
		t.Fatalf("GenerateURL() error = %v", err)
	}

	const wantKey = "42/valid-key"
	require.Equal(t, wantKey, formFields["key"])
}

func TestFinalizeUpload_RequiresAuthenticatedUser(t *testing.T) {
	t.Parallel()

	svc := nilService()
	ctx := context.Background()

	err := svc.FinalizeUpload(ctx, "valid-key")
	require.ErrorIs(t, err, usercontext.ErrNotFound)
}

func TestFinalizeUpload_RejectsInvalidKey(t *testing.T) {
	t.Parallel()

	svc := nilService()
	ctx := usercontext.CtxWithUser(context.Background(), 1)

	err := svc.FinalizeUpload(ctx, "has spaces")
	require.Error(t, err)
}

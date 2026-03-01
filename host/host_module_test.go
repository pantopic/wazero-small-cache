package wazero_small_cache

import (
	"context"
	_ "embed"
	"os"
	"testing"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

//go:embed test\.wasm
var testwasm []byte

func TestModule(t *testing.T) {
	var (
		ctx = context.Background()
	)
	r := wazero.NewRuntimeWithConfig(ctx, wazero.NewRuntimeConfig())
	wasi_snapshot_preview1.MustInstantiate(ctx, r)

	hostModule := New()
	hostModule.Register(ctx, r)

	compiled, err := r.CompileModule(ctx, testwasm)
	if err != nil {
		panic(err)
	}
	cfg := wazero.NewModuleConfig().WithStdout(os.Stdout).WithName(`a`)
	mod1, err := r.InstantiateModule(ctx, compiled, cfg)
	if err != nil {
		t.Errorf(`%v`, err)
		return
	}
	cfg = wazero.NewModuleConfig().WithStdout(os.Stdout).WithName(`b`)
	mod2, err := r.InstantiateModule(ctx, compiled, cfg)
	if err != nil {
		t.Errorf(`%v`, err)
		return
	}
	ctx, err = hostModule.InitContext(ctx, mod1)
	if err != nil {
		t.Fatalf(`%v`, err)
	}

	ctx = hostModule.ContextCopy(ctx, ctx)

	t.Run(`put`, func(t *testing.T) {
		_, err := mod1.ExportedFunction(`testLocalPut`).Call(ctx, uint64(1), uint64(2))
		if err != nil {
			t.Fatalf("%v", err)
		}
	})
	t.Run(`get`, func(t *testing.T) {
		stack, err := mod2.ExportedFunction(`testLocalGet`).Call(ctx, uint64(1))
		if err != nil {
			t.Fatalf("%v", err)
		}
		if stack[0] != uint64(2) {
			t.Fatalf("expected %d, got %d", 2, stack[0])
		}
	})
	t.Run(`put2`, func(t *testing.T) {
		_, err := mod1.ExportedFunction(`testLocalPut2`).Call(ctx, uint64(1), uint64(3))
		if err != nil {
			t.Fatalf("%v", err)
		}
	})
	t.Run(`get2`, func(t *testing.T) {
		stack, err := mod2.ExportedFunction(`testLocalGet2`).Call(ctx, uint64(1))
		if err != nil {
			t.Fatalf("%v", err)
		}
		if stack[0] != uint64(3) {
			t.Fatalf("expected %d, got %d", 3, stack[0])
		}
	})
	t.Run(`del`, func(t *testing.T) {
		_, err := mod2.ExportedFunction(`testLocalDel`).Call(ctx, uint64(1))
		if err != nil {
			t.Fatalf("%v", err)
		}
		stack, err := mod2.ExportedFunction(`testLocalGet`).Call(ctx, uint64(1))
		if err != nil {
			t.Fatalf("%v", err)
		}
		if stack[0] != uint64(0) {
			t.Fatalf("expected %d, got %d", 0, stack[0])
		}
	})

	hostModule.Stop()
}

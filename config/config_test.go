package config

import (
	"bytes"
	"testing"
	"time"

	"github.com/valyala/fasthttp"
	"github.com/zishang520/engine.io-server-go-fasthttp/v2/types"
)

func TestAttachOptionsDefauleValue(t *testing.T) {
	opts := AttachOptionsInterface(DefaultAttachOptions())
	t.Run("path", func(t *testing.T) {
		if path := opts.Path(); opts.GetRawPath() == nil && path != "/engine.io" {
			t.Fatalf(`*AttachOptions.Path() = %q, want match for %#q`, path, "/engine.io")
		}
	})

	t.Run("destroyUpgrade", func(t *testing.T) {
		if destroyUpgrade := opts.DestroyUpgrade(); opts.GetRawDestroyUpgrade() == nil && destroyUpgrade != true {
			t.Fatalf(`*AttachOptions.DestroyUpgrade() = %t, want match for %t`, destroyUpgrade, true)
		}
	})

	t.Run("destroyUpgradeTimeout", func(t *testing.T) {
		if destroyUpgradeTimeout := opts.DestroyUpgradeTimeout(); opts.GetRawDestroyUpgradeTimeout() == nil && destroyUpgradeTimeout != 1000*time.Millisecond {
			t.Fatalf(`*AttachOptions.DestroyUpgradeTimeout() = %d, want match for %d`, destroyUpgradeTimeout, 1000*time.Millisecond)
		}
	})

	t.Run("addTrailingSlash", func(t *testing.T) {
		if addTrailingSlash := opts.AddTrailingSlash(); opts.GetRawAddTrailingSlash() == nil && addTrailingSlash != true {
			t.Fatalf(`*AttachOptions.AddTrailingSlash() = %t, want match for %t`, addTrailingSlash, true)
		}
	})
}

func TestAttachOptionsSetValue(t *testing.T) {
	opts := AttachOptionsInterface(DefaultAttachOptions())
	t.Run("path", func(t *testing.T) {
		opts.SetPath("test")
		if path := opts.Path(); path != "test" {
			t.Fatalf(`*AttachOptions.Path() = %q, want match for %#q`, path, "test")
		}
	})

	t.Run("destroyUpgrade", func(t *testing.T) {
		opts.SetDestroyUpgrade(false)
		if destroyUpgrade := opts.DestroyUpgrade(); destroyUpgrade != false {
			t.Fatalf(`*AttachOptions.DestroyUpgrade() = %t, want match for %t`, destroyUpgrade, false)
		}
	})

	t.Run("destroyUpgradeTimeout", func(t *testing.T) {
		opts.SetDestroyUpgradeTimeout(5000 * time.Millisecond)
		if destroyUpgradeTimeout := opts.DestroyUpgradeTimeout(); destroyUpgradeTimeout != 5000*time.Millisecond {
			t.Fatalf(`*AttachOptions.DestroyUpgradeTimeout() = %d, want match for %d`, destroyUpgradeTimeout, 5000*time.Millisecond)
		}
	})
}

func TestServerOptionsDefauleValue(t *testing.T) {
	opts := ServerOptionsInterface(DefaultServerOptions())

	t.Run("pingTimeout", func(t *testing.T) {
		if pingTimeout := opts.PingTimeout(); opts.GetRawPingTimeout() == nil && pingTimeout != 20000*time.Millisecond {
			t.Fatalf(`*ServerOptions.PingTimeout() = %d, want match for %d`, pingTimeout, 20000*time.Millisecond)
		}
	})

	t.Run("pingInterval", func(t *testing.T) {
		if pingInterval := opts.PingInterval(); opts.GetRawPingInterval() == nil && pingInterval != 25000*time.Millisecond {
			t.Fatalf(`*ServerOptions.PingInterval() = %d, want match for %d`, pingInterval, 25000*time.Millisecond)
		}
	})

	t.Run("upgradeTimeout", func(t *testing.T) {
		if upgradeTimeout := opts.UpgradeTimeout(); opts.GetRawUpgradeTimeout() == nil && upgradeTimeout != 10000*time.Millisecond {
			t.Fatalf(`*ServerOptions.UpgradeTimeout() = %d, want match for %d`, upgradeTimeout, 10000*time.Millisecond)
		}
	})

	t.Run("maxHttpBufferSize", func(t *testing.T) {
		if maxHttpBufferSize := opts.MaxHttpBufferSize(); opts.GetRawMaxHttpBufferSize() == nil && maxHttpBufferSize != 100000 {
			t.Fatalf(`*ServerOptions.MaxHttpBufferSize() = %d, want match for %d`, maxHttpBufferSize, 100000)
		}
	})

	t.Run("allowRequest", func(t *testing.T) {
		if allowRequest := opts.AllowRequest(); opts.GetRawAllowRequest() == nil && allowRequest != nil {
			t.Fatalf(`*ServerOptions.AllowRequest() = %v, want match for nil`, allowRequest)
		}
	})

	t.Run("transports", func(t *testing.T) {
		if transports := opts.Transports(); opts.GetRawTransports() == nil && transports != nil && !(transports.Has("polling") && transports.Has("websocket")) {
			t.Fatalf(`*ServerOptions.Transports() = %s, want match for ["polling", "websocket")]`, transports.Keys())
		}
	})

	t.Run("allowUpgrades", func(t *testing.T) {
		if allowUpgrades := opts.AllowUpgrades(); opts.GetRawAllowUpgrades() == nil && allowUpgrades != true {
			t.Fatalf(`*ServerOptions.AllowUpgrades() = %t, want match for %t`, allowUpgrades, true)
		}
	})

	t.Run("perMessageDeflate", func(t *testing.T) {
		if perMessageDeflate := opts.PerMessageDeflate(); opts.GetRawPerMessageDeflate() == nil && perMessageDeflate != nil {
			t.Fatalf(`*ServerOptions.PerMessageDeflate() = %v, want match for nil`, perMessageDeflate)
		}
	})

	t.Run("httpCompression/threshold", func(t *testing.T) {
		if httpCompression := opts.HttpCompression(); opts.GetRawHttpCompression() == nil && httpCompression != nil && httpCompression.Threshold != 1024 {
			t.Fatalf(`*ServerOptions.HttpCompression().Threshold = %d, want match for %d`, httpCompression.Threshold, 1024)
		}
	})

	t.Run("initialPacket", func(t *testing.T) {
		if initialPacket := opts.InitialPacket(); opts.GetRawInitialPacket() == nil && initialPacket != nil {
			t.Fatalf(`*ServerOptions.InitialPacket() = %v, want match for nil`, initialPacket)
		}
	})

	t.Run("cookie", func(t *testing.T) {
		if cookie := opts.Cookie(); opts.GetRawCookie() == nil && cookie != nil {
			t.Fatalf(`*ServerOptions.Cookie() = %v, want match for nil`, cookie)
		}
	})

	t.Run("cors", func(t *testing.T) {
		if cors := opts.Cors(); opts.GetRawCors() == nil && cors != nil {
			t.Fatalf(`*ServerOptions.Cors() = %v, want match for nil`, cors)
		}
	})

	t.Run("allowEIO3", func(t *testing.T) {
		if allowEIO3 := opts.AllowEIO3(); opts.GetRawAllowEIO3() == nil && allowEIO3 != false {
			t.Fatalf(`*ServerOptions.AllowEIO3() = %t, want match for %t`, allowEIO3, false)
		}
	})
}

func TestServerOptionsSetValue(t *testing.T) {
	opts := ServerOptionsInterface(DefaultServerOptions())

	t.Run("pingTimeout", func(t *testing.T) {
		opts.SetPingTimeout(10 * time.Millisecond)
		if pingTimeout := opts.PingTimeout(); pingTimeout != 10*time.Millisecond {
			t.Fatalf(`*ServerOptions.PingTimeout() = %d, want match for %d`, pingTimeout, 10*time.Millisecond)
		}
	})

	t.Run("pingInterval", func(t *testing.T) {
		opts.SetPingInterval(15 * time.Millisecond)
		if pingInterval := opts.PingInterval(); pingInterval != 15*time.Millisecond {
			t.Fatalf(`*ServerOptions.PingInterval() = %d, want match for %d`, pingInterval, 15*time.Millisecond)
		}
	})

	t.Run("upgradeTimeout", func(t *testing.T) {
		opts.SetUpgradeTimeout(10000 * time.Millisecond)
		if upgradeTimeout := opts.UpgradeTimeout(); upgradeTimeout != 10000*time.Millisecond {
			t.Fatalf(`*ServerOptions.UpgradeTimeout() = %d, want match for %d`, upgradeTimeout, 10000*time.Millisecond)
		}
	})

	t.Run("maxHttpBufferSize", func(t *testing.T) {
		opts.SetMaxHttpBufferSize(999)
		if maxHttpBufferSize := opts.MaxHttpBufferSize(); maxHttpBufferSize != 999 {
			t.Fatalf(`*ServerOptions.MaxHttpBufferSize() = %d, want match for %d`, maxHttpBufferSize, 999)
		}
	})

	t.Run("allowRequest", func(t *testing.T) {
		opts.SetAllowRequest(nil)
		if allowRequest := opts.AllowRequest(); allowRequest != nil {
			t.Fatalf(`*ServerOptions.AllowRequest() = %v, want match for nil`, allowRequest)
		}
	})

	t.Run("transports", func(t *testing.T) {
		opts.SetTransports(types.NewSet("websocket", "polling"))
		if transports := opts.Transports(); transports != nil && !(transports.Has("polling") && transports.Has("websocket")) {
			t.Fatalf(`*ServerOptions.Transports() = %s, want match for ["polling", "websocket")]`, transports.Keys())
		}
	})

	t.Run("allowUpgrades", func(t *testing.T) {
		opts.SetAllowUpgrades(false)
		if allowUpgrades := opts.AllowUpgrades(); allowUpgrades != false {
			t.Fatalf(`*ServerOptions.AllowUpgrades() = %t, want match for %t`, allowUpgrades, false)
		}
	})

	t.Run("perMessageDeflate", func(t *testing.T) {
		input := &types.PerMessageDeflate{Threshold: 1024}
		opts.SetPerMessageDeflate(input)
		if perMessageDeflate := opts.PerMessageDeflate(); perMessageDeflate.Threshold != 1024 {
			t.Fatalf(`*ServerOptions.PerMessageDeflate().Threshold = %d, want match for %d`, perMessageDeflate.Threshold, 1024)
		}
	})

	t.Run("httpCompression/threshold", func(t *testing.T) {
		input := &types.HttpCompression{Threshold: 2048}
		opts.SetHttpCompression(input)
		if httpCompression := opts.HttpCompression(); httpCompression != nil && httpCompression.Threshold != 2048 {
			t.Fatalf(`*ServerOptions.HttpCompression().Threshold = %d, want match for %d`, httpCompression.Threshold, 2048)
		}
	})

	t.Run("initialPacket", func(t *testing.T) {
		input := bytes.NewBuffer([]byte{1})
		opts.SetInitialPacket(input)
		if initialPacket := opts.InitialPacket(); initialPacket != input {
			t.Fatalf(`*ServerOptions.InitialPacket() = %v, want match for %v`, initialPacket, input)
		}
	})

	t.Run("cookie", func(t *testing.T) {
		input := &fasthttp.Cookie{}
		input.SetKey("Test")
		opts.SetCookie(input)
		if cookie := opts.Cookie(); cookie != input {
			t.Fatalf(`*ServerOptions.Cookie() = %v, want match for %v`, cookie, input)
		}
	})

	t.Run("cors", func(t *testing.T) {
		input := &types.Cors{
			Origin: "http://localhost",
		}
		opts.SetCors(input)
		if cors := opts.Cors(); cors != input {
			t.Fatalf(`*ServerOptions.Cors() = %v, want match for %v`, cors, input)
		}
	})

	t.Run("allowEIO3", func(t *testing.T) {
		opts.SetAllowEIO3(true)
		if allowEIO3 := opts.AllowEIO3(); allowEIO3 != true {
			t.Fatalf(`*ServerOptions.AllowEIO3() = %t, want match for %t`, allowEIO3, true)
		}
	})
}

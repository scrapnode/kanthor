package signature_test

import (
	"testing"

	"github.com/scrapnode/kanthor/pkg/signature"
	"github.com/stretchr/testify/assert"
)

func TestSign(t *testing.T) {
	t.Run("external", func(st *testing.T) {
		secret := "20c9cf0526f6426d913c11e49b6571ab"

		assert.Equal(st, "6dd54bfedfedcaf5a805e3bf298b9bc8aed2b8c90fb3101d467aec9cb0b2f04c", signature.Sign(secret, `msg_2aWPa1QEfbS80peVx0LZmf07fqr.1704431421865.{"source":"console"}`))
		assert.Equal(st, "35692ce9dfd82fdf1b974005ec5c85d494a888c464c40e3cd571ab61ce8ca9ad", signature.Sign(secret, `msg_2aWQksgy08hUmVNF93p8mu8MW7n.1704432001472.{"source":"console"}`))
	})
}

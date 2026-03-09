package signer

import (
	"context"
	"encoding/hex"
	"net/http"
	"strings"

	"github.com/longbridge/openapi-go/internal/util"

	"github.com/pkg/errors"
)

const AlgHmacSHA256 = "HMAC-SHA256"
const AlgHmacSHA1 = "HMAC-SHA1"
const AlgHmacMD5 = "HMAC-MD5"

var ErrInvalidSignatureText = errors.New("invalid signature text")
var ErrInvalidHashAlg = errors.New("invalid hash algorithm")

var supportSigner = []string{AlgHmacSHA256, AlgHmacSHA1, AlgHmacMD5}

func isSupportedSigner(s string) bool {
	for _, item := range supportSigner {
		if s == item {
			return true
		}
	}

	return false
}

// Signer is for to sign request can access longbridge api
type Signer struct{}

func (_ *Signer) String() string {
	return "openapi"
}

func (s *Signer) parseSignature(str string) (signer string, signerHeaders []string, expect string, err error) {
	idx := strings.Index(str, " ")

	if idx == -1 {
		err = ErrInvalidSignatureText
		return
	}

	signer = str[:idx]

	if !isSupportedSigner(signer) {
		err = errors.Wrapf(ErrInvalidSignatureText, "unsupported sign method")
		return
	}

	str = str[idx+1:]

	parts := strings.Split(str, ",")

	kvs := make([]string, 0, len(parts)*2)

	for _, item := range parts {
		kv := strings.SplitN(item, "=", 2)

		if len(kv) != 2 {
			err = ErrInvalidSignatureText
			return
		}

		kvs = append(kvs, strings.TrimSpace(kv[0]), strings.TrimSpace(kv[1]))
	}

	for i := 0; i < len(kvs); i = i + 2 {
		k := kvs[i]
		v := kvs[i+1]

		switch k {
		case "SignedHeaders":
			signerHeaders = strings.Split(v, ";")
		case "Signature":
			expect = v
		}

		if len(signerHeaders) > 0 && expect != "" {
			break
		}
	}

	if len(signerHeaders) == 0 {
		err = errors.Wrapf(ErrInvalidSignatureText, "need SignedHeaders")
		return
	}

	return
}

// Sign signature method of Longbridge OpenApi
func (s *Signer) Sign(ctx context.Context, secret []byte, req *http.Request, body []byte) (signature string, plain string, equal bool, err error) {
	var (
		signedHeaders  []string
		signer, expect string
		encrypted      []byte
	)

	if signer, signedHeaders, expect, err = s.parseSignature(req.Header.Get("x-api-signature")); err != nil {
		return
	}

	plain = s.plainText(ctx, req, signedHeaders, body)

	textToSign := signer + "|" + hex.EncodeToString(util.SHA1(util.UnsafeStringToBytes(plain)))

	switch signer {
	case AlgHmacSHA1:
		encrypted, err = util.HmacSHA1(secret, util.UnsafeStringToBytes(textToSign))
	case AlgHmacSHA256:
		encrypted, err = util.HmacSHA256(secret, util.UnsafeStringToBytes(textToSign))
	case AlgHmacMD5:
		encrypted, err = util.HmacMD5(secret, util.UnsafeStringToBytes(textToSign))
	default:
		err = ErrInvalidHashAlg
	}

	if err != nil {
		return
	}

	signature = hex.EncodeToString(encrypted)
	equal = signature == expect

	return
}

func (s *Signer) plainText(ctx context.Context, req *http.Request, signedHeaders []string, body []byte) string {
	mtd := strings.ToUpper(req.Method)
	uri := req.URL.Path
	query := req.URL.RawQuery

	plain := mtd + "|" +
		uri + "|" +
		query

	headers := ""

	for _, item := range signedHeaders {
		key := strings.ToLower(strings.TrimSpace(item))
		headers = headers + key + ":" + strings.TrimSpace(req.Header.Get(key)) + "\n"
	}

	plain = plain + "|" +
		headers + "|" +
		strings.Join(signedHeaders, ";") + "|"

	if len(body) != 0 {
		hashBody := util.SHA1(body)
		plain = plain + hex.EncodeToString(hashBody)
	}

	return plain
}

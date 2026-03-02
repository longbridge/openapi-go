package http

import (
	"context"
	nhttp "net/http"
	"strconv"
	"strings"
	"time"

	"github.com/longportapp/openapi-go/internal/signer"
	"github.com/longportapp/openapi-go/internal/util"
)

const (
	headerTimestamp = "x-timestamp"
)

var sign = &signer.Signer{}

func signature(req *nhttp.Request, secret string, body []byte) error {
	// OAuth 2.0 mode: skip HMAC signature generation when using Bearer token
	authHeader := req.Header.Get("authorization")
	if strings.HasPrefix(authHeader, "Bearer ") {
		return nil
	}

	if v := req.Header.Get(headerTimestamp); v == "" {
		req.Header.Add(headerTimestamp, strconv.FormatInt(time.Now().UnixMilli(), 10))
	}

	req.Header.Add("x-api-signature", "HMAC-SHA256 SignedHeaders=authorization;x-api-key;x-timestamp")
	signstr, _, _, err := sign.Sign(context.Background(), util.UnsafeStringToBytes(secret), req, body)
	if err != nil {
		return err
	}
	req.Header.Set("x-api-signature", "HMAC-SHA256 SignedHeaders=authorization;x-api-key;x-timestamp, Signature="+signstr)
	return nil
}

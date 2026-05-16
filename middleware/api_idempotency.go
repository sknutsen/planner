package middleware

import (
	"bytes"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/sknutsen/planner/database"
	"github.com/sknutsen/planner/internal/apijson"
)

// IdempotencyBuffer hashes Idempotency-Key + request for duplicate detection, restores the body for handlers,
// and caches successful JSON responses for replay.
func IdempotencyBuffer(db database.DBTX) echo.MiddlewareFunc {
	q := database.New(db)
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			method := c.Request().Method
			if method == http.MethodGet || method == http.MethodHead || method == http.MethodOptions {
				return next(c)
			}

			key := strings.TrimSpace(c.Request().Header.Get("Idempotency-Key"))
			if key == "" {
				return next(c)
			}

			body, err := io.ReadAll(c.Request().Body)
			if err != nil {
				return apijson.Error(c, http.StatusBadRequest, "BAD_REQUEST", "Could not read request body for idempotency.")
			}
			c.Request().Body = io.NopCloser(bytes.NewReader(body))

			userID, _ := c.Get(CtxAPIUserID).(string)
			if userID == "" {
				return next(c)
			}

			sumKey := sha256.Sum256([]byte(key))
			keyHash := hex.EncodeToString(sumKey[:])
			hr := method + "\n" + c.Request().URL.Path + "\n"
			sumReq := sha256.Sum256(append([]byte(hr), body...))
			requestHash := hex.EncodeToString(sumReq[:])

			row, err := q.GetIdempotencyRecord(c.Request().Context(), database.GetIdempotencyRecordParams{
				UserID:  userID,
				KeyHash: keyHash,
			})
			if err == nil {
				if row.RequestHash != requestHash {
					return apijson.Error(c, http.StatusConflict, "IDEMPOTENCY_KEY_REUSE", "Same Idempotency-Key was used with a different request fingerprint.")
				}
				c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
				c.Response().WriteHeader(int(row.StatusCode))
				_, werr := c.Response().Write([]byte(row.ResponseBody))
				return werr
			}
			if !errors.Is(err, sql.ErrNoRows) {
				return apijson.Error(c, http.StatusInternalServerError, "SERVER_ERROR", "Idempotency lookup failed.")
			}

			rw := &captureWriter{ResponseWriter: c.Response().Writer}
			c.Response().Writer = rw

			nerr := next(c)
			if nerr != nil {
				return nerr
			}

			st := rw.status
			if st == 0 {
				st = http.StatusOK
			}
			if st < http.StatusOK || st >= http.StatusMultipleChoices {
				return nil
			}

			insErr := q.InsertIdempotencyResponse(c.Request().Context(), database.InsertIdempotencyResponseParams{
				UserID:       userID,
				KeyHash:      keyHash,
				RequestHash:  requestHash,
				ResponseBody: string(rw.buf),
				StatusCode:   int64(st),
			})
			if insErr != nil && strings.Contains(insErr.Error(), "UNIQUE") {
				row2, e2 := q.GetIdempotencyRecord(c.Request().Context(), database.GetIdempotencyRecordParams{
					UserID:  userID,
					KeyHash: keyHash,
				})
				if e2 == nil && row2.RequestHash == requestHash {
					// Another request completed the same idempotent call; response already sent to client.
					return nil
				}
			}
			return nil
		}
	}
}

type captureWriter struct {
	http.ResponseWriter
	status int
	buf    []byte
}

func (w *captureWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *captureWriter) Write(b []byte) (int, error) {
	if w.status == 0 {
		w.status = http.StatusOK
	}
	w.buf = append(w.buf, b...)
	return w.ResponseWriter.Write(b)
}

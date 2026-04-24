package grpcutil

import (
	"context"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
)

// ---- error_test ----

func TestIsCode(t *testing.T) {
	err := ErrNotFound("user not found")
	if !IsCode(err, codes.NotFound) {
		t.Fatal("expected NotFound")
	}
	if IsCode(err, codes.Internal) {
		t.Fatal("unexpected Internal")
	}
}

func TestIsNotFound(t *testing.T) {
	if !IsNotFound(ErrNotFound("x")) {
		t.Fatal("expected true")
	}
	if IsNotFound(ErrInternal("x")) {
		t.Fatal("expected false")
	}
}

func TestFromError_nil(t *testing.T) {
	code, msg := FromError(nil)
	if code != codes.OK || msg != "" {
		t.Fatalf("expected OK/empty, got %v/%q", code, msg)
	}
}

func TestCodeToHTTPStatus(t *testing.T) {
	cases := []struct {
		code codes.Code
		want int
	}{
		{codes.OK, 200},
		{codes.NotFound, 404},
		{codes.InvalidArgument, 400},
		{codes.Unauthenticated, 401},
		{codes.PermissionDenied, 403},
		{codes.Internal, 500},
		{codes.Unavailable, 503},
		{codes.DeadlineExceeded, 504},
	}
	for _, tc := range cases {
		got := CodeToHTTPStatus(tc.code)
		if got != tc.want {
			t.Errorf("CodeToHTTPStatus(%v) = %d, want %d", tc.code, got, tc.want)
		}
	}
}

func TestHTTPStatusToCode(t *testing.T) {
	if HTTPStatusToCode(404) != codes.NotFound {
		t.Error("expected NotFound")
	}
	if HTTPStatusToCode(200) != codes.OK {
		t.Error("expected OK")
	}
}

// ---- metadata_test ----

func TestMetaGet_missing(t *testing.T) {
	ctx := context.Background()
	if got := MetaGet(ctx, "x-foo"); got != "" {
		t.Fatalf("expected empty, got %q", got)
	}
}

func TestMetaGet_present(t *testing.T) {
	md := metadata.Pairs("x-request-id", "abc123")
	ctx := metadata.NewIncomingContext(context.Background(), md)
	if got := MetaGet(ctx, "x-request-id"); got != "abc123" {
		t.Fatalf("expected abc123, got %q", got)
	}
}

func TestExtractBearerToken(t *testing.T) {
	md := metadata.Pairs("authorization", "Bearer my-token-xyz")
	ctx := metadata.NewIncomingContext(context.Background(), md)
	if got := ExtractBearerToken(ctx); got != "my-token-xyz" {
		t.Fatalf("expected my-token-xyz, got %q", got)
	}
}

func TestExtractBearerToken_empty(t *testing.T) {
	ctx := context.Background()
	if got := ExtractBearerToken(ctx); got != "" {
		t.Fatalf("expected empty, got %q", got)
	}
}

// ---- pagination_test ----

func TestPageRequest_Offset(t *testing.T) {
	r := PageRequest{Page: 3, PageSize: 10}
	if got := r.Offset(); got != 20 {
		t.Fatalf("expected 20, got %d", got)
	}
}

func TestPageRequest_Offset_firstPage(t *testing.T) {
	r := PageRequest{Page: 1, PageSize: 10}
	if got := r.Offset(); got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
}

func TestNewPageResponse(t *testing.T) {
	pr := NewPageResponse(1, 10, 25)
	if pr.TotalPages != 3 {
		t.Fatalf("expected 3 pages, got %d", pr.TotalPages)
	}
	if !pr.HasMore {
		t.Fatal("expected HasMore=true")
	}
}

func TestNewPageResponse_lastPage(t *testing.T) {
	pr := NewPageResponse(3, 10, 25)
	if pr.HasMore {
		t.Fatal("expected HasMore=false on last page")
	}
}

func TestCursor_roundtrip(t *testing.T) {
	type payload struct {
		ID   string `json:"id"`
		Page int    `json:"page"`
	}
	original := payload{ID: "user-123", Page: 5}

	token, err := EncodeCursor(original)
	if err != nil {
		t.Fatalf("encode: %v", err)
	}

	var decoded payload
	if err := DecodeCursor(token, &decoded); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if decoded.ID != original.ID || decoded.Page != original.Page {
		t.Fatalf("roundtrip mismatch: %+v", decoded)
	}
}

func TestEncodeOffsetCursor(t *testing.T) {
	token, err := EncodeOffsetCursor(50)
	if err != nil {
		t.Fatalf("encode: %v", err)
	}
	offset, err := DecodeOffsetCursor(token)
	if err != nil {
		t.Fatalf("decode: %v", err)
	}
	if offset != 50 {
		t.Fatalf("expected 50, got %d", offset)
	}
}

func TestDecodeOffsetCursor_empty(t *testing.T) {
	offset, err := DecodeOffsetCursor("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if offset != 0 {
		t.Fatalf("expected 0, got %d", offset)
	}
}

// ---- validate_test ----

func TestValidationError_noErrors(t *testing.T) {
	err := Validate().Required("name", "John").Err()
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestValidationError_required(t *testing.T) {
	err := Validate().Required("name", "").Err()
	if err == nil {
		t.Fatal("expected error for empty required field")
	}
	if !IsInvalidArg(err) {
		t.Fatal("expected InvalidArgument code")
	}
}

func TestValidationError_multiple(t *testing.T) {
	err := Validate().
		Required("name", "").
		Positive("age", -1).
		MaxLen("bio", "this is way too long for our limit", 5).
		Err()
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestValidationError_positive(t *testing.T) {
	if Validate().Positive("count", 1).HasErrors() {
		t.Fatal("positive int should pass")
	}
	if !Validate().Positive("count", 0).HasErrors() {
		t.Fatal("zero should fail")
	}
}

func TestIsRetryable(t *testing.T) {
	if !IsRetryable(ErrUnavailable("down")) {
		t.Fatal("unavailable should be retryable")
	}
	if IsRetryable(ErrNotFound("x")) {
		t.Fatal("not found should not be retryable")
	}
}

func TestContextError_canceled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err := ContextError(ctx)
	if !IsCanceled(err) {
		t.Fatalf("expected Canceled, got %v", err)
	}
}

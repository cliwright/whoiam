package internal

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/sts"
)

func strPtr(s string) *string { return &s }

// --- AssertAccountAsExpected ---

func TestAssertAccountAsExpected_Match(t *testing.T) {
	identity := &sts.GetCallerIdentityOutput{Account: strPtr("123456789012")}
	if err := AssertAccountAsExpected(identity, "123456789012"); err != nil {
		t.Errorf("expected nil, got %v", err)
	}
}

func TestAssertAccountAsExpected_Mismatch(t *testing.T) {
	identity := &sts.GetCallerIdentityOutput{Account: strPtr("999999999999")}
	err := AssertAccountAsExpected(identity, "123456789012")
	if err == nil {
		t.Fatal("expected error on account mismatch, got nil")
	}
	if !strings.Contains(err.Error(), "123456789012") || !strings.Contains(err.Error(), "999999999999") {
		t.Errorf("error message should contain both expected and actual account numbers, got: %s", err.Error())
	}
}

func TestAssertAccountAsExpected_NilAccount(t *testing.T) {
	identity := &sts.GetCallerIdentityOutput{Account: nil}
	if err := AssertAccountAsExpected(identity, "123456789012"); err == nil {
		t.Error("expected error for nil account field, got nil")
	}
}

// --- PrintCallerIdentityTable ---

func TestPrintCallerIdentityTable_AllFields(t *testing.T) {
	identity := &sts.GetCallerIdentityOutput{
		Account: strPtr("123456789012"),
		Arn:     strPtr("arn:aws:iam::123456789012:user/test"),
		UserId:  strPtr("AIDAEXAMPLE"),
	}

	// Capture stdout to verify output without panicking
	orig := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	PrintCallerIdentityTable(identity, "production")

	w.Close()
	os.Stdout = orig

	var buf bytes.Buffer
	io.Copy(&buf, r)
	out := buf.String()

	if !strings.Contains(out, "production") {
		t.Errorf("expected output to contain account name, got: %s", out)
	}
	if !strings.Contains(out, "123456789012") {
		t.Errorf("expected output to contain account number, got: %s", out)
	}
}

func TestPrintCallerIdentityTable_NilFields(t *testing.T) {
	identity := &sts.GetCallerIdentityOutput{Account: nil, Arn: nil, UserId: nil}

	// Should not panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("PrintCallerIdentityTable panicked on nil fields: %v", r)
		}
	}()

	orig := os.Stdout
	_, w, _ := os.Pipe()
	os.Stdout = w
	PrintCallerIdentityTable(identity, "test")
	w.Close()
	os.Stdout = orig
}

package patterns

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

/**
 * ===========================================================================
 * GoF Design Patterns Tests
 * ===========================================================================
 *
 * For a Java developer:
 * - Go's patterns are often simpler because they rely on implicit interfaces.
 * - Singleton is implemented with `sync.Once`.
 * - Iterator is often implemented with a closure that returns values.
 * - Adapter is just a struct embedding or a thin wrapper.
 */

// ---------------------------------------------------------------------------
// Singleton
// ---------------------------------------------------------------------------

func TestSingletonIdentity(t *testing.T) {
	a := GetAppConfig()
	b := GetAppConfig()
	if a != b {
		t.Error("GetAppConfig should return the same pointer every time")
	}
}

// ---------------------------------------------------------------------------
// Factory Method
// ---------------------------------------------------------------------------

func TestNewNotifier(t *testing.T) {
	cases := []struct {
		nType, dest, wantPrefix string
	}{
		{"email", "a@b.com", "[Email"},
		{"sms", "+1", "[SMS"},
		{"slack", "dev", "[Slack"},
	}
	for _, tc := range cases {
		n := NewNotifier(tc.nType, tc.dest)
		out := n.Send("hi")
		if !strings.HasPrefix(out, tc.wantPrefix) {
			t.Errorf("NewNotifier(%q).Send() = %q, want prefix %q", tc.nType, out, tc.wantPrefix)
		}
	}
}

// ---------------------------------------------------------------------------
// Builder
// ---------------------------------------------------------------------------

func TestServerConfigBuilder(t *testing.T) {
	cfg := NewServerConfigBuilder().
		WithHost("0.0.0.0").
		WithPort(443).
		WithTLS("cert.pem").
		Build()
	if cfg.Host != "0.0.0.0" {
		t.Errorf("Expected host 0.0.0.0, got %s", cfg.Host)
	}
	if cfg.Port != 443 {
		t.Errorf("Expected port 443, got %d", cfg.Port)
	}
	if !cfg.TLS || cfg.CertFile != "cert.pem" {
		t.Errorf("Expected TLS with cert.pem, got TLS=%v cert=%s", cfg.TLS, cfg.CertFile)
	}
}

// ---------------------------------------------------------------------------
// Strategy
// ---------------------------------------------------------------------------

func TestSortStrategies(t *testing.T) {
	data := []int{5, 3, 1, 4, 2}
	expected := []int{1, 2, 3, 4, 5}

	for _, strategy := range []SortStrategy{bubbleSortStrategy{}, insertionSortStrategy{}} {
		sorter := NewSorter(strategy)
		result := sorter.Sort(data)
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("%s: expected %v, got %v", strategy.Name(), expected, result)
		}
	}
	// Verify original is unchanged
	if data[0] != 5 {
		t.Error("Sort strategies should not mutate original slice")
	}
}

// ---------------------------------------------------------------------------
// Observer
// ---------------------------------------------------------------------------

func TestEventBusPubSub(t *testing.T) {
	bus := NewEventBus()
	var received []string
	bus.Subscribe("test", func(event string, data interface{}) {
		received = append(received, fmt.Sprintf("%s:%v", event, data))
	})
	bus.Subscribe("test", func(event string, data interface{}) {
		received = append(received, "second")
	})
	bus.Publish("test", "payload")

	if len(received) != 2 {
		t.Errorf("Expected 2 notifications, got %d", len(received))
	}
	if received[0] != "test:payload" {
		t.Errorf("Expected 'test:payload', got %s", received[0])
	}
}

// ---------------------------------------------------------------------------
// Decorator
// ---------------------------------------------------------------------------

func TestDecoratorChain(t *testing.T) {
	var tr Transformer = baseTransformer{}
	tr = trimDecorator{inner: tr}
	tr = uppercaseDecorator{inner: tr}
	tr = bracketDecorator{inner: tr}

	result := tr.Transform("  hello  ")
	expected := "[HELLO]"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

// ---------------------------------------------------------------------------
// Iterator
// ---------------------------------------------------------------------------

func TestIntRange(t *testing.T) {
	var collected []int
	for v := range IntRange(3, 7) {
		collected = append(collected, v)
	}
	expected := []int{3, 4, 5, 6}
	if !reflect.DeepEqual(collected, expected) {
		t.Errorf("Expected %v, got %v", expected, collected)
	}
}

func TestFibonacci(t *testing.T) {
	fib := Fibonacci()
	expected := []int{0, 1, 1, 2, 3, 5, 8, 13}
	for i, want := range expected {
		got := fib()
		if got != want {
			t.Errorf("Fibonacci()[%d] = %d, want %d", i, got, want)
		}
	}
}

// ---------------------------------------------------------------------------
// Adapter
// ---------------------------------------------------------------------------

func TestLegacyToModernAdapter(t *testing.T) {
	legacy := legacyPrinterImpl{text: "data"}
	modern := NewModernPrinterFromLegacy(legacy)
	result := modern.PrintFormatted("PREFIX")
	expected := "PREFIX: data"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

// ---------------------------------------------------------------------------
// Template Method
// ---------------------------------------------------------------------------

func TestGenerateReport(t *testing.T) {
	report := GenerateReport(salesReport{})
	if !strings.Contains(report, "Sales Report") {
		t.Errorf("Expected report to contain 'Sales Report', got %q", report)
	}
	if !strings.Contains(report, "$1.2M") {
		t.Errorf("Expected report to contain data '$1.2M', got %q", report)
	}
}

// ---------------------------------------------------------------------------
// Command
// ---------------------------------------------------------------------------

func TestCommandWithUndo(t *testing.T) {
	editor := &TextEditor{}
	history := &CommandHistory{}

	history.Execute(&appendCommand{editor: editor, text: "Hello"})
	if editor.content != "Hello" {
		t.Errorf("Expected 'Hello', got %q", editor.content)
	}

	history.Execute(&appendCommand{editor: editor, text: " World"})
	if editor.content != "Hello World" {
		t.Errorf("Expected 'Hello World', got %q", editor.content)
	}

	history.Execute(&uppercaseCommand{editor: editor})
	if editor.content != "HELLO WORLD" {
		t.Errorf("Expected 'HELLO WORLD', got %q", editor.content)
	}

	history.Undo() // undo uppercase
	if editor.content != "Hello World" {
		t.Errorf("After undo uppercase expected 'Hello World', got %q", editor.content)
	}

	history.Undo() // undo second append
	if editor.content != "Hello" {
		t.Errorf("After undo append expected 'Hello', got %q", editor.content)
	}
}

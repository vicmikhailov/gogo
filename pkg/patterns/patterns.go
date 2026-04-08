// Package patterns showcases GoF design patterns implemented in idiomatic Go.
//
// For a Java developer:
//   - Go doesn't have classes, but it implements GoF patterns using structs,
//     interfaces, and functional features.
//   - Composition is preferred over inheritance (which doesn't exist).
//   - Java-ism to avoid: Over-engineering with too many design patterns.
//     Go prefers simple, direct code. Only use patterns when they actually
//     solve a concrete problem (like `sync.Once` for singletons or functional options
//     for complex configuration).
package patterns

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

// ---------------------------------------------------------------------------
// 1. Singleton (Creational) – uses sync.Once for thread-safe lazy init
// ---------------------------------------------------------------------------

// AppConfig is a process-wide configuration singleton.
//
// For a Java developer:
// - Java equivalent: enum Singleton { INSTANCE; } or Double-Checked Locking.
// - Go idiom: Use `sync.Once` to ensure a piece of code runs exactly once.
type AppConfig struct {
	AppName   string
	Version   string
	DebugMode bool
}

var (
	appConfigInstance *AppConfig
	appConfigOnce     sync.Once
)

// GetAppConfig returns the single AppConfig instance, creating it on first call.
// Java equivalent: `public static AppConfig getInstance()` with double-checked locking or an enum.
func GetAppConfig() *AppConfig {
	appConfigOnce.Do(func() {
		appConfigInstance = &AppConfig{
			AppName:   "GoShowcase",
			Version:   "1.0.0",
			DebugMode: false,
		}
	})
	return appConfigInstance
}

// ---------------------------------------------------------------------------
// 2. Factory Method (Creational) – returns interface, hides concrete types
// ---------------------------------------------------------------------------

// Notifier represents any notification channel.
//
// For a Java developer:
// - Java equivalent: Interface-based factory returns a concrete implementation.
// - Go idiom: Return an interface type, but keep concrete structs private (lowercase).
type Notifier interface {
	Send(message string) string
}

type emailNotifier struct{ address string }
type smsNotifier struct{ phone string }
type slackNotifier struct{ channel string }

func (e emailNotifier) Send(msg string) string {
	return fmt.Sprintf("[Email -> %s] %s", e.address, msg)
}
func (s smsNotifier) Send(msg string) string {
	return fmt.Sprintf("[SMS -> %s] %s", s.phone, msg)
}
func (s slackNotifier) Send(msg string) string {
	return fmt.Sprintf("[Slack -> #%s] %s", s.channel, msg)
}

// NewNotifier is a factory method that returns the correct Notifier by type name.
// Java equivalent: `public static Notifier createNotifier(String type, String destination)`
func NewNotifier(nType, destination string) Notifier {
	switch nType {
	case "email":
		return emailNotifier{address: destination}
	case "sms":
		return smsNotifier{phone: destination}
	case "slack":
		return slackNotifier{channel: destination}
	default:
		return emailNotifier{address: destination}
	}
}

// ---------------------------------------------------------------------------
// 3. Builder (Creational) – fluent API for complex object construction
// ---------------------------------------------------------------------------

// ServerConfig holds HTTP server settings built via a fluent builder.
//
// For a Java developer:
// - Java equivalent: Lombok @Builder or manual static Inner Builder class.
// - Go idiom: Fluent API using methods that return the builder pointer.
type ServerConfig struct {
	Host         string
	Port         int
	ReadTimeout  int
	WriteTimeout int
	MaxConns     int
	TLS          bool
	CertFile     string
}

// ServerConfigBuilder provides a fluent API for building ServerConfig.
//
// For a Java developer:
// - Similar to a static inner class `Builder` in a Java class.
// - All methods return a pointer to the builder to allow chaining.
type ServerConfigBuilder struct {
	config ServerConfig
}

// NewServerConfigBuilder creates a new instance of ServerConfigBuilder with default values.
// Java equivalent: `public static ServerConfigBuilder builder()`
func NewServerConfigBuilder() *ServerConfigBuilder {
	return &ServerConfigBuilder{
		config: ServerConfig{Host: "localhost", Port: 8080, MaxConns: 100},
	}
}

// WithHost sets the server host.
func (b *ServerConfigBuilder) WithHost(host string) *ServerConfigBuilder {
	b.config.Host = host
	return b
}

// WithPort sets the server port.
func (b *ServerConfigBuilder) WithPort(port int) *ServerConfigBuilder {
	b.config.Port = port
	return b
}

// WithReadTimeout sets the read timeout in seconds.
func (b *ServerConfigBuilder) WithReadTimeout(seconds int) *ServerConfigBuilder {
	b.config.ReadTimeout = seconds
	return b
}

// WithWriteTimeout sets the write timeout in seconds.
func (b *ServerConfigBuilder) WithWriteTimeout(seconds int) *ServerConfigBuilder {
	b.config.WriteTimeout = seconds
	return b
}

// WithMaxConns sets the maximum number of simultaneous connections.
func (b *ServerConfigBuilder) WithMaxConns(n int) *ServerConfigBuilder {
	b.config.MaxConns = n
	return b
}

// WithTLS enables TLS and sets the certificate file.
func (b *ServerConfigBuilder) WithTLS(certFile string) *ServerConfigBuilder {
	b.config.TLS = true
	b.config.CertFile = certFile
	return b
}

// Build returns the final ServerConfig.
// Java equivalent: `builder.build()`
func (b *ServerConfigBuilder) Build() ServerConfig {
	return b.config
}

// ---------------------------------------------------------------------------
// 4. Strategy (Behavioral) – interchangeable algorithms behind an interface
// ---------------------------------------------------------------------------

// SortStrategy defines a sorting algorithm.
//
// For a Java developer:
// - Java equivalent: Passing a Comparator to Collections.sort().
// - Go idiom: Define an interface for the algorithm and inject it.
// - This is the "behavioral" part of the pattern.
type SortStrategy interface {
	Sort(data []int) []int
	Name() string
}

// bubbleSortStrategy implements a simple bubble sort.
type bubbleSortStrategy struct{}

func (bubbleSortStrategy) Name() string { return "BubbleSort" }
func (bubbleSortStrategy) Sort(data []int) []int {
	n := len(data)
	result := make([]int, n)
	copy(result, data)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			if result[j] > result[j+1] {
				result[j], result[j+1] = result[j+1], result[j]
			}
		}
	}
	return result
}

// insertionSortStrategy implements insertion sort.
type insertionSortStrategy struct{}

func (insertionSortStrategy) Name() string { return "InsertionSort" }
func (insertionSortStrategy) Sort(data []int) []int {
	n := len(data)
	result := make([]int, n)
	copy(result, data)
	for i := 1; i < n; i++ {
		key := result[i]
		j := i - 1
		for j >= 0 && result[j] > key {
			result[j+1] = result[j]
			j--
		}
		result[j+1] = key
	}
	return result
}

// Sorter uses a pluggable SortStrategy (like Java's Comparator/Strategy).
//
// For a Java developer:
// - Similar to how `Collections.sort()` or `List.sort()` works.
type Sorter struct {
	strategy SortStrategy
}

// NewSorter creates a new Sorter with the given strategy.
// Java equivalent: `public Sorter(SortStrategy strategy)`
func NewSorter(strategy SortStrategy) *Sorter {
	return &Sorter{strategy: strategy}
}

// SetStrategy allows changing the sorting strategy at runtime.
// Java equivalent: `public void setStrategy(SortStrategy strategy)`
func (s *Sorter) SetStrategy(strategy SortStrategy) {
	s.strategy = strategy
}

// Sort delegates the sorting operation to the current strategy.
func (s *Sorter) Sort(data []int) []int {
	return s.strategy.Sort(data)
}

// ---------------------------------------------------------------------------
// 5. Observer (Behavioral) – event-driven pub/sub
// ---------------------------------------------------------------------------

// EventListener is the observer callback type.
//
// For a Java developer:
// - Java equivalent: PropertyChangeListener or Spring ApplicationEventPublisher.
// - Go idiom: Use a slice of functions (listeners) and a mutex for safety.
type EventListener func(eventName string, data interface{})

// EventBus is a simple publish-subscribe system (comparable to Java EventBus / listeners).
//
// For a Java developer:
// - Similar to Google Guava's EventBus or Spring's ApplicationEventPublisher.
type EventBus struct {
	mu        sync.RWMutex
	listeners map[string][]EventListener
}

// NewEventBus creates a new EventBus.
func NewEventBus() *EventBus {
	return &EventBus{listeners: make(map[string][]EventListener)}
}

// Subscribe registers a listener for a specific event.
// Java equivalent: `eventBus.addListener(eventType, listener)`
func (eb *EventBus) Subscribe(event string, listener EventListener) {
	eb.mu.Lock()
	defer eb.mu.Unlock()
	eb.listeners[event] = append(eb.listeners[event], listener)
}

// Publish fires an event, notifying all registered listeners.
// Java equivalent: `eventBus.post(event)`
func (eb *EventBus) Publish(event string, data interface{}) {
	eb.mu.RLock()
	defer eb.mu.RUnlock()
	for _, listener := range eb.listeners[event] {
		listener(event, data)
	}
}

// ---------------------------------------------------------------------------
// 6. Decorator (Structural) – wrapping behaviour around an interface
// ---------------------------------------------------------------------------

// Transformer processes a string (comparable to Java's Function<String,String>).
//
// For a Java developer:
// - Java equivalent: BufferedInputStream wraps a FileInputStream.
// - Go idiom: Wrap an interface with another struct that implements the same interface.
type Transformer interface {
	Transform(input string) string
}

// baseTransformer returns the input unchanged.
type baseTransformer struct{}

// Transform returns the input as-is for baseTransformer.
func (baseTransformer) Transform(input string) string { return input }

// uppercaseDecorator wraps a Transformer and uppercases its result.
type uppercaseDecorator struct{ inner Transformer }

// Transform for uppercaseDecorator converts the input to uppercase.
func (d uppercaseDecorator) Transform(input string) string {
	return strings.ToUpper(d.inner.Transform(input))
}

// trimDecorator wraps a Transformer and trims whitespace.
type trimDecorator struct{ inner Transformer }

// Transform for trimDecorator removes leading and trailing whitespace.
func (d trimDecorator) Transform(input string) string {
	return strings.TrimSpace(d.inner.Transform(input))
}

// bracketDecorator wraps a Transformer and adds brackets.
type bracketDecorator struct{ inner Transformer }

// Transform for bracketDecorator wraps the input in square brackets.
func (d bracketDecorator) Transform(input string) string {
	return "[" + d.inner.Transform(input) + "]"
}

// ---------------------------------------------------------------------------
// 7. Iterator (Behavioral) – channel-based and closure-based iterators
// ---------------------------------------------------------------------------

// IntRange returns a channel that yields integers in [start, end).
//
// For a Java developer:
// - Java equivalent: `IntStream.range(start, end).iterator()`.
// - Go idiom: Use a channel as an iterator. The `for range` loop makes it feel natural.
func IntRange(start, end int) <-chan int {
	ch := make(chan int)
	go func() {
		defer close(ch)
		for i := start; i < end; i++ {
			ch <- i
		}
	}()
	return ch
}

// Fibonacci returns a closure-based iterator that yields Fibonacci numbers.
//
// For a Java developer:
// - Java equivalent: A custom Iterator implementation with internal state.
// - Go idiom: Use a closure to capture and maintain state.
func Fibonacci() func() int {
	a, b := 0, 1
	return func() int {
		val := a
		a, b = b, a+b
		return val
	}
}

// ---------------------------------------------------------------------------
// 8. Adapter (Structural) – adapting incompatible interfaces
// ---------------------------------------------------------------------------

// LegacyPrinter has an old-style API.
//
// For a Java developer:
// - Java equivalent: Wrapper classes that convert one interface to another.
// - Go idiom: Embed the old type into a new struct that satisfies the new interface.
type LegacyPrinter interface {
	PrintStored() string
}

// ModernPrinter has a new-style API.
type ModernPrinter interface {
	PrintFormatted(prefix string) string
}

type legacyPrinterImpl struct{ text string }

// PrintStored returns the stored text for legacyPrinterImpl.
func (lp legacyPrinterImpl) PrintStored() string {
	return lp.text
}

// legacyToModernAdapter adapts LegacyPrinter to the ModernPrinter interface.
type legacyToModernAdapter struct {
	legacy LegacyPrinter
}

// PrintFormatted adapts the legacy PrintStored to the modern PrintFormatted.
func (a legacyToModernAdapter) PrintFormatted(prefix string) string {
	return fmt.Sprintf("%s: %s", prefix, a.legacy.PrintStored())
}

// NewModernPrinterFromLegacy wraps a LegacyPrinter as a ModernPrinter.
// Java equivalent: `public static ModernPrinter adapt(LegacyPrinter legacy)`
func NewModernPrinterFromLegacy(lp LegacyPrinter) ModernPrinter {
	return legacyToModernAdapter{legacy: lp}
}

// ---------------------------------------------------------------------------
// 9. Template Method (Behavioral) – via embedding + interface override
// ---------------------------------------------------------------------------

// ReportGenerator defines the steps to generate a report.
//
// For a Java developer:
// - Java equivalent: Abstract class with some final methods and some abstract ones.
// - Go idiom: Use an interface to define the "hooks" and a function to orchestrate them.
type ReportGenerator interface {
	CollectData() string
	FormatReport(data string) string
}

// GenerateReport is the template method that orchestrates the steps.
func GenerateReport(rg ReportGenerator) string {
	data := rg.CollectData()
	return rg.FormatReport(data)
}

type salesReport struct{}

// CollectData for salesReport simulates data collection.
func (salesReport) CollectData() string { return "Q4 Sales: $1.2M" }

// FormatReport for salesReport formats the collected data.
func (salesReport) FormatReport(data string) string {
	return fmt.Sprintf("=== Sales Report ===\n%s\n====================", data)
}

type inventoryReport struct{}

// CollectData for inventoryReport simulates data collection.
func (inventoryReport) CollectData() string { return "Widgets: 450, Gadgets: 230" }

// FormatReport for inventoryReport formats the collected data.
func (inventoryReport) FormatReport(data string) string {
	return fmt.Sprintf("--- Inventory Report ---\n%s\n------------------------", data)
}

// ---------------------------------------------------------------------------
// 10. Command (Behavioral) – encapsulate operations as objects
// ---------------------------------------------------------------------------

// Command encapsulates an action and its undo.
//
// For a Java developer:
// - Java equivalent: `Runnable` or a custom `Command` interface with `execute()` and `undo()`.
// - Go idiom: Structs that implement an Execute/Undo interface.
type Command interface {
	Execute() string
	Undo() string
}

// TextEditor maintains a buffer and supports command-based editing.
//
// For a Java developer:
// - This is the "receiver" in the Command pattern.
type TextEditor struct {
	content string
}

// appendCommand adds text to the editor.
type appendCommand struct {
	editor *TextEditor
	text   string
}

// Execute for appendCommand appends the text to the editor.
func (c *appendCommand) Execute() string {
	c.editor.content += c.text
	return fmt.Sprintf("Appended '%s' -> \"%s\"", c.text, c.editor.content)
}

// Undo for appendCommand removes the appended text from the editor.
func (c *appendCommand) Undo() string {
	c.editor.content = c.editor.content[:len(c.editor.content)-len(c.text)]
	return fmt.Sprintf("Undid append '%s' -> \"%s\"", c.text, c.editor.content)
}

// uppercaseCommand converts the editor's content to uppercase.
type uppercaseCommand struct {
	editor   *TextEditor
	previous string
}

// Execute for uppercaseCommand uppercases the editor's content.
func (c *uppercaseCommand) Execute() string {
	c.previous = c.editor.content
	c.editor.content = strings.ToUpper(c.editor.content)
	return fmt.Sprintf("Uppercased -> \"%s\"", c.editor.content)
}

// Undo for uppercaseCommand restores the content to its previous state.
func (c *uppercaseCommand) Undo() string {
	c.editor.content = c.previous
	return fmt.Sprintf("Undid uppercase -> \"%s\"", c.editor.content)
}

// CommandHistory records commands for undo support.
//
// For a Java developer:
// - Similar to a stack of Command objects.
type CommandHistory struct {
	history []Command
}

// Execute executes a command and adds it to the history.
// Java equivalent: `public String execute(Command cmd)`
func (h *CommandHistory) Execute(cmd Command) string {
	result := cmd.Execute()
	h.history = append(h.history, cmd)
	return result
}

// Undo reverts the last executed command from the history.
// Java equivalent: `public String undo()`
func (h *CommandHistory) Undo() string {
	if len(h.history) == 0 {
		return "Nothing to undo"
	}
	cmd := h.history[len(h.history)-1]
	h.history = h.history[:len(h.history)-1]
	return cmd.Undo()
}

// ---------------------------------------------------------------------------
// 11. Functional Options (Idiomatic Go) – flexible configuration
// ---------------------------------------------------------------------------

// Client represents a configured service client configured via functional options.
//
// For a Java developer:
//   - Java equivalent: Overloaded constructors or a Builder pattern.
//   - Go idiom: Functional Options allow for a clean API with sensible defaults.
//     It's widely used in popular libraries like gRPC and Zap.
type Client struct {
	host    string
	timeout time.Duration
	retry   int
}

// Option is a function type that modifies the Client.
type Option func(*Client)

// WithTimeout returns an Option that sets the client timeout.
// Java equivalent: `builder.setTimeout(t)`
func WithTimeout(t time.Duration) Option {
	return func(c *Client) {
		c.timeout = t
	}
}

// WithRetry returns an Option that sets the retry count.
// Java equivalent: `builder.setRetry(r)`
func WithRetry(r int) Option {
	return func(c *Client) {
		c.retry = r
	}
}

// NewClient creates a new Client with the given host and optional settings.
// Java equivalent: `public Client(String host, Option... options)`
func NewClient(host string, opts ...Option) *Client {
	// Default values
	client := &Client{
		host:    host,
		timeout: 5 * time.Second,
		retry:   3,
	}

	// Apply options
	for _, opt := range opts {
		opt(client)
	}

	return client
}

// ---------------------------------------------------------------------------
// RunPatternsDemo
// ---------------------------------------------------------------------------

// RunPatternsDemo showcases GoF design patterns implemented in idiomatic Go.
func RunPatternsDemo() {
	fmt.Println("--- GoF Design Patterns Demo ---")

	// 1. Singleton
	fmt.Println("1. Singleton (sync.Once):")
	cfg1 := GetAppConfig()
	cfg2 := GetAppConfig()
	fmt.Printf("   Same instance? %v (AppName: %s)\n", cfg1 == cfg2, cfg1.AppName)

	// 2. Factory Method
	fmt.Println("2. Factory Method:")
	notifiers := []Notifier{
		NewNotifier("email", "dev@example.com"),
		NewNotifier("sms", "+1234567890"),
		NewNotifier("slack", "engineering"),
	}
	for _, n := range notifiers {
		fmt.Printf("   %s\n", n.Send("Build succeeded"))
	}

	// 3. Builder
	fmt.Println("3. Builder (fluent API):")
	server := NewServerConfigBuilder().
		WithHost("0.0.0.0").
		WithPort(443).
		WithReadTimeout(30).
		WithWriteTimeout(60).
		WithMaxConns(1000).
		WithTLS("/etc/ssl/cert.pem").
		Build()
	fmt.Printf("   Server: %s:%d TLS=%v MaxConns=%d\n",
		server.Host, server.Port, server.TLS, server.MaxConns)

	// 4. Strategy
	fmt.Println("4. Strategy (swappable sort algorithms):")
	data := []int{64, 34, 25, 12, 22, 11, 90}
	sorter := NewSorter(bubbleSortStrategy{})
	fmt.Printf("   %s: %v\n", sorter.strategy.Name(), sorter.Sort(data))
	sorter.SetStrategy(insertionSortStrategy{})
	fmt.Printf("   %s: %v\n", sorter.strategy.Name(), sorter.Sort(data))

	// 5. Observer
	fmt.Println("5. Observer (EventBus pub/sub):")
	bus := NewEventBus()
	var log []string
	bus.Subscribe("user.created", func(event string, data interface{}) {
		log = append(log, fmt.Sprintf("Handler1: %s -> %v", event, data))
	})
	bus.Subscribe("user.created", func(event string, data interface{}) {
		log = append(log, fmt.Sprintf("Handler2: %s -> %v", event, data))
	})
	bus.Publish("user.created", "alice@example.com")
	for _, entry := range log {
		fmt.Printf("   %s\n", entry)
	}

	// 6. Decorator
	fmt.Println("6. Decorator (composable transformations):")
	var t Transformer = baseTransformer{}
	t = trimDecorator{inner: t}
	t = uppercaseDecorator{inner: t}
	t = bracketDecorator{inner: t}
	fmt.Printf("   Result: %s\n", t.Transform("  hello world  "))

	// 7. Iterator
	fmt.Println("7. Iterator (channel-based & closure-based):")
	fmt.Print("   Range [0,5): ")
	for val := range IntRange(0, 5) {
		fmt.Printf("%d ", val)
	}
	fmt.Println()
	fmt.Print("   Fibonacci(10): ")
	fib := Fibonacci()
	for i := 0; i < 10; i++ {
		fmt.Printf("%d ", fib())
	}
	fmt.Println()

	// 8. Adapter
	fmt.Println("8. Adapter (legacy -> modern interface):")
	legacy := legacyPrinterImpl{text: "important data"}
	modern := NewModernPrinterFromLegacy(legacy)
	fmt.Printf("   %s\n", modern.PrintFormatted("Adapted"))

	// 9. Template Method
	fmt.Println("9. Template Method (via interface composition):")
	for _, rg := range []ReportGenerator{salesReport{}, inventoryReport{}} {
		report := GenerateReport(rg)
		for _, line := range strings.Split(report, "\n") {
			fmt.Printf("   %s\n", line)
		}
	}

	// 10. Command (with undo)
	fmt.Println("10. Command (execute + undo):")
	editor := &TextEditor{}
	history := &CommandHistory{}
	fmt.Printf("   %s\n", history.Execute(&appendCommand{editor: editor, text: "Hello"}))
	fmt.Printf("   %s\n", history.Execute(&appendCommand{editor: editor, text: " World"}))
	fmt.Printf("   %s\n", history.Execute(&uppercaseCommand{editor: editor}))
	fmt.Printf("   %s\n", history.Undo())
	fmt.Printf("   %s\n", history.Undo())

	// 11. Functional Options
	fmt.Println("11. Functional Options (idiomatic Go configuration):")
	c1 := NewClient("localhost")
	c2 := NewClient("api.example.com", WithTimeout(30*time.Second), WithRetry(5))
	fmt.Printf("   Client 1: host=%s, timeout=%v, retry=%d\n", c1.host, c1.timeout, c1.retry)
	fmt.Printf("   Client 2: host=%s, timeout=%v, retry=%d\n", c2.host, c2.timeout, c2.retry)

	fmt.Println("--- GoF Design Patterns Demo End ---")
}

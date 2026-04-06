// Package iosystem showcases how Go excels in I/O and System-level programming.
package iosystem

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

// RunIOSystemDemo showcases Go's capabilities in I/O and System-level programming.
//
// For a Java developer:
//   - Go's `os` and `io` packages are equivalent to `java.io` and `java.nio`.
//   - `defer` is extensively used for resource management (like try-with-resources).
//   - The `flag` package is the standard way to build CLI tools (similar to Apache Commons CLI).
//   - The `net` package provides low-level networking (like `java.net.Socket`).
//   - `os/exec` is equivalent to `ProcessBuilder` for running external commands.
//   - `os/signal` is equivalent to `Runtime.addShutdownHook()`.
func RunIOSystemDemo() {
	fmt.Println("--- I/O and System-Level Programming Demo ---")

	// 1. File Manipulation
	// Java comparison: `Files.write()`, `Files.readAllBytes()`, `Scanner`.
	// Go idiom: `os.WriteFile` and `os.ReadFile` for small files; `bufio` for large files.
	fmt.Println("1. File Manipulation (I/O):")
	tempFile := "test_demo.txt"
	defer os.Remove(tempFile) // Cleanup

	// Write (Create/Overwrite)
	err := os.WriteFile(tempFile, []byte("Hello, Go I/O!\nLine 2\nLine 3\n"), 0644)
	if err != nil {
		fmt.Printf("   Error writing file: %v\n", err)
	} else {
		fmt.Println("   Created and wrote to test_demo.txt")
	}

	// Append
	// Java comparison: `StandardOpenOption.APPEND`.
	f, err := os.OpenFile(tempFile, os.O_APPEND|os.O_WRONLY, 0644)
	if err == nil {
		f.WriteString("Appended Line\n")
		f.Close()
		fmt.Println("   Appended text to file")
	}

	// Read entire file
	content, _ := os.ReadFile(tempFile)
	fmt.Printf("   File Content (size: %d bytes):\n", len(content))
	fmt.Print("   " + strings.ReplaceAll(string(content), "\n", "\n   "))

	// Buffered Reading (Line by Line)
	// Java comparison: `BufferedReader.readLine()`.
	fmt.Println("\n   Buffered reading (Line by Line):")
	f, _ = os.Open(tempFile)
	scanner := bufio.NewScanner(f)
	lineCount := 0
	for scanner.Scan() {
		lineCount++
		fmt.Printf("      Line %d: %s\n", lineCount, scanner.Text())
	}
	f.Close()

	// File Stats
	// Java comparison: `Files.readAttributes()`.
	info, _ := os.Stat(tempFile)
	fmt.Printf("   File Stats: Mode: %v, ModTime: %s\n", info.Mode(), info.ModTime().Format(time.Kitchen))

	// 2. Directory Walking
	// Java comparison: `Files.walkFileTree()`.
	fmt.Println("\n2. Directory Walking (filepath):")
	count := 0
	_ = filepath.WalkDir(".", func(path string, d os.DirEntry, err error) error {
		if err != nil || count >= 5 {
			return nil
		}
		if !d.IsDir() {
			fmt.Printf("   Found file: %s\n", path)
			count++
		}
		return nil
	})

	// 3. System: Environment Variables & CLI (Simulated)
	// Java comparison: `System.getenv()` and CLI argument parsing libraries.
	fmt.Println("\n3. System: Environment & CLI Flags:")
	os.Setenv("GO_DEMO_ENV", "ExpertMode")
	fmt.Printf("   Env variable GO_DEMO_ENV: %s\n", os.Getenv("GO_DEMO_ENV"))

	// Demonstrating how flags work (without parsing os.Args which would affect the whole process)
	fs := flag.NewFlagSet("demo", flag.ContinueOnError)
	msg := fs.String("msg", "default", "a message to print")
	_ = fs.Parse([]string{"-msg", "Go is great for CLI!"})
	fmt.Printf("   CLI Flag 'msg': %s\n", *msg)

	// 4. System: Executing External Commands
	// Java comparison: `ProcessBuilder` or `Runtime.exec()`.
	fmt.Println("\n4. System: Executing Commands (os/exec):")
	cmd := exec.Command("echo", "Hello from external command!")
	output, _ := cmd.Output()
	fmt.Printf("   Command output: %s", string(output))

	// 5. Networking: Simple TCP Connection (Loopback)
	// Java comparison: `ServerSocket` and `Socket`.
	fmt.Println("\n5. Networking: Simple TCP Echo (Loopback):")
	ln, err := net.Listen("tcp", "127.0.0.1:0") // Listen on random available port
	if err == nil {
		addr := ln.Addr().String()
		fmt.Printf("   Listening on %s\n", addr)

		go func() {
			conn, _ := ln.Accept()
			if conn != nil {
				io.Copy(conn, conn) // Echo back
				conn.Close()
			}
		}()

		conn, _ := net.Dial("tcp", addr)
		if conn != nil {
			fmt.Fprintf(conn, "Ping")
			reply := make([]byte, 4)
			conn.Read(reply)
			fmt.Printf("   Sent: Ping, Received: %s\n", string(reply))
			conn.Close()
		}
		ln.Close()
	}

	// 6. Signal Handling (Conceptual)
	// Java comparison: `Runtime.getRuntime().addShutdownHook()`.
	// Go idiom: Use a channel to receive OS signals like SIGINT or SIGTERM.
	fmt.Println("\n6. Signal Handling (Conceptual):")
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("   Ready to catch SIGINT/SIGTERM (Ctrl+C).")

	// Create a context that times out after 50ms so we don't block the demo
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	select {
	case sig := <-sigChan:
		fmt.Printf("   Received signal: %v\n", sig)
	case <-ctx.Done():
		fmt.Println("   (Signal demo timeout - normally you'd wait here for real shutdown)")
	}

	fmt.Println("--- I/O and System-Level Programming Demo End ---")
}

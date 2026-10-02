// Package main is the entrypoint for the gogo showcase.
//
// For a Java developer:
//   - In Java, a Spring Boot or Enterprise application typically starts in a
//     class with `public static void main(String[] args)`.
//   - In Go, executable programs belong to package `main` with a top-level `func main()`.
//   - This orchestrator runs each showcase package in logical progression,
//     culminating in an active HTTP web server.
package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"gogo/pkg/advanced"
	"gogo/pkg/basictypes"
	"gogo/pkg/collections"
	"gogo/pkg/commonlibs"
	"gogo/pkg/concurrency"
	"gogo/pkg/errors"
	"gogo/pkg/generics"
	"gogo/pkg/interfaces"
	"gogo/pkg/iosystem"
	"gogo/pkg/patterns"
	"gogo/pkg/syntax"
	"gogo/pkg/web"
)

func main() {
	fmt.Println("======================================================================")
	fmt.Println("             gogo: The Go Showcase for Java Developers               ")
	fmt.Println("======================================================================")

	// 1. Core Syntax (Variables, pointers, defer, zero-values, multiple returns)
	syntax.RunSyntaxDemo()
	fmt.Println()

	// 2. Concurrency (Goroutines, channels, sync primitives, worker pools)
	concurrency.RunConcurrencyDemo()
	fmt.Println()

	// 3. Interfaces (Duck typing, polymorphism, type assertions)
	interfaces.RunInterfacesDemo()
	fmt.Println()

	// 4. Generics (Type parameters, generic collections, map/filter)
	generics.RunGenericsDemo()
	fmt.Println()

	// 5. Error Handling (Errors as values, custom errors, errors.Is/As wrapping)
	errors.RunErrorsDemo()
	fmt.Println()

	// 6. Collections (Java Stream-like operations: Filter, Reduce, GroupBy, Zip)
	collections.RunCollectionsDemo()
	fmt.Println()

	// 7. Basic Types (Slices, maps, strings, and JSON serialization)
	basictypes.RunBasicTypesDemo()
	fmt.Println()

	// 8. Design Patterns & Go Antipatterns (GoF patterns + Java-to-Go pitfalls)
	patterns.RunPatternsDemo()
	fmt.Println()

	// 9. Advanced Features (Embedding, reflection, struct tags, static asset embedding)
	advanced.RunAdvancedDemo()
	fmt.Println()

	// 10. I/O and System Programming (File I/O, directory walking, CLI flags, exec)
	iosystem.RunIOSystemDemo()
	fmt.Println()

	// 11. Common 3rd Party Libraries (zap logging, uuid, testify assertions, gin)
	commonlibs.RunCommonLibsDemo()
	fmt.Println()

	// 12. Web Server (Standard net/http, custom middleware, REST endpoints)
	web.StartWebServer("8080")

	// Wait for user input or interrupt signal before terminating
	fmt.Println("\n[Server Ready] Press Enter or Ctrl+C to stop the showcase...")
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	enterPressed := make(chan struct{})
	go func() {
		reader := bufio.NewReader(os.Stdin)
		_, _ = reader.ReadString('\n')
		close(enterPressed)
	}()

	select {
	case <-stop:
		fmt.Println("\nShutdown signal received. Goodbye!")
	case <-enterPressed:
		fmt.Println("\nShowcase completed. Goodbye!")
	}
}

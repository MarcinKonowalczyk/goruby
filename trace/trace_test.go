package trace_test

import (
	"fmt"
	"testing"

	"github.com/MarcinKonowalczyk/goruby/trace"
	"github.com/MarcinKonowalczyk/goruby/trace/printer"
	"github.com/MarcinKonowalczyk/goruby/utils"
)

func TestNewTracer(t *testing.T) {
	tracer := trace.NewTracer()
	utils.AssertEqual(t, len(tracer.Messages()), 0)
}

////////////////////////////////////////////////////////////////////////////////

// called third
func third(tr trace.Tracer) {
	defer tr.Un(tr.Trace())
	tr.Message("Second")
}

// called second
func second(tr trace.Tracer) {
	defer tr.Un(tr.Trace())
	third(tr)
	tr.Message("Third") // third since _third is called before
}

// called first
func first(tr trace.Tracer) {
	defer tr.Un(tr.Trace())
	tr.Message("First")
	second(tr)
}

func TestTraceNormal(t *testing.T) {
	tracer := trace.NewTracer()
	first(tracer)
	tracer.Done()
	messages := tracer.Messages()
	utils.AssertEqual(t, len(messages), 3)
	for _, message := range messages {
		printer.PrettyPrint(message, printer.MULTILINE)
	}
}

////////////////////////////////////////////////////////////////////////////////

func first_gap(tr trace.Tracer) {
	defer tr.Un(tr.Trace())
	tr.Message("First")
	second_gap(tr)
}

func second_gap(tr trace.Tracer) {
	// oops! we do not register the intermediate function with the tracer!
	// defer tr.Un(tr.Trace())
	tr.Message("Second")
	third_gap(tr)
}

func third_gap(tr trace.Tracer) {
	// and we did not register this one either!
	// defer tr.Un(tr.Trace())
	tr.Message("Third")
}

func TestTraceGap(t *testing.T) {
	tracer := trace.NewTracer()
	first_gap(tracer)
	tracer.Done()
	messages := tracer.Messages()
	utils.AssertEqual(t, len(messages), 3)
	for _, message := range messages {
		// fmt.Printf("%s\n", message)
		printer.PrettyPrint(message, printer.MULTILINE)
	}
}

////////////////////////////////////////////////////////////////////////////////

func TestTraceWithout(t *testing.T) {
	tracer := trace.NewTracer()
	tracer.Message("First")
	tracer.Done()
	messages := tracer.Messages()
	utils.AssertEqual(t, len(messages), 1)
	for _, message := range messages {
		// fmt.Printf("%s\n", message)
		printer.PrettyPrint(message, printer.MULTILINE)
	}
}

////////////////////////////////////////////////////////////////////////////////

func TestTraceToWalkable(t *testing.T) {
	tracer := trace.NewTracer()
	first(tracer)
	_, err := tracer.ToWalkable()
	utils.AssertError(t, err, "not walkable")
	tracer.Done()
	walkable, err := tracer.ToWalkable()
	utils.AssertNoError(t, err)
	fmt.Println("Walkable:")
	walkable.Walk(func(node trace.Node) error {
		// fmt.Println(node, "prev:", node.Prev(), "next:", node.Next())
		fmt.Println(node)
		return nil
	})
}

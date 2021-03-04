package goauto

import (
	"path/filepath"
	"testing"
)

func TestPipeline(t *testing.T) {
	p := NewPipeline("Pipline Name", Silent)
	if p == nil {
		t.Errorf("Failed to create a Pipeline\n")
	}
	wf := Workflow{}
	p.Add(&wf)
}

func TestPipelineRec(t *testing.T) {
	p := NewPipeline("Test Pipeline", Silent)
	tp := filepath.Join("src", "github.com", "argylelabcoat", "goauto")
	err := p.WatchRecursive(tp, IgnoreHidden)
	if err != nil {
		t.Errorf("WatchRecursive failed %v\n", err)
	}

	tp = filepath.Join("src", "bogus", "bogus", "bogus")
	err = p.WatchRecursive(tp, IgnoreHidden)
	if err == nil {
		t.Errorf("WatchRecursive allowed bogus path %v\n", tp)
	}
}

/* Not a reliable test. Depends on speed of the fsnotify events
func TestPipelineConcurrency(t *testing.T) {
	p := NewPipeline("Test Pipeline", Verbose)
	tp := filepath.Join("src", "github.com", "argylelabcoat", "goauto", "testing")
	err := p.WatchRecursive(tp, IgnoreHidden)
	if err != nil {
		t.Errorf("WatchRecursive failed %v\n", err)
	}

	wf := NewWorkflow()
	p.Add(wf)

	// Run Pipeline concurrently
	go p.Start()

	atp, err := AbsPath(tp)
	if err != nil {
		t.Fatal(err)
	}

	// Add sub directories to detect changes
	pl := make([]string, 0, 2)
	for i := 0; i < 2; i++ {
		n := filepath.Join(atp, strconv.Itoa(i))
		pl = append(pl, n)
		os.Mkdir(n, 0744)
	}

	// This sucks to test!
	// No gurantees about how long before fsnotify is triggered
	time.Sleep(5 * time.Second)

	var fnd int
	for i, n := range pl {
		for _, v := range p.Watches {
			if n == v {
				fnd = i
			}
		}
		if fnd != i {
			t.Errorf("Pipeline failed to detect new dir %v\n", n)
		}
	}

	for i := 0; i < 10; i++ {
		n := filepath.Join(atp, strconv.Itoa(i))
		os.Remove(n)
	}

	err = p.Stop()
	if err != nil {
		t.Error(err)
	}
}
*/

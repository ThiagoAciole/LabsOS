package catalog

import "testing"

func TestCatalogSourcesExposePriorityAdapters(t *testing.T) {
	sources := DefaultSources()
	if len(sources) < 4 {
		t.Fatalf("sources = %d, want priority plus prepared sources", len(sources))
	}
	if sources[0].Name != "BigBearCasaOS" || sources[1].Name != "CasaOS/ZimaOS AppStore" {
		t.Fatalf("priority sources = %+v", sources[:2])
	}
}

package gomodule

import (
	"bytes"
	"strings"
	"testing"

	"github.com/google/blueprint"
	"github.com/roman-mazur/bood"
)

func TestZipArchiveFactory(t *testing.T) {
	ctx := blueprint.NewContext()

	ctx.MockFileSystem(map[string][]byte{
		"Blueprints": []byte(`
			zip_archive {
			  name: "test-out",
			  srcs: ["test-src.go", "another-test-src.txt"],
			}
		`),
		"test-src.go":          nil,
		"another-test-src.txt": nil,
	})

	ctx.RegisterModuleType("zip_archive", ZipArchiveFactory)

	cfg := bood.NewConfig()

	_, errs := ctx.ParseBlueprintsFiles(".", cfg)
	if len(errs) != 0 {
		t.Fatalf("Syntax errors in the test blueprint file: %s", errs)
	}

	_, errs = ctx.PrepareBuildActions(cfg)
	if len(errs) != 0 {
		t.Errorf("Unexpected errors while preparing build actions: %s", errs)
	}
	buffer := new(bytes.Buffer)
	if err := ctx.WriteBuildFile(buffer); err != nil {
		t.Errorf("Error writing ninja file: %s", err)
	} else {
		text := buffer.String()
		t.Logf("Gennerated ninja build file:\n%s", text)
		if !strings.Contains(text, "out/archives/test-out: g.gomodule.zipArchive") {
			t.Errorf("Generated ninja file does not have build of archivator module")
		}
		if !strings.Contains(text, "inputFiles = test-src.go another-test-src.txt") {
			t.Errorf("Generated archivator input files are not correct")
		}
		if !strings.Contains(text, "outputPath = out/archives/test-out") {
			t.Errorf("Generated archivator output path is not correct")
		}

	}
}

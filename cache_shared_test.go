package main

import "testing"

func TestListPackage_ResolvesTestScopedDependencyForGeneratedTestMain(t *testing.T) {
	t.Parallel()

	prevSharedCache := sharedCache
	t.Cleanup(func() {
		sharedCache = prevSharedCache
	})

	const (
		testMainImportPath = "example.com/project/internal/rebalancing.test"
		plainImportPath    = "example.com/project/mocks/rebalancing"
	)

	// When building a test binary, Go can list a dependency only under its
	// test-scoped path "plain/import/path [pkg.test]" rather than its plain
	// import path. The generated "pkg.test" main package still refers to that
	// dependency via the plain path, so garble must resolve both forms.
	testScopedImportPath := plainImportPath + " [" + testMainImportPath + "]"
	testMainPkg := &listedPackage{
		Name:       "main",
		ImportPath: testMainImportPath,
		Imports:    []string{testScopedImportPath},
	}
	testScopedDepPkg := &listedPackage{
		Name:       "rebalancing",
		ImportPath: testScopedImportPath,
		ForTest:    "example.com/project/internal/rebalancing",
	}

	sharedCache = &sharedCacheType{
		ListedPackages: map[string]*listedPackage{
			"runtime":                  {ImportPath: "runtime"},
			testMainPkg.ImportPath:    testMainPkg,
			testScopedDepPkg.ImportPath: testScopedDepPkg,
		},
	}

	pkg, err := listPackage(testMainPkg, plainImportPath)
	if err != nil {
		t.Fatalf("listPackage() error = %v", err)
	}
	if pkg.ImportPath != testScopedDepPkg.ImportPath {
		t.Fatalf("listPackage() returned %q, want %q", pkg.ImportPath, testScopedDepPkg.ImportPath)
	}
}

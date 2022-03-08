package gocovshtest

import (
	"runtime"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/require"
)

func isWindows() bool { return runtime.GOOS == "windows" }

func TestErrorFlows(t *testing.T) {
	g := goldie.New(t, goldie.WithFixtureDir("testdata/errors"))

	t.Run("missing coverage file", func(t *testing.T) {
		mt := &modelTest{
			T:               t,
			profileFilename: "missing.cover",
			codeRoot:        "testdata/general",
		}
		initCmd := mt.init()
		initMsg := initCmd()

		mm, cmd := mt.sendWindowSizeMsg(60, 20)
		require.NotNil(t, mm)
		require.Nil(t, cmd)

		mm, cmd = mt.sendErrorMsg(initMsg)
		require.NotNil(t, mm)
		require.Nil(t, cmd)
		if isWindows() {
			g.Assert(t, "error_flows_missing_coverage_file_windows", []byte(mm.View()))
		} else {
			g.Assert(t, "error_flows_missing_coverage_file", []byte(mm.View()))
		}

		mm, cmd = mt.sendLetterKey('f')
		require.NotNil(t, mm)
		require.Equal(t, tea.Quit(), cmd())

		mm, cmd = mt.sendErrorMsg("")
		require.NotNil(t, mm)
		require.Nil(t, cmd)
	})

	t.Run("missing go.mod file", func(t *testing.T) {
		mt := &modelTest{
			T:               t,
			profileFilename: "coverage.out",
			codeRoot:        "testdata/no-go.mod",
		}
		initCmd := mt.init()
		initMsg := initCmd()

		mm, cmd := mt.sendWindowSizeMsg(60, 20)
		require.NotNil(t, mm)
		require.Nil(t, cmd)

		mm, cmd = mt.sendErrorMsg(initMsg)
		require.NotNil(t, mm)
		require.Nil(t, cmd)
		if isWindows() {
			g.Assert(t, "error_flows_missing_go.mod_file_windows", []byte(mm.View()))
		} else {
			g.Assert(t, "error_flows_missing_go.mod_file", []byte(mm.View()))
		}
	})

	t.Run("invalid coverage file", func(t *testing.T) {
		mt := &modelTest{
			T:               t,
			profileFilename: "invalid.profile",
			codeRoot:        "testdata/errors",
		}
		initCmd := mt.init()
		initMsg := initCmd()

		mm, cmd := mt.sendWindowSizeMsg(60, 20)
		require.NotNil(t, mm)
		require.Nil(t, cmd)

		mm, cmd = mt.sendErrorMsg(initMsg)
		require.NotNil(t, mm)
		require.Nil(t, cmd)
		g.Assert(t, "error_flows_invalid_coverage_file", []byte(mm.View()))
	})

	t.Run("no profiles", func(t *testing.T) {
		mt := &modelTest{
			T:               t,
			profileFilename: "empty.profile",
			codeRoot:        "testdata/errors",
		}
		initCmd := mt.init()
		initMsg := initCmd()

		mm, cmd := mt.sendWindowSizeMsg(60, 20)
		require.NotNil(t, mm)
		require.Nil(t, cmd)

		mm, cmd = mt.sendErrorMsg(initMsg)
		require.NotNil(t, mm)
		require.Nil(t, cmd)
		g.Assert(t, "error_flows_empty_coverage_file", []byte(mm.View()))
	})

	t.Run("invalid source file name", func(t *testing.T) {
		mt := &modelTest{
			T:               t,
			profileFilename: "coverage.profile",
			codeRoot:        "testdata/errors",
		}
		initCmd := mt.init()
		initMsg := initCmd()

		mm, cmd := mt.sendWindowSizeMsg(60, 20)
		require.NotNil(t, mm)
		require.Nil(t, cmd)

		mm, cmd = mt.sendProfilesMsg(initMsg)
		require.NotNil(t, mm)
		require.Nil(t, cmd)

		mm, cmd = mt.sendEnterKey()
		require.NotNil(t, mm)
		require.NotNil(t, cmd)

		errMsg := cmd()
		require.NotNil(t, errMsg)

		mm, cmd = mt.sendErrorMsg(errMsg)
		require.NotNil(t, mm)
		require.Nil(t, cmd)
		if isWindows() {
			g.Assert(t, "error_flows_missing_source_file_windows", []byte(mm.View()))
		} else {
			g.Assert(t, "error_flows_missing_source_file", []byte(mm.View()))
		}
	})

	t.Run("invalid go.mod", func(t *testing.T) {
		mt := &modelTest{
			T:               t,
			profileFilename: "cover.profile",
			codeRoot:        "testdata/errors/badmodule",
		}
		initCmd := mt.init()
		initMsg := initCmd()

		mm, cmd := mt.sendWindowSizeMsg(60, 20)
		require.NotNil(t, mm)
		require.Nil(t, cmd)

		mm, cmd = mt.sendErrorMsg(initMsg)
		require.NotNil(t, mm)
		require.Nil(t, cmd)
		g.Assert(t, "error_flows_invalid_go.mod", []byte(mm.View()))
	})
}

package gocovshtest

import (
	"path"
	"testing"

	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/require"
)

func TestHappyFlow(t *testing.T) {
	t.Run("no requested files", func(t *testing.T) {
		testHappyFlow(t, "not-requested", nil, nil)
	})

	t.Run("requested files", func(t *testing.T) {
		testHappyFlow(t, "requested", []string{"covered.go"}, nil)
	})

	t.Run("filtered lines", func(t *testing.T) {
		testHappyFlow(
			t,
			"filtered",
			[]string{"covered.go"},
			map[string][]int{
				"covered.go": {4},
			},
		)
	})
}

func testHappyFlow(t *testing.T, prefix string, requestedFiles []string, filteredLines map[string][]int) {
	dir := path.Join("testdata", "general", prefix)
	g := goldie.New(t, goldie.WithFixtureDir(dir))

	mt := &modelTest{
		T:               t,
		profileFilename: "profile.cover",
		codeRoot:        "testdata/general",
		requestedFiles:  requestedFiles,
		filteredLines:   filteredLines,
	}

	t.Run("initial setup", func(t *testing.T) {
		initCmd := mt.init()
		initMsg := initCmd()

		mm, cmd := mt.sendWindowSizeMsg(60, 20)
		require.NotNil(t, mm)
		require.Nil(t, cmd)

		mm, cmd = mt.sendProfilesMsg(initMsg)
		require.NotNil(t, mm)
		require.Nil(t, cmd)

		g.Assert(t, "happy_flow_initial_setup", []byte(mm.View()))
	})

	t.Run("first file", func(t *testing.T) {
		mm, cmd := mt.sendEnterKey()
		require.NotNil(t, mm)
		require.NotNil(t, cmd)

		// load file from the returned command
		fileMsg := cmd()
		require.NotNil(t, fileMsg)

		mm, cmd = mt.sendFileContentsMsg(fileMsg)
		require.NotNil(t, mm)
		require.Nil(t, cmd)

		g.Assert(t, "happy_flow_first_file", []byte(mm.View()))
	})

	t.Run("back to list", func(t *testing.T) {
		mm, cmd := mt.sendEscKey()
		require.NotNil(t, mm)
		require.Nil(t, cmd)

		g.Assert(t, "happy_flow_back_to_list", []byte(mm.View()))
	})

	t.Run("no exit on esc", func(t *testing.T) {
		mm, cmd := mt.sendEscKey()
		require.NotNil(t, mm)
		require.Nil(t, cmd)

		g.Assert(t, "happy_flow_no_exit_on_esc", []byte(mm.View()))
	})

	t.Run("select second file", func(t *testing.T) {
		mm, cmd := mt.sendLetterKey('j')
		require.NotNil(t, mm)
		require.NotNil(t, cmd) // command is not nit but irrelevant

		g.Assert(t, "happy_flow_select_second_file", []byte(mm.View()))
	})

	t.Run("view second file", func(t *testing.T) {
		mm, cmd := mt.sendEnterKey()
		require.NotNil(t, mm)
		require.NotNil(t, cmd)

		// load file from the returned command
		fileMsg := cmd()
		require.NotNil(t, fileMsg)

		mm, cmd = mt.sendFileContentsMsg(fileMsg)
		require.NotNil(t, mm)
		require.Nil(t, cmd)

		g.Assert(t, "happy_flow_view_second_file", []byte(mm.View()))
	})

	t.Run("codeview navigation", func(t *testing.T) {
		t.Run("bottom", func(t *testing.T) {
			mm, cmd := mt.sendLetterKey('G')
			require.NotNil(t, mm)
			require.Nil(t, cmd)

			g.Assert(t, "happy_flow_codeview_navigation_bottom", []byte(mm.View()))
		})

		t.Run("top", func(t *testing.T) {
			mm, cmd := mt.sendLetterKey('g')
			require.NotNil(t, mm)
			require.Nil(t, cmd)

			g.Assert(t, "happy_flow_codeview_navigation_top", []byte(mm.View()))
		})

		t.Run("back", func(t *testing.T) {
			mm, cmd := mt.sendEscKey()
			require.NotNil(t, mm)
			require.Nil(t, cmd)

			g.Assert(t, "happy_flow_codeview_navigation_back", []byte(mm.View()))
		})
	})

	t.Run("toggle help", func(t *testing.T) {
		t.Run("full", func(t *testing.T) {
			mm, cmd := mt.sendLetterKey('?')
			require.NotNil(t, mm)
			require.Nil(t, cmd)

			g.Assert(t, "happy_flow_toggle_help_full", []byte(mm.View()))
		})

		t.Run("none", func(t *testing.T) {
			mm, cmd := mt.sendLetterKey('?')
			require.NotNil(t, mm)
			require.Nil(t, cmd)

			g.Assert(t, "happy_flow_toggle_help_none", []byte(mm.View()))
		})

		t.Run("short", func(t *testing.T) {
			mm, cmd := mt.sendLetterKey('?')
			require.NotNil(t, mm)
			require.Nil(t, cmd)

			g.Assert(t, "happy_flow_toggle_help_short", []byte(mm.View()))
		})
	})
}

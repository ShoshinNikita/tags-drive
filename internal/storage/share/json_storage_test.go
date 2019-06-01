package share

import (
	"os"
	"testing"

	clog "github.com/ShoshinNikita/log/v2"
	"github.com/stretchr/testify/assert"

	"github.com/tags-drive/core/internal/storage/files"
)

const testJsonFile = "test.json"

func TestMain(m *testing.M) {
	code := m.Run()

	os.Remove(testJsonFile)

	os.Exit(code)
}

// Test of filesIDs type

func TestHasID(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		arr  []int
		find int
		res  bool
	}{
		{
			arr:  []int{1, 2, 3, 4},
			find: 5,
			res:  false,
		},
		{
			arr:  []int{1, 2, 3, 4},
			find: 1,
			res:  true,
		},
		{
			arr:  []int{1, 2, 3, 4},
			find: -4,
			res:  false,
		},
	}

	for i, tt := range tests {
		ids := newFileIDs(tt.arr)
		res := ids.hasID(tt.find)

		assert.Equalf(tt.res, res, "iteration #%d", i+1)
	}
}

func TestDeleteID(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		arr    []int
		delete int
		res    []int
	}{
		{
			arr:    []int{1, 2, 3, 4},
			delete: 5,
			res:    []int{1, 2, 3, 4},
		},
		{
			arr:    []int{1, 2, 3, 4},
			delete: 1,
			res:    []int{2, 3, 4},
		},
		{
			arr:    []int{1, 2, 3, 4},
			delete: 3,
			res:    []int{1, 2, 4},
		},
		{
			arr:    []int{1, 2, 3, 4},
			delete: 2,
			res:    []int{1, 3, 4},
		},
		{
			arr:    []int{1, 2, 3, 4},
			delete: 4,
			res:    []int{1, 2, 3},
		},
	}

	for i, tt := range tests {
		ids := newFileIDs(tt.arr)
		ids.deleteID(tt.delete)

		assert.Equalf(tt.res, []int(ids), "iteration #%d", i+1)
	}
}

// Test of jsonShareStorage type

func TestCheckFile(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		startTokens map[string]filesIDs
		checkToken  string
		checkID     int
		result      bool
	}{
		{
			startTokens: map[string]filesIDs{
				"1": []int{1, 2, 3, 4},
				"2": []int{2, 3, 4, 5},
				"3": []int{2, 6},
				"4": []int{1, 10, 23, 24, 15},
			},
			checkToken: "1",
			checkID:    1,
			result:     true,
		},
		{
			startTokens: map[string]filesIDs{
				"1": []int{1, 2, 3, 4},
				"2": []int{2, 3, 4, 5},
				"3": []int{2, 6},
				"4": []int{1, 10, 23, 24, 15},
			},
			checkToken: "4",
			checkID:    23,
			result:     true,
		},
		{
			startTokens: map[string]filesIDs{
				"1": []int{1, 2, 3, 4},
				"2": []int{2, 3, 4, 5},
				"3": []int{2, 6},
				"4": []int{1, 10, 23, 24, 15},
			},
			checkToken: "4",
			checkID:    100,
			result:     false,
		},
		{
			startTokens: map[string]filesIDs{
				"1": []int{1, 2, 3, 4},
				"2": []int{2, 3, 4, 5},
				"3": []int{2, 6},
				"4": []int{1, 10, 23, 24, 15},
			},
			checkToken: "5",
			checkID:    1,
			result:     false,
		},
	}

	for i, tt := range tests {
		st := newStorage()
		st.tokens = tt.startTokens

		res := st.CheckFile(tt.checkToken, tt.checkID)

		assert.Equalf(tt.result, res, "iteration #%d", i+1)

		assert.NoError(st.Shutdown())
	}
}

func TestDeleteFile(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		startIDs [][]int
		deleteID int
		result   [][]int
	}{
		{
			startIDs: [][]int{
				[]int{1, 2, 3, 4},
				[]int{2, 3, 4, 5},
				[]int{2, 6},
				[]int{1, 10, 23, 24, 15},
			},
			deleteID: 1,
			result: [][]int{
				[]int{2, 3, 4},
				[]int{2, 3, 4, 5},
				[]int{2, 6},
				[]int{10, 23, 24, 15},
			},
		},
		{
			startIDs: [][]int{
				[]int{1, 2, 3, 4},
				[]int{2, 3, 4, 5},
				[]int{2, 6},
				[]int{1, 10, 23, 24, 15},
			},
			deleteID: 10,
			result: [][]int{
				[]int{1, 2, 3, 4},
				[]int{2, 3, 4, 5},
				[]int{2, 6},
				[]int{1, 23, 24, 15},
			},
		},
		{
			startIDs: [][]int{
				[]int{1, 2, 3, 4},
				[]int{2, 3, 4, 5},
				[]int{2, 6},
				[]int{1, 10, 23, 24, 15},
			},
			deleteID: 100,
			result: [][]int{
				[]int{1, 2, 3, 4},
				[]int{2, 3, 4, 5},
				[]int{2, 6},
				[]int{1, 10, 23, 24, 15},
			},
		},
	}

	for i, tt := range tests {
		st := newStorage()
		var tokens []string
		for j := range tt.startIDs {
			t := st.CreateToken(tt.startIDs[j])
			tokens = append(tokens, t)
		}

		st.DeleteFile(tt.deleteID)

		var res [][]int
		for _, token := range tokens {
			res = append(res, st.tokens[token])
		}

		assert.Equalf(tt.result, res, "iteration #%d", i+1)

		assert.NoError(st.Shutdown())
	}
}

func TestFilterFiles(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		tokens map[string]filesIDs

		filterToken string
		files       []files.File

		isError bool
		res     []files.File
	}{
		{
			tokens: map[string]filesIDs{
				"1": []int{1, 2, 3},
			},
			filterToken: "1",
			//
			isError: false,
			files: []files.File{
				{ID: 1, Filename: "test"}, {ID: 2, Filename: "xyz"}, {ID: 3, Filename: "ppp"},
				{ID: 4, Filename: "asd"}, {ID: 5, Filename: "ghj"}, {ID: 6, Filename: "io"},
				{ID: 7, Filename: "git"}, {ID: 8, Filename: "cvb"}, {ID: 9, Filename: "txt"},
			},
			res: []files.File{
				{ID: 1, Filename: "test"}, {ID: 2, Filename: "xyz"}, {ID: 3, Filename: "ppp"},
			},
		},
		{
			tokens: map[string]filesIDs{
				"1": []int{1, 3, 9},
			},
			filterToken: "1",
			//
			isError: false,
			files: []files.File{
				{ID: 1, Filename: "test"}, {ID: 2, Filename: "xyz"}, {ID: 3, Filename: "ppp"},
				{ID: 4, Filename: "asd"}, {ID: 5, Filename: "ghj"}, {ID: 6, Filename: "io"},
				{ID: 7, Filename: "git"}, {ID: 8, Filename: "cvb"}, {ID: 9, Filename: "txt"},
			},
			res: []files.File{
				{ID: 1, Filename: "test"}, {ID: 3, Filename: "ppp"}, {ID: 9, Filename: "txt"},
			},
		},
		{
			tokens: map[string]filesIDs{
				"1": []int{1, 3, 8, 9, 10, 11, 12, 13, 14, 15, 16},
			},
			filterToken: "1",
			//
			isError: false,
			files: []files.File{
				{ID: 1, Filename: "test"}, {ID: 2, Filename: "xyz"}, {ID: 3, Filename: "ppp"},
				{ID: 4, Filename: "asd"}, {ID: 5, Filename: "ghj"}, {ID: 6, Filename: "io"},
				{ID: 7, Filename: "git"}, {ID: 8, Filename: "cvb"}, {ID: 9, Filename: "txt"},
			},
			res: []files.File{
				{ID: 1, Filename: "test"}, {ID: 3, Filename: "ppp"}, {ID: 8, Filename: "cvb"},
				{ID: 9, Filename: "txt"},
			},
		},
		{
			tokens: map[string]filesIDs{
				"1": []int{1, 3, 9},
				"2": []int{20, 30, 31, 32},
			},
			filterToken: "2",
			//
			isError: false,
			files: []files.File{
				{ID: 1, Filename: "test"}, {ID: 2, Filename: "xyz"}, {ID: 3, Filename: "ppp"},
				{ID: 4, Filename: "asd"}, {ID: 5, Filename: "ghj"}, {ID: 6, Filename: "io"},
				{ID: 7, Filename: "git"}, {ID: 8, Filename: "cvb"}, {ID: 9, Filename: "txt"},
			},
			res: []files.File{},
		},
		{
			tokens: map[string]filesIDs{
				"1": []int{1, 2, 3},
			},
			filterToken: "5",
			//
			isError: true,
			files: []files.File{
				{ID: 1, Filename: "test"}, {ID: 2, Filename: "xyz"}, {ID: 3, Filename: "ppp"},
				{ID: 4, Filename: "asd"}, {ID: 5, Filename: "ghj"}, {ID: 6, Filename: "io"},
				{ID: 7, Filename: "git"}, {ID: 8, Filename: "cvb"}, {ID: 9, Filename: "txt"},
			},
			res: []files.File{},
		},
	}

	for i, tt := range tests {
		st := newStorage()
		st.tokens = tt.tokens

		res, err := st.FilterFiles(tt.filterToken, tt.files)

		assert.Equal(err != nil, tt.isError, "iteration #%d", i+1)

		if !tt.isError {
			assert.Equal(tt.res, res, "iteration #%d", i+1)
		}

		assert.NoError(st.Shutdown())
	}
}

func newStorage() *jsonShareStorage {
	st := newJsonShareStorage(Config{
		ShareTokenJSONFile: testJsonFile,
		Encrypt:            false,
	}, clog.NewProdLogger())

	err := st.init()
	if err != nil {
		panic(err)
	}

	return st
}
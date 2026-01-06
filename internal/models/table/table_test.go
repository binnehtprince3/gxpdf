package table

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTable(t *testing.T) {
	tests := []struct {
		name      string
		rowCount  int
		colCount  int
		expectErr bool
	}{
		{"valid 3x4", 3, 4, false},
		{"valid 1x1", 1, 1, false},
		{"valid 10x5", 10, 5, false},
		{"invalid rows", 0, 3, true},
		{"invalid cols", 3, 0, true},
		{"negative rows", -1, 3, true},
		{"negative cols", 3, -1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			table, err := NewTable(tt.rowCount, tt.colCount)

			if tt.expectErr {
				assert.Error(t, err)
				assert.Nil(t, table)
			} else {
				require.NoError(t, err)
				require.NotNil(t, table)
				assert.Equal(t, tt.rowCount, table.RowCount)
				assert.Equal(t, tt.colCount, table.ColCount)
				assert.Len(t, table.Rows, tt.rowCount)

				// Check all cells are initialized
				for r := 0; r < tt.rowCount; r++ {
					assert.Len(t, table.Rows[r], tt.colCount)
					for c := 0; c < tt.colCount; c++ {
						cell := table.Rows[r][c]
						assert.NotNil(t, cell)
						assert.Equal(t, r, cell.Row)
						assert.Equal(t, c, cell.Column)
						assert.Equal(t, "", cell.Text)
					}
				}
			}
		})
	}
}

func TestTable_GetCell(t *testing.T) {
	table, err := NewTable(3, 4)
	require.NoError(t, err)

	// Valid positions
	cell := table.GetCell(0, 0)
	assert.NotNil(t, cell)
	assert.Equal(t, 0, cell.Row)
	assert.Equal(t, 0, cell.Column)

	cell = table.GetCell(2, 3)
	assert.NotNil(t, cell)
	assert.Equal(t, 2, cell.Row)
	assert.Equal(t, 3, cell.Column)

	// Invalid positions
	assert.Nil(t, table.GetCell(-1, 0))
	assert.Nil(t, table.GetCell(0, -1))
	assert.Nil(t, table.GetCell(3, 0)) // Out of bounds
	assert.Nil(t, table.GetCell(0, 4)) // Out of bounds
}

func TestTable_SetCell(t *testing.T) {
	table, err := NewTable(3, 3)
	require.NoError(t, err)

	// Valid set
	cell := NewCell("Test", 1, 2)
	err = table.SetCell(1, 2, cell)
	assert.NoError(t, err)
	assert.Equal(t, "Test", table.GetCell(1, 2).Text)

	// Cell position should be updated to match table position
	wrongCell := NewCell("Wrong", 99, 99)
	err = table.SetCell(0, 1, wrongCell)
	assert.NoError(t, err)
	retrieved := table.GetCell(0, 1)
	assert.Equal(t, 0, retrieved.Row)
	assert.Equal(t, 1, retrieved.Column)

	// Invalid positions
	err = table.SetCell(-1, 0, cell)
	assert.Error(t, err)

	err = table.SetCell(0, -1, cell)
	assert.Error(t, err)

	err = table.SetCell(3, 0, cell)
	assert.Error(t, err)

	err = table.SetCell(0, 3, cell)
	assert.Error(t, err)
}

func TestTable_GetRow(t *testing.T) {
	table, err := NewTable(3, 4)
	require.NoError(t, err)

	// Set some data
	table.SetCell(1, 0, NewCell("A", 1, 0))
	table.SetCell(1, 1, NewCell("B", 1, 1))
	table.SetCell(1, 2, NewCell("C", 1, 2))
	table.SetCell(1, 3, NewCell("D", 1, 3))

	// Get row
	row := table.GetRow(1)
	require.Len(t, row, 4)
	assert.Equal(t, "A", row[0].Text)
	assert.Equal(t, "B", row[1].Text)
	assert.Equal(t, "C", row[2].Text)
	assert.Equal(t, "D", row[3].Text)

	// Invalid rows
	assert.Nil(t, table.GetRow(-1))
	assert.Nil(t, table.GetRow(3))
}

func TestTable_GetColumn(t *testing.T) {
	table, err := NewTable(3, 4)
	require.NoError(t, err)

	// Set some data
	table.SetCell(0, 2, NewCell("X", 0, 2))
	table.SetCell(1, 2, NewCell("Y", 1, 2))
	table.SetCell(2, 2, NewCell("Z", 2, 2))

	// Get column
	col := table.GetColumn(2)
	require.Len(t, col, 3)
	assert.Equal(t, "X", col[0].Text)
	assert.Equal(t, "Y", col[1].Text)
	assert.Equal(t, "Z", col[2].Text)

	// Invalid columns
	assert.Nil(t, table.GetColumn(-1))
	assert.Nil(t, table.GetColumn(4))
}

func TestTable_IsEmpty(t *testing.T) {
	table, err := NewTable(2, 2)
	require.NoError(t, err)

	// Initially empty
	assert.True(t, table.IsEmpty())

	// Add content
	table.SetCell(0, 0, NewCell("Test", 0, 0))
	assert.False(t, table.IsEmpty())
}

func TestTable_CellCount(t *testing.T) {
	tests := []struct {
		name     string
		rows     int
		cols     int
		expected int
	}{
		{"3x4", 3, 4, 12},
		{"1x1", 1, 1, 1},
		{"5x2", 5, 2, 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			table, err := NewTable(tt.rows, tt.cols)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, table.CellCount())
		})
	}
}

func TestTable_NonEmptyCellCount(t *testing.T) {
	table, err := NewTable(3, 3)
	require.NoError(t, err)

	assert.Equal(t, 0, table.NonEmptyCellCount())

	// Add some content
	table.SetCell(0, 0, NewCell("A", 0, 0))
	assert.Equal(t, 1, table.NonEmptyCellCount())

	table.SetCell(1, 1, NewCell("B", 1, 1))
	assert.Equal(t, 2, table.NonEmptyCellCount())

	table.SetCell(2, 2, NewCell("C", 2, 2))
	assert.Equal(t, 3, table.NonEmptyCellCount())
}

func TestTable_HasMergedCells(t *testing.T) {
	table, err := NewTable(3, 3)
	require.NoError(t, err)

	assert.False(t, table.HasMergedCells())

	// Add merged cell
	mergedCell := NewCell("Merged", 0, 0).WithRowSpan(2)
	table.SetCell(0, 0, mergedCell)
	assert.True(t, table.HasMergedCells())
}

func TestTable_ToStringGrid(t *testing.T) {
	table, err := NewTable(2, 3)
	require.NoError(t, err)

	// Populate table
	table.SetCell(0, 0, NewCell("A1", 0, 0))
	table.SetCell(0, 1, NewCell("B1", 0, 1))
	table.SetCell(0, 2, NewCell("C1", 0, 2))
	table.SetCell(1, 0, NewCell("A2", 1, 0))
	table.SetCell(1, 1, NewCell("B2", 1, 1))
	table.SetCell(1, 2, NewCell("C2", 1, 2))

	grid := table.ToStringGrid()
	require.Len(t, grid, 2)
	require.Len(t, grid[0], 3)
	require.Len(t, grid[1], 3)

	assert.Equal(t, "A1", grid[0][0])
	assert.Equal(t, "B1", grid[0][1])
	assert.Equal(t, "C1", grid[0][2])
	assert.Equal(t, "A2", grid[1][0])
	assert.Equal(t, "B2", grid[1][1])
	assert.Equal(t, "C2", grid[1][2])
}

func TestTable_Validate(t *testing.T) {
	// Valid table
	table, err := NewTable(3, 4)
	require.NoError(t, err)
	assert.NoError(t, table.Validate())

	// Invalid row count
	invalidTable := &Table{
		RowCount: 0,
		ColCount: 3,
		Rows:     [][]*Cell{},
	}
	assert.Error(t, invalidTable.Validate())

	// Invalid column count
	invalidTable = &Table{
		RowCount: 3,
		ColCount: 0,
		Rows:     make([][]*Cell, 3),
	}
	assert.Error(t, invalidTable.Validate())

	// Row count mismatch
	invalidTable = &Table{
		RowCount: 3,
		ColCount: 3,
		Rows:     make([][]*Cell, 2), // Should be 3
	}
	assert.Error(t, invalidTable.Validate())

	// Column count mismatch
	table, _ = NewTable(2, 3)
	table.Rows[0] = make([]*Cell, 2) // Should be 3
	assert.Error(t, table.Validate())

	// Wrong cell indices
	table, _ = NewTable(2, 2)
	wrongCell := NewCell("Wrong", 99, 99)
	table.Rows[0][0] = wrongCell
	assert.Error(t, table.Validate())
}

func TestTable_String(t *testing.T) {
	table, err := NewTable(2, 2)
	require.NoError(t, err)
	table.Method = "Lattice"
	table.PageNum = 5

	table.SetCell(0, 0, NewCell("A", 0, 0))
	table.SetCell(0, 1, NewCell("B", 0, 1))

	str := table.String()
	assert.Contains(t, str, "rows=2")
	assert.Contains(t, str, "cols=2")
	assert.Contains(t, str, "method=Lattice")
	assert.Contains(t, str, "page=5")
	assert.Contains(t, str, "\"A\"")
	assert.Contains(t, str, "\"B\"")
}

func TestTable_WithMetadata(t *testing.T) {
	table, err := NewTable(3, 3)
	require.NoError(t, err)

	// Set metadata
	table.PageNum = 10
	table.Method = "Stream"
	table.Bounds = NewRectangle(100, 200, 400, 300)

	assert.Equal(t, 10, table.PageNum)
	assert.Equal(t, "Stream", table.Method)
	assert.Equal(t, 100.0, table.Bounds.X)
	assert.Equal(t, 200.0, table.Bounds.Y)
	assert.Equal(t, 400.0, table.Bounds.Width)
	assert.Equal(t, 300.0, table.Bounds.Height)
}

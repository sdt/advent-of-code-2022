package aoc

type Grid[T any] struct {
    w, h int
    cell []T
}

func NewGrid[T any](w, h int, defaultValue T) *Grid[T] {
    grid := Grid[T]{
        w: w,
        h: h,
        cell: make([]T, w * h),
    }
    for i := range grid.cell {
        grid.cell[i] = defaultValue
    }
    return &grid
}

func (this Grid[T]) Get(x, y int) T {
    return this.cell[this.offset(x, y)]
}

func (this *Grid[T]) Set(x, y int, value T) {
    this.cell[this.offset(x, y)] = value
}

func (this Grid[T]) GetMaybe(x, y int) (T, bool) {
    if this.Contains(x, y) {
        return this.cell[this.offset(x, y)], true
    }
    var nothing T
    return nothing, false
}

func (this Grid[T]) Width() int {
    return this.w
}

func (this Grid[T]) Height() int {
    return this.h
}

func (this Grid[T]) offset(x, y int) int {
    return this.w * y + x
}

func (this Grid[T]) Contains(x, y int) bool {
    return x >= 0 && y >= 0 && x < this.w && y < this.h
}

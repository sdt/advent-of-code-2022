package main

import (
    "advent-of-code/aoc"
    "fmt"
    "strings"
)

type File struct {
    name string
    size int
}

type Directory struct {
    name string
    files []File
    directories map[string]*Directory
    parent *Directory
    size int
}

func main() {
    filename := aoc.GetFilename()
    lines := aoc.GetInputLines(filename)

    fmt.Println(part1(lines))
}

func part1(lines []string) int {
    root := ParseInput(lines)
    ComputeSize(root)
    return SumSmallDirs(root)
}

func SumSmallDirs(dir *Directory) int {
    total := 0
    for _, child := range dir.directories {
        total += SumSmallDirs(child)
    }
    if dir.size <= 100000 {
        total += dir.size
    }
    return total
}

func ComputeSize(dir *Directory) int {
    for _, file := range dir.files {
        dir.size += file.size
    }
    for _, child := range dir.directories {
        dir.size += ComputeSize(child)
    }
    return dir.size
}

func ParseInput(lines[] string) *Directory {
    cwd := NewDirectory("/", nil)

    for _, line := range lines {
        cwd = ParseLine(cwd, line)
    }

    // Walk back up to the root directory
    for cwd.parent != nil {
        cwd = cwd.parent
    }
    return cwd
}

func ParseLine(cwd *Directory, line string) *Directory {
    if line == "$ ls" || line == "$ cd /" {
        return cwd
    }

    if line == "$ cd .." {
        return cwd.parent
    }

    if strings.HasPrefix(line, "$ cd ") {
        start := len("$ cd ")
        name := line[start:]
        child := cwd.directories[name]
        return child
    }

    if strings.HasPrefix(line, "dir ") {
        start := len("dir ")
        name := line[start:]
        cwd.directories[name] = NewDirectory(name, cwd)
        return cwd
    }

    // File
    words := strings.Split(line, " ")
    cwd.files = append(cwd.files, NewFile(words[1], aoc.ParseInt(words[0])))

    return cwd
}

func NewDirectory(name string, parent *Directory) *Directory {
    return &Directory{
        name: name,
        parent: parent,
        directories: make(map[string]*Directory),
    }
}

func NewFile(name string, size int) File {
    return File{ name: name, size: size }
}

package storage

import (
	"os"

	"github.com/Anacardo89/kanboards/fsops"
	"github.com/Anacardo89/kanboards/logger"
	"gopkg.in/yaml.v2"
)

type Menu struct {
	Projects []Project `yaml:"projects"`
}

type Project struct {
	Id     int64   `yaml:"id"`
	Title  string  `yaml:"title"`
	Labels []Label `yaml:"labels"`
	Boards []Board `yaml:"boards"`
}

type Board struct {
	Id    int64  `yaml:"id"`
	Pos   int    `yaml:"position"`
	Title string `yaml:"title"`
	Cards []Card `yaml:"cards"`
}

type Label struct {
	Id    int64  `yaml:"id"`
	Title string `yaml:"title"`
	Color string `yaml:"color"`
}

type Card struct {
	Id          int64       `yaml:"id"`
	Title       string      `yaml:"title"`
	Description string      `yaml:"description"`
	CheckList   []CheckItem `yaml:"checklist"`
	CardLabels  []Label     `yaml:"card_labels"`
}

type CheckItem struct {
	Id    int64  `yaml:"id"`
	Title string `yaml:"title"`
	Check bool   `yaml:"check"`
}

func (m *Menu) ToYAML() string {
	data, err := yaml.Marshal(m)
	if err != nil {
		logger.Error.Println("Cannot export to YAML", err)
	}
	datastr := string(data)
	return datastr
}

func FromYAML(data []byte) *Menu {
	m := Menu{}
	err := yaml.Unmarshal(data, &m)
	if err != nil {
		logger.Error.Println("Cannot import from YAML", err)
	}
	return &m
}

func ToFile(data string) {
	f, err := os.Open(fsops.YamlPath)
	if err == nil {
		os.Remove(fsops.YamlPath)
	} else {
		logger.Error.Println(err)
	}
	f.Close()
	f, err = os.Create(fsops.YamlPath)
	if err != nil {
		logger.Error.Println(err)
	}
	defer f.Close()
	f.WriteString(data)
}

func FromFile() []byte {
	data, err := os.ReadFile(fsops.YamlPath)
	if err != nil {
		logger.Error.Println(err)
	}
	return data
}

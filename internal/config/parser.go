package config

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/gabefiori/gsp/internal/finder"
)

type Parser struct {
	line int
	sc   *bufio.Scanner
	cfg  *Config
}

func NewParser(r io.Reader, cfg *Config) *Parser {
	return &Parser{
		line: 1,
		sc:   bufio.NewScanner(r),
		cfg:  cfg,
	}
}

func (p *Parser) Run() error {
	for ; p.sc.Scan(); p.line++ {
		line := strings.TrimSpace(p.sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		eq := strings.IndexByte(line, '=')
		if eq == -1 {
			return p.lineErr("invalid field")
		}

		key := strings.TrimSpace(line[:eq])
		val := strings.TrimSpace(line[eq+1:])

		if err := p.field(key, val); err != nil {
			return err
		}
	}

	return nil
}

func (p *Parser) field(k, v string) error {
	switch k {
	case "selector":
		p.cfg.Selector = v
	case "sort":
		p.cfg.Sort = v
	case "expand-output":
		p.cfg.ExpandOutput = v == "true"
	case "unique":
		p.cfg.Unique = v == "true"
	case "source":
		sep := strings.IndexByte(v, ':')
		if sep == -1 {
			return p.lineErr("invalid source")
		}

		path := v[sep+1:]
		depth, err := strconv.ParseUint(v[:sep], 10, 8)
		if err != nil {
			return p.lineErr(err.Error())
		}

		p.cfg.Sources = append(p.cfg.Sources, finder.Source{
			Depth:        uint8(depth),
			OriginalPath: path,
		})
	}

	return nil
}

func (p *Parser) lineErr(msg string) error {
	return fmt.Errorf("failed to parse config %q on line %d.", msg, p.line)
}

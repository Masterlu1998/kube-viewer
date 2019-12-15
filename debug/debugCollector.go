package debug

import (
	"fmt"
	"time"
)

type Message struct {
	where   string
	level   messageLevel
	content string
	t       time.Time
}

type messageLevel string

const (
	Info  messageLevel = "Info"
	Error messageLevel = "Error"
	Warn  messageLevel = "Warn"
)

func NewDebugMessage(level messageLevel, content, where string) Message {
	return Message{
		where:   where,
		level:   level,
		content: content,
		t:       time.Now(),
	}
}

func (d Message) Format() string {
	switch d.level {
	case Info:
		return fmt.Sprintf("[[%s]](fg:green,mod:bold) %s [[%s]](fg:blue): %s",
			d.level,
			d.t.Format(time.RFC3339),
			d.where,
			d.content,
		)
	case Warn:
		return fmt.Sprintf("[[%s]](fg:yellow,mod:bold) %s [[%s]](fg:blue): %s",
			d.level,
			d.t.Format(time.RFC3339),
			d.where,
			d.content,
		)
	case Error:
		return fmt.Sprintf("[[%s]](fg:red,mod:bold) %s [[%s]](fg:blue): %s",
			d.level,
			d.t.Format(time.RFC3339),
			d.where,
			d.content,
		)
	}

	return ""
}

func NewDebugCollector() *DebugCollector {
	return &DebugCollector{
		debugCh: make(chan Message),
	}
}

type DebugCollector struct {
	debugCh chan Message
}

func (e *DebugCollector) GetDebugMessageCh() chan Message {
	return e.debugCh
}

func (e *DebugCollector) Collect(m Message) {
	e.debugCh <- m
}

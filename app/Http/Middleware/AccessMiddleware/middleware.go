package accessmiddleware

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
)

type AccessLogger struct {
	mu   sync.Mutex
	dir  string
	date string
	file *os.File
}

func New() fiber.Handler {
	l := &AccessLogger{dir: filepath.FromSlash("storages/logs")}
	_ = os.MkdirAll(l.dir, 0o755)

	runtime.SetFinalizer(l, func(al *AccessLogger) {
		al.mu.Lock()
		defer al.mu.Unlock()
		if al.file != nil {
			_ = al.file.Close()
		}
	})

	return func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		end := time.Now()

		l.write(end, formatLine(c, start, end))
		return err
	}
}

func (l *AccessLogger) write(now time.Time, line string) {
	date := now.Format("060102") // YYMMDD

	l.mu.Lock()
	defer l.mu.Unlock()

	if l.file == nil || l.date != date {
		if l.file != nil {
			_ = l.file.Close()
		}

		path := filepath.Join(l.dir, fmt.Sprintf("access-%s.log", date))
		f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
		if err != nil {
			fmt.Fprintln(os.Stderr, "access log open failed:", err)
			l.file = nil
			l.date = ""
			return
		}

		l.file = f
		l.date = date
	}

	if _, err := l.file.WriteString(line); err != nil {
		fmt.Fprintln(os.Stderr, "access log write failed:", err)
	}
}

func formatLine(c *fiber.Ctx, start, end time.Time) string {
	method := sanitize(string(c.Method()))
	path := sanitize(string(c.Path()))
	ip := sanitize(c.IP())
	ua := sanitize(string(c.Get("User-Agent")))
	status := c.Response().StatusCode()
	latencyMs := end.Sub(start).Milliseconds()
	ts := end.Format(time.RFC3339Nano)

	return fmt.Sprintf("%s %s %s %d %dms %s %q\n", ts, method, path, status, latencyMs, ip, ua)
}

func sanitize(s string) string {
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.ReplaceAll(s, "\r", " ")
	return strings.TrimSpace(s)
}

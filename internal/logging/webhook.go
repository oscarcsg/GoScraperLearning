// internal/logging/webhook.go
package logging

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync/atomic"
	"time"

	"go.uber.org/zap/zapcore"
)

// Configuration from env
var externalConfig ExternalLogConfig

// 100 msgs buffer
var alertChannel = make(chan alertPayload, 100)

// Counter atomically used
var droppedLogs atomic.Uint64

func startAlertWorker(config ExternalLogConfig) {
	// If no external logs provider, just do nothing
	if config.Provider == "none" || config.Provider == "" {
		return
	}

	externalConfig = config

	// Parallel
	go func() {
		// 5 seconds to group errors
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop() // defer executes this line at the end of the function

		var errorCount uint64

		for {
			select {
			case payload := <-alertChannel:
				errorCount++
				if errorCount <= 3 {
					sendAlert(payload)
				}

			case <-ticker.C:
				dropped := droppedLogs.Swap(0)

				if errorCount > 3 || dropped > 0 {
					omittedFromQueue := uint64(0)
					
					if errorCount > 3 {
						omittedFromQueue = errorCount - 3
					}
					totalOmitted := omittedFromQueue + dropped

					summaryMsg := fmt.Sprintf("FLOOD: %d additional errors omitted during the last 5 seconds to avoid an overload in the HTTP petitions. Go to the local log file to see every log.", totalOmitted)

					sendAlert(alertPayload{
						Level: FATAL,
						Message: summaryMsg,
						Fields: nil,
					})
				}
				errorCount = 0
			}
		}
	}()
}

func queueAlert(lvl Level, msg string, fields []Field) {
	if externalConfig.Provider == "none" || externalConfig.Provider == "" {
		return
	}

	select {
	case alertChannel <- alertPayload{Level: lvl, Message: msg, Fields: fields}:
		// Stored in the buffer
	default:
		// The buffer is full, add 1 directly to the atomic counter
		droppedLogs.Add(1)
	}
}

func sendAlert(payload alertPayload) {
	switch externalConfig.Provider {
	case "telegram":
		sendTelegram(payload)
	case "webhook":
		sendGenericWebhook(payload)
	}
}

// --- Auxiliar Conversors ---

func levelToString(lvl Level) string {
	switch lvl {
	case ALWAYS: return "ALWAYS"
	case DEBUG:  return "DEBUG"
	case INFO:   return "INFO"
	case WARN:   return "WARN"
	case ERROR:  return "ERROR"
	case FATAL:  return "FATAL"
	default:     return "UNKNOWN"
	}
}

func getEmoji(lvl Level) string {
	switch lvl {
	case ALWAYS: return "📢"
	case DEBUG:  return "🐞"
	case INFO:   return "ℹ️"
	case WARN:   return "⚠️"
	case ERROR:  return "❌"
	case FATAL:  return "🚨"
	default:     return "🐛"
	}
}

// --- HTTP Integrations ---

func sendTelegram(payload alertPayload) {
	lvlStr := levelToString(payload.Level)
	
	var sb strings.Builder
	fmt.Fprintf(
		&sb,
		"%s <b>[%s]</b> - %s\n\n%s",
		getEmoji(payload.Level),
		lvlStr,
		time.Now().UTC().Format("2006-01-02 15:04:05 MST"),
		payload.Message,
	)

	// Volcamos los fields de Zap a un mapa para visualizarlos en Telegram
	if len(payload.Fields) > 0 {
		enc := zapcore.NewMapObjectEncoder()
		for _, f := range payload.Fields {
			f.AddTo(enc)
		}
		
		sb.WriteString("\n\n<b><i>Details:</i></b>\n<code>")
		for k, v := range enc.Fields {
			fmt.Fprintf(
				&sb,
				"%s: %+v\n",
				k,
				v,
			)
		}
		sb.WriteString("</code>")
	}

	body, err := json.Marshal(map[string]string{
		"chat_id":    externalConfig.TelChatId,
		"text":       sb.String(),
		"parse_mode": "HTML",
	})
	// If failed creating the json, cancel
	if err != nil {
		return
	}

	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", externalConfig.TelBotToken)
	
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return 
	}
	req.Header.Set("Content-Type", "application/json")

	// 5 seconds of timeout to avoid infinite waiting
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	
	// If petition had no errors, close the TCP connection
	if err == nil {
		resp.Body.Close()
	}
}

func sendGenericWebhook(payload alertPayload) {
	enc := zapcore.NewMapObjectEncoder()
	for _, f := range payload.Fields {
		f.AddTo(enc)
	}

	body, err := json.Marshal(map[string]any{
		"level":   levelToString(payload.Level),
		"message": payload.Message,
		"time":    time.Now().UTC().Format(time.RFC3339),
		"fields":  enc.Fields,
	})
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", externalConfig.WebhookURL, bytes.NewBuffer(body))
	// Check the error. If it is not nil, it means req is no valid, so just don't use it and close the method
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	
	if externalConfig.WebhookAuthHeader != "" {
		req.Header.Set("Authorization", externalConfig.WebhookAuthHeader)
	}

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	
	// If petition had no errors, close the TCP connection
	if err == nil {
		resp.Body.Close()
	}
}

package main

import "log"

// OverrideLogger is a dummy logger to kill the gogm logger
type OverrideLogger struct {
	Level string
}

func (d OverrideLogger) Debug(s string) {
	if d.Level == "DEBUG" {
		log.Println("[DEBUG] " + s)
	}
}

func (d OverrideLogger) Debugf(s string, vals ...interface{}) {
	if d.Level == "DEBUG" {
		log.Printf("[DEBUG] "+s+"\n", vals...)
	}
}

func (d OverrideLogger) Info(s string) {
	log.Println("[INFO] " + s)
}

func (d OverrideLogger) Infof(s string, vals ...interface{}) {
	log.Printf("[INFO] "+s+"\n", vals...)
}

func (d OverrideLogger) Warn(s string) {
	log.Println("[WARN] " + s)
}

func (d OverrideLogger) Warnf(s string, vals ...interface{}) {
	log.Printf("[WARN] "+s+"\n", vals...)
}

func (d OverrideLogger) Error(s string) {
	log.Println("[ERROR] " + s)
}

func (d OverrideLogger) Errorf(s string, vals ...interface{}) {
	log.Printf("[ERROR] "+s+"\n", vals...)
}

func (d OverrideLogger) Fatal(s string) {
	log.Fatalln("[FATAL] " + s)
}

func (d OverrideLogger) Fatalf(s string, vals ...interface{}) {
	log.Fatalf("[FATAL] "+s+"\n", vals...)
}

package rnd

import (
	"fmt"
	"math/rand"
	"time"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = r_char()
	}
	return string(b)
}

// varying char
func r_char() byte {
	return letters[rand.Intn(len(letters))]
}

// varying char
func r_var_char(min, max int) string {
	if (max - min) == 0 {
		return RandString(max)
	}
	s := min + rand.Intn(max-min)
	return RandString(s)
}

// text
func r_text() string {
	return r_var_char(1024, 2048)
}

// SMALLINT, INTEGER, AND BIGINT
func r_int(max int) int {
	return rand.Intn(max)
}

// NUMERIC with varying precision
func r_numeric(max float64) float64 {
	return (0.1 + rand.Float64()) * (max - 0.1)
}

// timestamp etc.
func r_datetime(relative int8) time.Time {
	if relative == 0 {
		return time.Now()
	}

	mod := rand.Intn(999999)
	if relative > 0 {
		return time.Now().Add(time.Duration(mod) * time.Minute)
	}

	return time.Now().Truncate(time.Duration(mod) * time.Minute)
}

// interval       16bytes
func r_interval(relative int8) time.Duration {
	return 10 * time.Second
}

// https://www.postgresql.org/docs/9.3/static/datatype-binary.html
func r_byte_array(min, max int) []byte {
	out := r_var_char(min, max)
	return []byte(out)
}

func PSQL_char() string {
	r := r_char()
	return string(r)
}
func PSQL_var_char(min, max int) string {
	r := r_var_char(min, max)
	return fmt.Sprintf("'%s'", r)
}

func PSQL_text() string {
	r := r_text()
	return r
}

func PSQL_int(max int) int {
	return r_int(max)
}

func PSQL_numeric(max float64, places uint8) float64 {
	r := r_numeric(max)
	// TODO format places
	return r
}

func PSQL_datetime(rel int8, fm uint8) string {

	// 0 timestamp       8 bytes WITHOUT timezone
	// 1 timestamp       8 bytes WITH timezone
	// 2 date            4 bytes (no time of day)
	// 3 time            8 bytes WITHOUT timezone
	// 4 time           12 bytes WITH timezone

	r := r_datetime(rel)
	rs := ""
	switch fm {
	case 0:
		rs = r.Format("2006-01-02 15:04:05")
		break
	case 1:
		rs = r.Format("2006-01-02 15:04:05")
		break
	case 2:
		rs = r.Format("2006-01-02")
		break
	case 3:
		rs = r.Format("15:04:05")
		break
	case 4:
		rs = r.Format("15:04:05")
		break
	}

	return fmt.Sprintf("'%s'", rs)
}

func PSQL_interval() string {
	r := r_interval(1)
	return fmt.Sprintf("%s", r)
}

func PSQL_byte_array(min, max int) string {
	return string(r_byte_array(min, max))
}

func PSQL_bool() string {
	if (rand.Intn(1000) % 2) > 0 {
		return "true"
	}
	return "false"
}

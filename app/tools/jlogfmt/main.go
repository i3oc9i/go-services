// This program takes the structured log output and makes it readable.
package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
)

var service string

func init() {
	flag.StringVar(&service, "service", "", "filter which service to see")
}

func main() {
	flag.Parse()
	var b strings.Builder

	// Scan standard input for log data per line.
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		s := scanner.Text()

		// Convert the JSON to a map for processing.
		m := make(map[string]any)
		err := json.Unmarshal([]byte(s), &m)
		if err != nil {
			if service == "" {
				fmt.Println(s)
			}
			continue
		}

		// If a service filter was provided, check.
		if service != "" && m["service"] != service {
			continue
		}

		// I like always having a traceid present in the logs.
		traceID := "00000000-0000-0000-0000-000000000000"
		if v, ok := m["traceid"]; ok {
			traceID = fmt.Sprintf("%v", v)
		}

		// Build out the know portions of the log in the order I want them in.
		b.Reset()
		b.WriteString(fmt.Sprintf("%s: %s: %s: %s: %s: %s: ",
			m["service"],
			m["ts"],
			m["level"],
			traceID,
			m["caller"],
			m["msg"],
		))

		// if status key exist then we print just after.
		if val, ok := m["status"]; ok {
			b.WriteString(fmt.Sprintf("status[%v]: ", val))
		}

		// if statusCode key exist then we print jsut next.
		if val, ok := m["statuscode"]; ok {
			b.WriteString(fmt.Sprintf("statuscode[%v]: ", val))
		}

		// Sort the keys of the map into a list
		keys := make([]string, 0, len(m))

		for k := range m {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		// Add the rest of the keys ignoring the ones we already added.
		for _, k := range keys {
			switch k {
			case "service", "ts", "level", "traceid", "caller", "msg", "status", "statuscode":
				continue
			}

			b.WriteString(fmt.Sprintf("%s[%v]: ", k, m[k]))
		}

		// Write the new log format, removing the last ':'
		out := b.String()
		fmt.Println(out[:len(out)-2])
	}

	if err := scanner.Err(); err != nil {
		log.Println(err)
	}
}

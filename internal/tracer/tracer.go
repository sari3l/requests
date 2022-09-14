package tracer

import (
	"crypto/tls"
	"fmt"
	"net/http/httptrace"
	"time"
)

const (
	TraceGotConn int = iota
	TraceDNSStart
	TraceDNSDone
	TraceConnectStart
	TraceConnectDone
	TraceGetConn
	TraceGotFirstResponseByte
	TraceTLSHandshakeStart
	TraceTLSHandshakeDone
)

type tracerLog struct {
	tType   int
	time    time.Time
	message string
}

type Tracer struct {
	connInfo    *httptrace.GotConnInfo
	tracerLogs  map[int]*tracerLog
	Enable      bool
	StartTime   time.Time
	EndTime     time.Time
	ClientTrace *httptrace.ClientTrace
	Output      struct {
		Prefix            string
		remoteAddr        string
		timeDNSLookup     time.Duration
		timeTCPConnect    time.Duration
		timeTLSHandshake  time.Duration
		timeConnect       time.Duration
		timeFirstResponse time.Duration
		timeResponse      time.Duration
		timeTotal         time.Duration
		timeConnIdle      time.Duration
	}
}

func (t *Tracer) Init() *Tracer {
	t.tracerLogs = make(map[int]*tracerLog, 0)
	t.ClientTrace = &httptrace.ClientTrace{
		GetConn: func(hostPort string) {
			t.tracerLogs[TraceGetConn] = &tracerLog{
				tType:   TraceGetConn,
				time:    time.Now(),
				message: fmt.Sprintf("[*] Get Conn: %s\n", hostPort),
			}
		},
		DNSStart: func(info httptrace.DNSStartInfo) {
			t.tracerLogs[TraceDNSStart] = &tracerLog{
				tType:   TraceDNSStart,
				time:    time.Now(),
				message: fmt.Sprintln("[*] DNS Start"),
			}
		},
		DNSDone: func(info httptrace.DNSDoneInfo) {
			t.tracerLogs[TraceDNSDone] = &tracerLog{
				tType:   TraceDNSDone,
				time:    time.Now(),
				message: fmt.Sprintf("[*] DNS Done: %+v\n", info),
			}
		},
		ConnectStart: func(network string, addr string) {
			t.tracerLogs[TraceConnectStart] = &tracerLog{
				tType:   TraceConnectStart,
				time:    time.Now(),
				message: fmt.Sprintf("[*] Connect Start: %s, Addr: %s\n", network, addr),
			}
		},
		ConnectDone: func(network string, addr string, err error) {
			t.tracerLogs[TraceConnectDone] = &tracerLog{
				tType:   TraceConnectDone,
				time:    time.Now(),
				message: fmt.Sprintf("[*] Connect Done: %s, Addr: %s, Error: %+v\n", network, addr, err),
			}
		},
		TLSHandshakeStart: func() {
			t.tracerLogs[TraceTLSHandshakeStart] = &tracerLog{
				tType:   TraceTLSHandshakeStart,
				time:    time.Now(),
				message: fmt.Sprintln("[*] TLS Hand shake Start"),
			}
		},
		TLSHandshakeDone: func(state tls.ConnectionState, err error) {
			t.tracerLogs[TraceTLSHandshakeDone] = &tracerLog{
				tType:   TraceTLSHandshakeDone,
				time:    time.Now(),
				message: fmt.Sprintf("[*] TLSHandshakeDone: %+v, Error: %+v\n", state, err),
			}
		},
		GotConn: func(connInfo httptrace.GotConnInfo) {
			t.tracerLogs[TraceGotConn] = &tracerLog{
				tType:   TraceGotConn,
				time:    time.Now(),
				message: fmt.Sprintf("[*] Got Conn: %+v\n", connInfo),
			}
			t.connInfo = &connInfo
		},
		GotFirstResponseByte: func() {
			t.tracerLogs[TraceGotFirstResponseByte] = &tracerLog{
				tType:   TraceGotFirstResponseByte,
				time:    time.Now(),
				message: fmt.Sprintln("[*] Got First Response Byte"),
			}
		},
	}
	return t
}

func (t *Tracer) ToString() string {
	if t.connInfo != nil {
		t.Output.timeConnIdle = t.connInfo.IdleTime
		t.Output.remoteAddr = t.connInfo.Conn.RemoteAddr().String()
	}

	if s := t.tracerLogs[TraceDNSStart]; s != nil {
		if d := t.tracerLogs[TraceDNSDone]; d != nil {
			t.Output.timeDNSLookup = d.time.Sub(s.time)
		} else {
			t.Output.timeDNSLookup = t.EndTime.Sub(s.time)
		}
	}

	if s := t.tracerLogs[TraceDNSDone]; s != nil {
		if d := t.tracerLogs[TraceConnectStart]; d != nil {
			t.Output.timeTCPConnect = d.time.Sub(s.time)
		}
	}

	if s := t.tracerLogs[TraceConnectStart]; s != nil {
		if d := t.tracerLogs[TraceConnectDone]; d != nil {
			t.Output.timeConnect = d.time.Sub(s.time)
		} else {
			t.Output.timeConnect = t.EndTime.Sub(s.time)
		}
	}

	if s := t.tracerLogs[TraceTLSHandshakeStart]; s != nil {
		if d := t.tracerLogs[TraceTLSHandshakeDone]; d != nil {
			t.Output.timeTLSHandshake = d.time.Sub(s.time)
		} else {
			t.Output.timeTLSHandshake = t.EndTime.Sub(s.time)
		}
	}

	if s := t.tracerLogs[TraceGotConn]; s != nil {
		if d := t.tracerLogs[TraceGotFirstResponseByte]; d != nil {
			t.Output.timeFirstResponse = d.time.Sub(s.time)
			t.Output.timeResponse = t.EndTime.Sub(d.time)
		} else {
			t.Output.timeFirstResponse = t.EndTime.Sub(s.time)
		}
	}

	t.Output.timeTotal = t.EndTime.Sub(t.StartTime)

	return fmt.Sprintf(
		"%s\n"+
			"[>] %-20s: %20s\n"+
			"[v] %-20s: %20v\n"+
			"[v] %-20s: %20v\n"+
			"[v] %-20s: %20v\n"+
			"[v] %-20s: %20v\n"+
			"[v] %-20s: %20v\n"+
			"[v] %-20s: %20v\n"+
			"[-] %-20s: %20v\n"+
			"[x] %-20s: %20v\n",
		t.Output.Prefix,
		"RemoteAddr", t.Output.remoteAddr,
		"DNSLookupCost", t.Output.timeDNSLookup,
		"TCPConnectCost", t.Output.timeTCPConnect,
		"TLSHandshakeCost", t.Output.timeTLSHandshake,
		"ConnectCost", t.Output.timeConnect,
		"FirstResponseCost", t.Output.timeFirstResponse,
		"ResponseTimeCost", t.Output.timeResponse,
		"TotalCost", t.Output.timeTotal,
		"Idle", t.Output.timeConnIdle)
}

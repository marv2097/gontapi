package ntapi

/*
#cgo CFLAGS: -I/opt/napatech3/include
#cgo LDFLAGS: -L/opt/napatech3/lib -lntapi -lpthread -lm
#include <stdlib.h>
#include <nt.h>
#include "macro_wrapper.h"

char errorBuffer[128];
*/
import "C"

import (
	"unsafe"
	"time"
	"errors"
	"bytes"

)

// Struct to hold NTPL Response data (Basic content)
type NtplInfoType struct {
    NtplId 			        uint32				// ID of the NTPL command
    StreamId		        int					// The selected stream ID
    ErrCode			        int32				// Ntpl Parser Error code
    ErrDesc                 []string			// Ntpl Parser Error description
}

// Struct to hold the metadata about the capture packet
type CaptureInfo struct {
    Timestamp 		        time.Time			// Timestamp the packet was captured
    CaptureLength 	        int					// CaptureLength is the total number of bytes read off of the wire.
    Length 			        int					// Length is the size of the original packet.  Should always be >= CaptureLength.
}

// Struct to hold the stats for a port - based on C.struct_NtRMON1Counters_s
type statsPortType struct {
    DropEvents				uint64				//  (ntapi._Ctype_uint64_t) 0,
    Octets					uint64				//  (ntapi._Ctype_uint64_t) 63979910,
    Pkts					uint64				//  (ntapi._Ctype_uint64_t) 294517,
    BroadcastPkts			uint64				//  (ntapi._Ctype_uint64_t) 0,
    MulticastPkts			uint64              //  (ntapi._Ctype_uint64_t) 0,
    CrcAlignErrors  		uint64              //  (ntapi._Ctype_uint64_t) 0,
    UndersizePkts  			uint64              //  (ntapi._Ctype_uint64_t) 0,
    OversizePkts  			uint64              //  (ntapi._Ctype_uint64_t) 0,
    Fragments  				uint64              //  (ntapi._Ctype_uint64_t) 0,
    Jabbers  				uint64              //  (ntapi._Ctype_uint64_t) 0,
    Collisions  			uint64              //  (ntapi._Ctype_uint64_t) 0,
    Pkts64Octets  			uint64              //  (ntapi._Ctype_uint64_t) 1,
    Pkts65to127Octets  		uint64              //  (ntapi._Ctype_uint64_t) 9463,
    Pkts128to255Octets  	uint64              //  (ntapi._Ctype_uint64_t) 283717,
    Pkts256to511Octets  	uint64              //  (ntapi._Ctype_uint64_t) 1156,
    Pkts512to1023Octets  	uint64              //  (ntapi._Ctype_uint64_t) 157,
    Pkts1024to1518Octets	uint64              //  (ntapi._Ctype_uint64_t) 20
}

var hCfgStream C.NtConfigStream_t
var hStatStream C.NtStatStream_t
var ntplInfo C.NtNtplInfo_t
var hNetRx C.NtNetStreamRx_t
var hNetBuf C.NtNetBuf_t
var hStat C.NtStatistics_t


// Initialise the NTAPI 
func NtInit()(err error){
	status := C.NT_Init(C.NTAPI_VERSION)
	if status != C.NT_SUCCESS {
		C.NT_ExplainError(status, &C.errorBuffer[0], 127)
		err = errors.New("NT Init Failed: " + C.GoString(&C.errorBuffer[0]))
	}
	return
} 


// Open a Stat stream
func NtStatOpen(name string)(err error){
    status := C.NT_StatOpen(&hStatStream, C.CString(name))
    if status != C.NT_SUCCESS {
        C.NT_ExplainError(status, &C.errorBuffer[0], 127)
        err = errors.New("NT Stat Stream Open Failed: " + C.GoString(&C.errorBuffer[0]))
    }
    return
}


// Open a configuration
func NtConfigOpen(name string)(err error){
	status := C.NT_ConfigOpen(&hCfgStream, C.CString(name))
	if status != C.NT_SUCCESS {
		C.NT_ExplainError(status, &C.errorBuffer[0], 127)
		err = errors.New("NT Config Open Failed: " + C.GoString(&C.errorBuffer[0]))
	}
	return
}


// Apply a NTPL Statement to a stream
func NtNtpl(ntplBuffer string)(err error, ntplInfo NtplInfoType){
	var CntplInfo C.NtNtplInfo_t
	ntplInfo.ErrDesc = make([]string, 3)

	status := C.NT_NTPL(hCfgStream, C.CString(ntplBuffer), &CntplInfo, C.NT_NTPL_PARSER_VALIDATE_NORMAL)
	// Assign traffic to stream ID 1 and mask all traffic matching the assign statement color=7.
	// eg buffer: "Assign[streamid=1;color=7;txport=1] = All"
	if status != C.NT_SUCCESS {
		C.NT_ExplainError(status, &C.errorBuffer[0], 127)
		err = errors.New("NT NTPL Failed: " + C.GoString(&C.errorBuffer[0]))

		// Map errors (the C struct is not fully defined in go may be due to union statement)
		ntplInfo.ErrDesc[0] = string(CntplInfo.u[:bytes.IndexByte(CntplInfo.u[:C.NT_MAX_NTPL_BUFFER_SIZE], 0)])
		ntplInfo.ErrDesc[1] = string(CntplInfo.u[C.NT_MAX_NTPL_BUFFER_SIZE:bytes.IndexByte(CntplInfo.u[C.NT_MAX_NTPL_BUFFER_SIZE:C.NT_MAX_NTPL_BUFFER_SIZE*2], 0)+C.NT_MAX_NTPL_BUFFER_SIZE])
		ntplInfo.ErrDesc[2] = string(CntplInfo.u[C.NT_MAX_NTPL_BUFFER_SIZE*2:bytes.IndexByte(CntplInfo.u[C.NT_MAX_NTPL_BUFFER_SIZE:C.NT_MAX_NTPL_BUFFER_SIZE*3], 0)+C.NT_MAX_NTPL_BUFFER_SIZE*2])
		ntplInfo.ErrCode = int32(unpack32(CntplInfo.u[C.NT_MAX_NTPL_BUFFER_SIZE*3:]))	
	} else {
		// Success
		ntplInfo.NtplId = uint32(CntplInfo.ntplId)
	}
	return
}


// Close a configuration
func NtConfigClose()(err error){
	status := C.NT_ConfigClose(hCfgStream)
	if status != C.NT_SUCCESS {
		C.NT_ExplainError(status, &C.errorBuffer[0], 127)
		err = errors.New("NT Config Close Failed: " + C.GoString(&C.errorBuffer[0]))
	}
	return
}


// Clear Stats
func NtStatClear()(err error){
    hStat.cmd = C.NT_STATISTICS_READ_CMD_QUERY_V2

    // Configure the Stats to retrieve
    q_v3_ptr := (*C.struct_NtStatisticsQuery_v2_s)(unsafe.Pointer(&hStat.u))
    q_v3_ptr.poll = 0
    q_v3_ptr.clear = 1

    status := C.NT_StatRead(hStatStream, &hStat)
    if status != C.NT_SUCCESS {
        C.NT_ExplainError(status, &C.errorBuffer[0], 127)
        err = errors.New("NT Stat Read Failed: " + C.GoString(&C.errorBuffer[0]))
    }
    return 
}


// Read Stats for a port
func NtStatReadPort(port int)(err error, ntStats statsPortType){
    hStat.cmd = C.NT_STATISTICS_READ_CMD_QUERY_V2

    // Configure the Stats to retrieve
    q_v3_ptr := (*C.struct_NtStatisticsQuery_v2_s)(unsafe.Pointer(&hStat.u)) // This has to be cast to access the union
    q_v3_ptr.poll = 1
    q_v3_ptr.clear = 0

    status := C.NT_StatRead(hStatStream, &hStat)
    if status != C.NT_SUCCESS {
        C.NT_ExplainError(status, &C.errorBuffer[0], 127)
        err = errors.New("NT Stat Read Failed: " + C.GoString(&C.errorBuffer[0]))
    }

    // Return the stats by mapping over to our ntStats struct
    ntStats.DropEvents = uint64(q_v3_ptr.data.port.aPorts[port].rx.RMON1.dropEvents)
    ntStats.Octets = uint64(q_v3_ptr.data.port.aPorts[port].rx.RMON1.octets)
    ntStats.Pkts = uint64(q_v3_ptr.data.port.aPorts[port].rx.RMON1.pkts)
    ntStats.BroadcastPkts = uint64(q_v3_ptr.data.port.aPorts[port].rx.RMON1.broadcastPkts)
    ntStats.MulticastPkts = uint64(q_v3_ptr.data.port.aPorts[port].rx.RMON1.multicastPkts)
    ntStats.CrcAlignErrors = uint64(q_v3_ptr.data.port.aPorts[port].rx.RMON1.crcAlignErrors)
    ntStats.UndersizePkts = uint64(q_v3_ptr.data.port.aPorts[port].rx.RMON1.undersizePkts)
    ntStats.OversizePkts = uint64(q_v3_ptr.data.port.aPorts[port].rx.RMON1.oversizePkts)
    ntStats.Fragments = uint64(q_v3_ptr.data.port.aPorts[port].rx.RMON1.fragments)
    ntStats.Jabbers = uint64(q_v3_ptr.data.port.aPorts[port].rx.RMON1.jabbers)
    ntStats.Collisions = uint64(q_v3_ptr.data.port.aPorts[port].rx.RMON1.collisions)
    ntStats.Pkts64Octets = uint64(q_v3_ptr.data.port.aPorts[port].rx.RMON1.pkts64Octets)
    ntStats.Pkts65to127Octets = uint64(q_v3_ptr.data.port.aPorts[port].rx.RMON1.pkts65to127Octets)
    ntStats.Pkts128to255Octets = uint64(q_v3_ptr.data.port.aPorts[port].rx.RMON1.pkts128to255Octets)
    ntStats.Pkts256to511Octets = uint64(q_v3_ptr.data.port.aPorts[port].rx.RMON1.pkts256to511Octets)
    ntStats.Pkts512to1023Octets = uint64(q_v3_ptr.data.port.aPorts[port].rx.RMON1.pkts512to1023Octets)
    ntStats.Pkts1024to1518Octets = uint64(q_v3_ptr.data.port.aPorts[port].rx.RMON1.pkts1024to1518Octets)
    return 
}

// Get a stream handle with the stream ID. NT_NET_INTERFACE_PACKET specify that we will receive data in a packet based matter.
func NtNetRxOpen(name string, stream uint32)(err error){
	status := C.NT_NetRxOpen(&hNetRx, C.CString(name), C.NT_NET_INTERFACE_PACKET, C.uint32_t(stream), -1)
	if status != C.NT_SUCCESS {
		C.NT_ExplainError(status, &C.errorBuffer[0], 127)
		err = errors.New("NT NetRxOpen Failed: " + C.GoString(&C.errorBuffer[0]))
	}
	return
}


// Gets a single packet from the api and copies it to a new slice
func NtNetRxGetSlice() (ci CaptureInfo, data []byte, err error){

	// Get a packet
	status := C.NT_NetRxGet(hNetRx, &hNetBuf, 1000)
	if status != C.NT_SUCCESS {
		// Get the status code as text
		C.NT_ExplainError(status, &C.errorBuffer[0], 127)
		err = errors.New("NT NetRxGet Failed: " + C.GoString(&C.errorBuffer[0]))
		return
	}

	// Get the Metadata 
	ts_sec := int64(C.nt_net_get_pkt_timestamp(hNetBuf)) / 100000000 
	ts_nsec := (int64(C.nt_net_get_pkt_timestamp(hNetBuf)) - (ts_sec * 100000000)) * 10
	ci.Timestamp = time.Unix(ts_sec, ts_nsec) // NT_TIMESTAMP_TYPE_NATIVE_UNIX 64-bit 10 ns resolution timer from a base of January 1, 1970
	ci.CaptureLength = int(C.nt_net_get_pkt_cap_length(hNetBuf) - C.NT_DESCR_NT_LENGTH)
	ci.Length = int(C.int(C.nt_net_get_pkt_wire_length(hNetBuf)))
	
	// Get the data (GoBytes copies data from ptr into a new slice then returns it)
	data = C.GoBytes(unsafe.Pointer(C.nt_net_get_pkt_l2_ptr(hNetBuf)), C.nt_net_get_pkt_cap_length(hNetBuf) - C.NT_DESCR_NT_LENGTH)
	
	// Release the packet
	status = C.NT_NetRxRelease(hNetRx, hNetBuf)
	if status != C.NT_SUCCESS {
		C.NT_ExplainError(status, &C.errorBuffer[0], 127)
		err = errors.New("NT NetRxRelease Failed: " + C.GoString(&C.errorBuffer[0]))
	}

	// Return the packet info and data
	return
}


// Gets a single packet from the api and copies it to an existing slice
func NtNetRxGetTo(data []byte) (ci CaptureInfo, err error){

	// Get a packet
	status := C.NT_NetRxGet(hNetRx, &hNetBuf, 1000)
	if status != C.NT_SUCCESS {
		// Get the status code as text
		C.NT_ExplainError(status, &C.errorBuffer[0], 127)
		err = errors.New("NT NetRxGet Failed: " + C.GoString(&C.errorBuffer[0]))
		return
	}

	// Get the Metadata 
	ts_sec := int64(C.nt_net_get_pkt_timestamp(hNetBuf)) / 100000000 
	ts_nsec := (int64(C.nt_net_get_pkt_timestamp(hNetBuf)) - (ts_sec * 100000000)) * 10
	ci.Timestamp = time.Unix(ts_sec, ts_nsec) // NT_TIMESTAMP_TYPE_NATIVE_UNIX 64-bit 10 ns resolution timer from a base of January 1, 1970
	ci.CaptureLength = int(C.nt_net_get_pkt_cap_length(hNetBuf) - C.NT_DESCR_NT_LENGTH)
	ci.Length = int(C.int(C.nt_net_get_pkt_wire_length(hNetBuf)))
	
	// Get the data (memcpy copies data from ptr into our existing slice)
	C.memcpy(unsafe.Pointer(&data[0]), C.nt_net_get_pkt_l2_ptr(hNetBuf), C.size_t(C.nt_net_get_pkt_cap_length(hNetBuf) - C.NT_DESCR_NT_LENGTH))

	// Release the packet
	status = C.NT_NetRxRelease(hNetRx, hNetBuf)
	if status != C.NT_SUCCESS {
		C.NT_ExplainError(status, &C.errorBuffer[0], 127)
		err = errors.New("NT NetRxRelease Failed: " + C.GoString(&C.errorBuffer[0]))
	}

	// Return the packet info
	return
}


// Convert a slice of 4 bytes into an integer - Big Endian
func unpack32(ipBytes []byte) uint32 {
	return uint32(ipBytes[0])<<24 + uint32(ipBytes[1])<<16 + uint32(ipBytes[2])<<8 + uint32(ipBytes[3])
}

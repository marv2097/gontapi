#include <nt.h>

unsigned long long nt_net_get_pkt_timestamp(NtNetBuf_t x){
    return NT_NET_GET_PKT_TIMESTAMP(x);
}

int nt_net_get_pkt_timestamp_type(NtNetBuf_t x){
    return NT_NET_GET_PKT_TIMESTAMP(x);
}

int nt_net_get_pkt_wire_length(NtNetBuf_t x){
    return NT_NET_GET_PKT_WIRE_LENGTH(x);
}

int nt_net_get_pkt_cap_length(NtNetBuf_t x){
    return NT_NET_GET_PKT_CAP_LENGTH(x);
}

void* nt_net_get_pkt_l2_ptr(NtNetBuf_t x){
    return NT_NET_GET_PKT_L2_PTR(x);
}
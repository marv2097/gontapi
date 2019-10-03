unsigned long long nt_net_get_pkt_timestamp(NtNetBuf_t x);
int nt_net_get_pkt_timestamp_type(NtNetBuf_t x);
int nt_net_get_pkt_wire_length(NtNetBuf_t x);
int nt_net_get_pkt_cap_length(NtNetBuf_t x);
void* nt_net_get_pkt_l2_ptr(NtNetBuf_t x);
/*
typedef struct goNtNtplInfo_s {   
    enum NtNTPLReturnType_e eType;          //!< Returned status   
    uint32_t                ntplId;         //!< ID of the NTPL command
    int                     streamId;       //!< The selected stream ID   
    uint64_t                ts;             //!< Time when the NTPL command is in effect   
    enum NtTimestampType_e  timestampType;  //!< The time stamp type of NtNtplInfo_t::ts   
    /**    * NTPL return data.    * Error or filter information.    

#ifndef DOXYGEN_INTERNAL_ONLY   
    uint32_t reserved[50]; 
#endif  
    union NtplReturnData_u {     
        struct NtNtplParserErrorData_s   errorData;          //!< Error code and error text 
#ifndef DOXYGEN_INTERNAL_ONLY     
        struct NtNtplFilterCounters_s    aFilterInfo[10];    // Deprecated. Use info stream instead 
//@G-- 
#define NT_NTPL_TEXT_BUFFER_SIZE 3*4*1024     
    char textBuffer[NT_NTPL_TEXT_BUFFER_SIZE];        //!< Text buffer for ntpl debug info 
//@G++ 
#endif   
    } u; 
} NtNtplInfo_t;
*/
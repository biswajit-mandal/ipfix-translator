package msghandler

type Message struct {
	AgentID   string         `json:"AgentID"`
	Header    MessageHeader  `json:"Header"`
	DataSets  []IpfixDataSet `json:"DataSets"`
	Timestamp int64          `json:"Timestamp"`
	RoomKey   string         `json:"roomKey"`
}

type AugmentedMessage struct {
	AgentID   string        `json:"AgentID"`
	Header    MessageHeader `json:"Header"`
	DataSets  IpfixDataSet  `json:"DataSets"`
	Timestamp int64         `json:"Timestamp"`
	RoomKey   string        `json:"roomKey"`
}

// MessageHeader represents IPFIX message header
type MessageHeader struct {
	Version    uint16 // Version of IPFIX to which this Message conforms
	Length     uint16 // Total length of the IPFIX Message, measured in octets
	ExportTime uint64 // Time at which the IPFIX Message Header leaves the Exporter
	SequenceNo uint32 // Incremental sequence counter modulo 2^32
	DomainID   uint32 // A 32-bit id that is locally unique to the Exporting Process
}

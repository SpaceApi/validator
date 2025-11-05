package v3

import (
	"bytes"
	"encoding/json"
)

type AddressFamily int

const (
	Any AddressFamily = iota
	IPv6
	IPv4
)

var toString = map[AddressFamily]string{
	Any:  "any",
	IPv6: "ipv6",
	IPv4: "ipv4",
}

var toAddressFamily = map[string]AddressFamily{
	"any":  Any,
	"ipv6": IPv6,
	"ipv4": IPv4,
}

var toTcpStack = map[AddressFamily]string{
	Any:  "tcp",
	IPv6: "tcp6",
	IPv4: "tcp4",
}

func (f AddressFamily) String() string {
	return toString[f]
}

func (f AddressFamily) TcpStack() string {
	return toTcpStack[f]
}

func (f AddressFamily) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(toString[f])
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (f *AddressFamily) UnmarshalJSON(b []byte) error {
	var str string
	err := json.Unmarshal(b, &str)
	if err != nil {
		return err
	}
	*f = toAddressFamily[str]
	return nil
}

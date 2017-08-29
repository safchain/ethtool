package ethtool

import (
	"fmt"
)

type sff8079 struct {
	ExtIdentifier   string
	Connector       string
	TransceiverType string
	Encoding        string
	BRNominal       string
}

func ParseSFF8079(id []byte) (sff8079, error) {
	if id[0] != 0x03 && id[1] != 0x04 {
		return sff8079{}, fmt.Errorf("unknown eeprom format, not sff-8079")
	}

	sff := sff8079{}

	// External Identifier
	sff.ExtIdentifier = fmt.Sprintf("0x%02x ", id[1])
	switch id[1] {
	case 0x00:
		sff.ExtIdentifier += "(GBIC not specified / not MOD_DEF compliant)"
	case 0x04:
		sff.ExtIdentifier += "(GBIC/SFP defined by 2-wire interface ID)"
	case 0x05, 0x06, 0x07:
		sff.ExtIdentifier += fmt.Sprintf("(GBIC compliant with MOD_DEF %u)", id[1])
	default:
		sff.ExtIdentifier += "(unknown)"
	}

	// Connector
	sff.Connector = fmt.Sprintf("0x%02x ", id[2])
	switch id[2] {
	case 0x00:
		sff.Connector += "(unknown or unspecified)"
	case 0x01:
		sff.Connector += "(SC)"
	case 0x02:
		sff.Connector += "(Fibre Channel Style 1 copper)"
	case 0x03:
		sff.Connector += "(Fibre Channel Style 2 copper)"
	case 0x04:
		sff.Connector += "(BNC/TNC)"
	case 0x05:
		sff.Connector += "(Fibre Channel coaxial headers)"
	case 0x06:
		sff.Connector += "(FibreJack)"
	case 0x07:
		sff.Connector += "(LC)"
	case 0x08:
		sff.Connector += "(MT-RJ)"
	case 0x09:
		sff.Connector += "(MU)"
	case 0x0a:
		sff.Connector += "(SG)"
	case 0x0b:
		sff.Connector += "(Optical pigtail)"
	case 0x0c:
		sff.Connector += "(MPO Parallel Optic)"
	case 0x20:
		sff.Connector += "(HSSDC II)"
	case 0x21:
		sff.Connector += "(Copper pigtail)"
	case 0x22:
		sff.Connector += "(RJ45)"
	default:
		sff.Connector += "(reserved or unknown)"
	}

	return sff, nil
}

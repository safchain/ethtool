package ethtool

import (
	"fmt"
)

type sff8079 struct {
	ExtIdentifier string
	Connector     string
	TransCodes    string
	TransTypes    []string
}

func ParseSFF8079(id []byte) (sff8079, error) {
	if id[0] != 0x03 && id[1] != 0x04 {
		return sff8079{}, fmt.Errorf("unknown eeprom format, not sff-8079")
	}

	sff := sff8079{}

	// External identifier
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

	// Transceiver codes
	sff.TransCodes = fmt.Sprintf("0x%02x 0x%02x 0x%02x 0x%02x 0x%02x 0x%02x 0x%02x 0x%02x",
		id[3], id[4], id[5], id[6],
		id[7], id[8], id[9], id[10])

	/* 10G Ethernet Compliance Codes */
	if id[3]&(1<<7) != 0 {
		sff.TransTypes = append(sff.TransTypes, "10G Ethernet: 10G Base-ER [SFF-8472 rev10.4 only]")
	}
	if id[3]&(1<<6) != 0 {
		sff.TransTypes = append(sff.TransTypes, "10G Ethernet: 10G Base-LRM")
	}
	if id[3]&(1<<5) != 0 {
		sff.TransTypes = append(sff.TransTypes, "10G Ethernet: 10G Base-LR")
	}
	if id[3]&(1<<4) != 0 {
		sff.TransTypes = append(sff.TransTypes, "10G Ethernet: 10G Base-SR")
	}

	/* Infiniband Compliance Codes */
	if id[3]&(1<<3) != 0 {
		sff.TransTypes = append(sff.TransTypes, "Infiniband: 1X SX")
	}
	if id[3]&(1<<2) != 0 {
		sff.TransTypes = append(sff.TransTypes, "Infiniband: 1X LX")
	}
	if id[3]&(1<<1) != 0 {
		sff.TransTypes = append(sff.TransTypes, "Infiniband: 1X Copper Active")
	}
	if id[3]&(1<<0) != 0 {
		sff.TransTypes = append(sff.TransTypes, "Infiniband: 1X Copper Passive")
	}

	/* ESCON Compliance Codes */
	if id[4]&(1<<7) != 0 {
		sff.TransTypes = append(sff.TransTypes, "ESCON: ESCON MMF, 1310nm LED")
	}
	if id[4]&(1<<6) != 0 {
		sff.TransTypes = append(sff.TransTypes, "ESCON: ESCON SMF, 1310nm Laser")
	}

	/* SONET Compliance Codes */
	if id[4]&(1<<5) != 0 {
		sff.TransTypes = append(sff.TransTypes, "SONET: OC-192, short reach")
	}
	if id[4]&(1<<4) != 0 {
		sff.TransTypes = append(sff.TransTypes, "SONET: SONET reach specifier bit 1")
	}
	if id[4]&(1<<3) != 0 {
		sff.TransTypes = append(sff.TransTypes, "SONET: SONET reach specifier bit 2")
	}
	if id[4]&(1<<2) != 0 {
		sff.TransTypes = append(sff.TransTypes, "SONET: OC-48, long reach")
	}
	if id[4]&(1<<1) != 0 {
		sff.TransTypes = append(sff.TransTypes, "SONET: OC-48, intermediate reach")
	}
	if id[4]&(1<<0) != 0 {
		sff.TransTypes = append(sff.TransTypes, "SONET: OC-48, short reach")
	}
	if id[5]&(1<<6) != 0 {
		sff.TransTypes = append(sff.TransTypes, "SONET: OC-12, single mode, long reach")
	}
	if id[5]&(1<<5) != 0 {
		sff.TransTypes = append(sff.TransTypes, "SONET: OC-12, single mode, inter. reach")
	}
	if id[5]&(1<<4) != 0 {
		sff.TransTypes = append(sff.TransTypes, "SONET: OC-12, short reach")
	}
	if id[5]&(1<<2) != 0 {
		sff.TransTypes = append(sff.TransTypes, "SONET: OC-3, single mode, long reach")
	}
	if id[5]&(1<<1) != 0 {
		sff.TransTypes = append(sff.TransTypes, "SONET: OC-3, single mode, inter. reach")
	}
	if id[5]&(1<<0) != 0 {
		sff.TransTypes = append(sff.TransTypes, "SONET: OC-3, short reach")
	}

	/* Ethernet Compliance Codes */
	if id[6]&(1<<7) != 0 {
		sff.TransTypes = append(sff.TransTypes, "Ethernet: BASE-PX")
	}
	if id[6]&(1<<6) != 0 {
		sff.TransTypes = append(sff.TransTypes, "Ethernet: BASE-BX10")
	}
	if id[6]&(1<<5) != 0 {
		sff.TransTypes = append(sff.TransTypes, "Ethernet: 100BASE-FX")
	}
	if id[6]&(1<<4) != 0 {
		sff.TransTypes = append(sff.TransTypes, "Ethernet: 100BASE-LX/LX10")
	}
	if id[6]&(1<<3) != 0 {
		sff.TransTypes = append(sff.TransTypes, "Ethernet: 1000BASE-T")
	}
	if id[6]&(1<<2) != 0 {
		sff.TransTypes = append(sff.TransTypes, "Ethernet: 1000BASE-CX")
	}
	if id[6]&(1<<1) != 0 {
		sff.TransTypes = append(sff.TransTypes, "Ethernet: 1000BASE-LX")
	}
	if id[6]&(1<<0) != 0 {
		sff.TransTypes = append(sff.TransTypes, "Ethernet: 1000BASE-SX")
	}

	/* Fibre Channel link length */
	if id[7]&(1<<7) != 0 {
		sff.TransTypes = append(sff.TransTypes, "FC: very long distance (V)")
	}
	if id[7]&(1<<6) != 0 {
		sff.TransTypes = append(sff.TransTypes, "FC: short distance (S)")
	}
	if id[7]&(1<<5) != 0 {
		sff.TransTypes = append(sff.TransTypes, "FC: intermediate distance (I)")
	}
	if id[7]&(1<<4) != 0 {
		sff.TransTypes = append(sff.TransTypes, "FC: long distance (L)")
	}
	if id[7]&(1<<3) != 0 {
		sff.TransTypes = append(sff.TransTypes, "FC: medium distance (M)")
	}

	/* Fibre Channel transmitter technology */
	if id[7]&(1<<2) != 0 {
		sff.TransTypes = append(sff.TransTypes, "FC: Shortwave laser, linear Rx (SA)")
	}
	if id[7]&(1<<1) != 0 {
		sff.TransTypes = append(sff.TransTypes, "FC: Longwave laser (LC)")
	}
	if id[7]&(1<<0) != 0 {
		sff.TransTypes = append(sff.TransTypes, "FC: Electrical inter-enclosure (EL)")
	}
	if id[8]&(1<<7) != 0 {
		sff.TransTypes = append(sff.TransTypes, "FC: Electrical intra-enclosure (EL)")
	}
	if id[8]&(1<<6) != 0 {
		sff.TransTypes = append(sff.TransTypes, "FC: Shortwave laser w/o OFC (SN)")
	}
	if id[8]&(1<<5) != 0 {
		sff.TransTypes = append(sff.TransTypes, "FC: Shortwave laser with OFC (SL)")
	}
	if id[8]&(1<<4) != 0 {
		sff.TransTypes = append(sff.TransTypes, "FC: Longwave laser (LL)")
	}
	if id[8]&(1<<3) != 0 {
		sff.TransTypes = append(sff.TransTypes, "FC: Copper Active")
	}
	if id[8]&(1<<2) != 0 {
		sff.TransTypes = append(sff.TransTypes, "FC: Copper Passive")
	}
	if id[8]&(1<<1) != 0 {
		sff.TransTypes = append(sff.TransTypes, "FC: Copper FC-BaseT")
	}

	/* Fibre Channel transmission media */
	if id[9]&(1<<7) != 0 {
		sff.TransTypes = append(sff.TransTypes, "FC: Twin Axial Pair (TW)")
	}
	if id[9]&(1<<6) != 0 {
		sff.TransTypes = append(sff.TransTypes, "FC: Twisted Pair (TP)")
	}
	if id[9]&(1<<5) != 0 {
		sff.TransTypes = append(sff.TransTypes, "FC: Miniature Coax (MI)")
	}
	if id[9]&(1<<4) != 0 {
		sff.TransTypes = append(sff.TransTypes, "FC: Video Coax (TV)")
	}
	if id[9]&(1<<3) != 0 {
		sff.TransTypes = append(sff.TransTypes, "FC: Multimode, 62.5um (M6)")
	}
	if id[9]&(1<<2) != 0 {
		sff.TransTypes = append(sff.TransTypes, "FC: Multimode, 50um (M5)")
	}
	if id[9]&(1<<0) != 0 {
		sff.TransTypes = append(sff.TransTypes, "FC: Single Mode (SM)")
	}

	/* Fibre Channel speed */
	if id[10]&(1<<7) != 0 {
		sff.TransTypes = append(sff.TransTypes, "FC: 1200 MBytes/sec")
	}
	if id[10]&(1<<6) != 0 {
		sff.TransTypes = append(sff.TransTypes, "FC: 800 MBytes/sec")
	}
	if id[10]&(1<<4) != 0 {
		sff.TransTypes = append(sff.TransTypes, "FC: 400 MBytes/sec")
	}
	if id[10]&(1<<2) != 0 {
		sff.TransTypes = append(sff.TransTypes, "FC: 200 MBytes/sec")
	}
	if id[10]&(1<<0) != 0 {
		sff.TransTypes = append(sff.TransTypes, "FC: 100 MBytes/sec")
	}

	return sff, nil
}

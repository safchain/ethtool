package ethtool

import (
	"fmt"
	"strings"
	"time"
	"unsafe"
)

type SFF8079 struct {
	ExtIdentifier  string    `json:"external_identifier"`
	Connector      string    `json:"connector"`
	TransCodes     string    `json:"transceiver_codes"`
	TransTypes     []string  `json:"transceiver_types"`
	Encoding       string    `json:"encoding"`
	BRNominalMBd   uint      `json:"br_nominal_mbd"`
	RateIdentifier string    `json:"rate_identifier"`
	CableSMFLenKm  uint      `json:"cable_smf_length_km,omitempty"`
	CableSMFLenM   uint      `json:"cable_smf_length_m,omitempty"`
	Cable50umLenM  uint      `json:"cable_50um_length_m,omitempty"`
	Cable625umLenM uint      `json:"cable_62_5um_length_m,omitempty"`
	CableCoprLenM  uint      `json:"cable_copper_length_m,omitempty"`
	CableOM3LenM   uint      `json:"cable_om3_length_m,omitempty"`
	PasveCuCompl   string    `json:"passive_cu_compliant,omitempty"`
	ActveCuCompl   string    `json:"active_cu_compliant,omitempty"`
	LaserWavelen   string    `json:"laser_wavelength,omitempty"`
	VendorName     string    `json:"vendor_name"`
	VendorOUI      string    `json:"vendor_oui"`
	VendorPN       string    `json:"vendor_pn"`
	VendorRev      string    `json:"vendor_rev"`
	OptionVals     string    `json:"option_vals"`
	Option         string    `json:"option"`
	BRMargMaxPerc  uint      `json:"br_margin_max_perc"`
	BRMargMinPerc  uint      `json:"br_margin_min_perc"`
	VendorSN       string    `json:"vendor_sn"`
	VendorDate     time.Time `json:"vendor_date"`
	DateCode       string    `json:"date_code"`
}

func ParseSFF8079(id []byte) (*SFF8079, error) {
	if id[0] != 0x03 && id[1] != 0x04 {
		return nil, fmt.Errorf("unknown eeprom format, not sff-8079")
	}

	sff := SFF8079{}

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

	// Encoding
	sff.Encoding = fmt.Sprintf("0x%02x ", id[11])
	switch id[11] {
	case 0x00:
		sff.Encoding += "(unspecified)"
	case 0x01:
		sff.Encoding += "(8B/10B)"
	case 0x02:
		sff.Encoding += "(4B/5B)"
	case 0x03:
		sff.Encoding += "(NRZ)"
	case 0x04:
		sff.Encoding += "(Manchester)"
	case 0x05:
		sff.Encoding += "(SONET Scrambled)"
	case 0x06:
		sff.Encoding += "(64B/66B)"
	default:
		sff.Encoding += "(reserved or unknown)"
	}

	// BR nominal
	v := *(*uint8)(unsafe.Pointer(&id[12]))
	sff.BRNominalMBd = uint(v) * 100

	// Rate identifier
	sff.RateIdentifier = fmt.Sprintf("0x%02x ", id[13])
	switch id[13] {
	case 0x00:
		sff.RateIdentifier += "(unspecified)"
	case 0x01:
		sff.RateIdentifier += "(4/2/1G Rate_Select & AS0/AS1)"
	case 0x02:
		sff.RateIdentifier += "(8/4/2G Rx Rate_Select only)"
	case 0x03:
		sff.RateIdentifier += "(8/4/2G Independent Rx & Tx Rate_Select)"
	case 0x04:
		sff.RateIdentifier += "(8/4/2G Tx Rate_Select only)"
	default:
		sff.RateIdentifier += "(reserved or unknown)"
	}

	// Length SMF km
	v = *(*uint8)(unsafe.Pointer(&id[14]))
	sff.CableSMFLenKm = uint(v)

	// Length SMF
	v = *(*uint8)(unsafe.Pointer(&id[15]))
	sff.CableSMFLenM = uint(v) * 100

	// Length 50um
	v = *(*uint8)(unsafe.Pointer(&id[16]))
	sff.Cable50umLenM = uint(v) * 10

	// Length 62.5um
	v = *(*uint8)(unsafe.Pointer(&id[17]))
	sff.Cable625umLenM = uint(v) * 10

	// Length copper
	v = *(*uint8)(unsafe.Pointer(&id[18]))
	sff.CableCoprLenM = uint(v)

	// Length OM3
	v = *(*uint8)(unsafe.Pointer(&id[19]))
	sff.CableOM3LenM = uint(v) * 10

	// Passive cu compliance
	// Active cu compliance
	// Laser wave length
	if id[8]&(1<<2) != 0 {
		sff.PasveCuCompl = fmt.Sprintf("0x%02x ", id[60])
		switch id[60] {
		case 0x00:
			sff.PasveCuCompl += "(unspecified)"
		case 0x01:
			sff.PasveCuCompl += "(SFF-8431 appendix E)"
		default:
			sff.PasveCuCompl += "(unknown)"
		}
		sff.PasveCuCompl += " [SFF-8472 rev10.4 only]"
	} else if id[8]&(1<<3) != 0 {
		sff.ActveCuCompl = fmt.Sprintf("0x%02x ", id[60])
		switch id[60] {
		case 0x00:
			sff.ActveCuCompl += "(unspecified)"
		case 0x01:
			sff.ActveCuCompl += "(SFF-8431 appendix E)"
		case 0x04:
			sff.ActveCuCompl += "(SFF-8431 limiting)"
		default:
			sff.ActveCuCompl += "(unknown)"
		}
		sff.ActveCuCompl += " [SFF-8472 rev10.4 only]"
	} else {
		sff.LaserWavelen = fmt.Sprintf("%u%s", (id[60]<<8)|id[61], "nm")
	}

	// Vendor name
	for i := 20; i <= 35; i++ {
		sff.VendorName += string(id[i])
	}
	sff.VendorName = strings.TrimSpace(sff.VendorName)

	// Vendor OUI
	sff.VendorOUI = fmt.Sprintf("%02x:%02x:%02x", id[37], id[38], id[39])

	// Vendor PN
	for i := 40; i <= 55; i++ {
		sff.VendorPN += string(id[i])
	}
	sff.VendorPN = strings.TrimSpace(sff.VendorPN)

	// Vendor rev
	for i := 56; i <= 59; i++ {
		sff.VendorRev += string(id[i])
	}
	sff.VendorRev = strings.TrimSpace(sff.VendorRev)

	// Options values
	sff.OptionVals = fmt.Sprintf("0x%02x 0x%02x", id[64], id[65])
	if id[65]&(1<<1) != 0 {
		sff.Option += "RX_LOS implemented"
	}
	if id[65]&(1<<2) != 0 {
		sff.Option += "RX_LOS implemented, inverted"
	}
	if id[65]&(1<<3) != 0 {
		sff.Option += "TX_FAULT implemented"
	}
	if id[65]&(1<<4) != 0 {
		sff.Option += "TX_DISABLE implemented"
	}
	if id[65]&(1<<5) != 0 {
		sff.Option += "RATE_SELECT implemented"
	}
	if id[65]&(1<<6) != 0 {
		sff.Option += "Tunable transmitter technology"
	}
	if id[65]&(1<<7) != 0 {
		sff.Option += "Receiver decision threshold implemented"
	}
	if id[64]&(1<<0) != 0 {
		sff.Option += "Linear receiver output implemented"
	}
	if id[64]&(1<<1) != 0 {
		sff.Option += "Power level 2 requirement"
	}
	if id[64]&(1<<2) != 0 {
		sff.Option += "Cooled transceiver implemented"
	}
	if id[64]&(1<<3) != 0 {
		sff.Option += "Retimer or CDR implemented"
	}
	if id[64]&(1<<4) != 0 {
		sff.Option += "Paging implemented"
	}
	if id[64]&(1<<5) != 0 {
		sff.Option += "Power level 3 requirement"
	}

	// BR margin max
	v = *(*uint8)(unsafe.Pointer(&id[66]))
	sff.BRMargMaxPerc = uint(v)

	// BR margin min
	v = *(*uint8)(unsafe.Pointer(&id[67]))
	sff.BRMargMinPerc = uint(v)

	// Vendor SN
	for i := 68; i <= 83; i++ {
		sff.VendorSN += string(id[i])
	}
	sff.VendorSN = strings.TrimSpace(sff.VendorSN)

	// Vendor Date
	t, err := time.Parse("2006-01-02", fmt.Sprintf("20%s%s-%s%s-%s%s",
		string(id[84]), string(id[85]), string(id[86]), string(id[87]), string(id[88]), string(id[89])))
	if err != nil {
		return nil, fmt.Errorf("parse date: %v", err)
	}
	sff.VendorDate = t

	// Date code
	for i := 84; i <= 91; i++ {
		sff.DateCode += string(id[i])
	}
	sff.DateCode = strings.TrimSpace(sff.DateCode)

	return &sff, nil
}
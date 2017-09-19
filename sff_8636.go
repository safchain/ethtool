package ethtool

import (
	"fmt"
)

// TODO:
// - Separate struct for data and test descr

type SFF8636 struct {
	Identifier         string   `json:"identifier"`
	ExtIdentifier      string   `json:"extIdentifier"`
	ExtIdentifierDescr []string `json:"extIdentifierDescr"`
}

func sff8636ShowIdentifier(id []byte) string {
	return sff8024ShowIdentifier(id, SFF8636_ID_OFFSET)
}

func sff8636ShowExtIdentifier(id []byte) string {
	return fmt.Sprintf("0x%02x", id[SFF8636_EXT_ID_OFFSET])
}

func sff8636ShowExtIdentifierDescr(id []byte) []string {
	descr := []string{}

	// Probably doesn't work properly with AND
	switch id[SFF8636_EXT_ID_OFFSET] & SFF8636_EXT_ID_PWR_CLASS_MASK {
	case SFF8636_EXT_ID_PWR_CLASS_1:
		descr = append(descr, "1.5W max. Power consumption")
	case SFF8636_EXT_ID_PWR_CLASS_2:
		descr = append(descr, "2.0W max. Power consumption")
	case SFF8636_EXT_ID_PWR_CLASS_3:
		descr = append(descr, "2.5W max. Power consumption")
	case SFF8636_EXT_ID_PWR_CLASS_4:
		descr = append(descr, "3.5W max. Power consumption")
	}

	if id[SFF8636_EXT_ID_OFFSET]&SFF8636_EXT_ID_CDR_TX_MASK != 0 {
		descr = append(descr, "CDR present in TX,")
	} else {
		descr = append(descr, "No CDR in TX,")
	}

	if id[SFF8636_EXT_ID_OFFSET]&SFF8636_EXT_ID_CDR_RX_MASK != 0 {
		descr = append(descr, "CDR present in RX")
	} else {
		descr = append(descr, "No CDR in RX")
	}

	// Probably doesn't work properly with AND
	switch id[SFF8636_EXT_ID_OFFSET] & SFF8636_EXT_ID_EPWR_CLASS_MASK {
	//	case SFF8636_EXT_ID_PWR_CLASS_LEGACY:
	case SFF8636_EXT_ID_PWR_CLASS_5:
		descr = append(descr, "4.0W max. Power consumption,")
	case SFF8636_EXT_ID_PWR_CLASS_6:
		descr = append(descr, "4.5W max. Power consumption,")
	case SFF8636_EXT_ID_PWR_CLASS_7:
		descr = append(descr, "5.0W max. Power consumption,")
	}

	if id[SFF8636_PWR_MODE_OFFSET]&SFF8636_HIGH_PWR_ENABLE != 0 {
		descr = append(descr, "High Power Class (> 3.5 W) enabled")
	} else {
		descr = append(descr, "High Power Class (> 3.5 W) not enabled")
	}

	return descr
}

func sff8636ShowConnector(id []byte) string {
	return sff8024ShowConnector(id, SFF8636_CTOR_OFFSET)
}

func sff8636ShowTransceiverCodes(id []byte) string {
	return fmt.Sprintf("0x%02x 0x%02x 0x%02x 0x%02x 0x%02x 0x%02x 0x%02x 0x%02x",
		id[SFF8636_ETHERNET_COMP_OFFSET],
		id[SFF8636_SONET_COMP_OFFSET],
		id[SFF8636_SAS_COMP_OFFSET],
		id[SFF8636_GIGE_COMP_OFFSET],
		id[SFF8636_FC_LEN_OFFSET],
		id[SFF8636_FC_TECH_OFFSET],
		id[SFF8636_FC_TRANS_MEDIA_OFFSET],
		id[SFF8636_FC_SPEED_OFFSET])
}

func sff8636ShowTransceiverType(id []byte) string {
	/* 10G/40G Ethernet Compliance Codes */
	if id[SFF8636_ETHERNET_COMP_OFFSET]&SFF8636_ETHERNET_10G_LRM != 0 {
		return "10G Ethernet: 10G Base-LRM"
	}
	if id[SFF8636_ETHERNET_COMP_OFFSET]&SFF8636_ETHERNET_10G_LR != 0 {
		return "10G Ethernet: 10G Base-LR"
	}
	if id[SFF8636_ETHERNET_COMP_OFFSET]&SFF8636_ETHERNET_10G_SR != 0 {
		return "10G Ethernet: 10G Base-SR"
	}
	if id[SFF8636_ETHERNET_COMP_OFFSET]&SFF8636_ETHERNET_40G_CR4 != 0 {
		return "40G Ethernet: 40G Base-CR4"
	}
	if id[SFF8636_ETHERNET_COMP_OFFSET]&SFF8636_ETHERNET_40G_SR4 != 0 {
		return "40G Ethernet: 40G Base-SR4"
	}
	if id[SFF8636_ETHERNET_COMP_OFFSET]&SFF8636_ETHERNET_40G_LR4 != 0 {
		return "40G Ethernet: 40G Base-LR4"
	}
	if id[SFF8636_ETHERNET_COMP_OFFSET]&SFF8636_ETHERNET_40G_ACTIVE != 0 {
		return "40G Ethernet: 40G Active Cable (XLPPI)"
	}

	/* Extended Specification Compliance Codes from SFF-8024 */
	if id[SFF8636_ETHERNET_COMP_OFFSET]&SFF8636_ETHERNET_RSRVD != 0 {
		switch id[SFF8636_OPTION_1_OFFSET] {
		case SFF8636_ETHERNET_UNSPECIFIED:
			return "(reserved or unknown)"
		case SFF8636_ETHERNET_100G_AOC:
			return "100G Ethernet: 100G AOC or 25GAUI C2M AOC with worst BER of 5x10^(-5)"
		case SFF8636_ETHERNET_100G_SR4:
			return "100G Ethernet: 100G Base-SR4 or 25GBase-SR"
		case SFF8636_ETHERNET_100G_LR4:
			return "100G Ethernet: 100G Base-LR4"
		case SFF8636_ETHERNET_100G_ER4:
			return "100G Ethernet: 100G Base-ER4"
		case SFF8636_ETHERNET_100G_SR10:
			return "100G Ethernet: 100G Base-SR10"
		case SFF8636_ETHERNET_100G_CWDM4_FEC:
			return "100G Ethernet: 100G CWDM4 MSA with FEC"
		case SFF8636_ETHERNET_100G_PSM4:
			return "100G Ethernet: 100G PSM4 Parallel SMF"
		case SFF8636_ETHERNET_100G_ACC:
			return "100G Ethernet: 100G ACC or 25GAUI C2M ACC with worst BER of 5x10^(-5)"
		case SFF8636_ETHERNET_100G_CWDM4_NO_FEC:
			return "100G Ethernet: 100G CWDM4 MSA without FEC"
		case SFF8636_ETHERNET_100G_RSVD1:
			return "(reserved or unknown)"
		case SFF8636_ETHERNET_100G_CR4:
			return "100G Ethernet: 100G Base-CR4 or 25G Base-CR CA-L"
		case SFF8636_ETHERNET_25G_CR_CA_S:
			return "25G Ethernet: 25G Base-CR CA-S"
		case SFF8636_ETHERNET_25G_CR_CA_N:
			return "25G Ethernet: 25G Base-CR CA-N"
		case SFF8636_ETHERNET_40G_ER4:
			return "40G Ethernet: 40G Base-ER4"
		case SFF8636_ETHERNET_4X10_SR:
			return "4x10G Ethernet: 10G Base-SR"
		case SFF8636_ETHERNET_40G_PSM4:
			return "40G Ethernet: 40G PSM4 Parallel SMF"
		case SFF8636_ETHERNET_G959_P1I1_2D1:
			return "Ethernet: G959.1 profile P1I1-2D1 (10709 MBd, 2km, 1310nm SM)"
		case SFF8636_ETHERNET_G959_P1S1_2D2:
			return "Ethernet: G959.1 profile P1S1-2D2 (10709 MBd, 40km, 1550nm SM)"
		case SFF8636_ETHERNET_G959_P1L1_2D2:
			return "Ethernet: G959.1 profile P1L1-2D2 (10709 MBd, 80km, 1550nm SM)"
		case SFF8636_ETHERNET_10GT_SFI:
			return "10G Ethernet: 10G Base-T with SFI electrical interface"
		case SFF8636_ETHERNET_100G_CLR4:
			return "100G Ethernet: 100G CLR4"
		case SFF8636_ETHERNET_100G_AOC2:
			return "100G Ethernet: 100G AOC or 25GAUI C2M AOC with worst BER of 10^(-12)"
		case SFF8636_ETHERNET_100G_ACC2:
			return "100G Ethernet: 100G ACC or 25GAUI C2M ACC with worst BER of 10^(-12)"
		}
		return "(reserved or unknown)"
	}

	/* SONET Compliance Codes */
	if id[SFF8636_SONET_COMP_OFFSET]&SFF8636_SONET_40G_OTN != 0 {
		return "40G OTN (OTU3B/OTU3C)"
	}
	if id[SFF8636_SONET_COMP_OFFSET]&SFF8636_SONET_OC48_LR != 0 {
		return "SONET: OC-48, long reach"
	}
	if id[SFF8636_SONET_COMP_OFFSET]&SFF8636_SONET_OC48_IR != 0 {
		return "SONET: OC-48, intermediate reach"
	}
	if id[SFF8636_SONET_COMP_OFFSET]&SFF8636_SONET_OC48_SR != 0 {
		return "SONET: OC-48, short reach"
	}

	/* SAS/SATA Compliance Codes */
	if id[SFF8636_SAS_COMP_OFFSET]&SFF8636_SAS_6G != 0 {
		return "SAS 6.0G"
	}
	if id[SFF8636_SAS_COMP_OFFSET]&SFF8636_SAS_3G != 0 {
		return "SAS 3.0G"
	}

	/* Ethernet Compliance Codes */
	if id[SFF8636_GIGE_COMP_OFFSET]&SFF8636_GIGE_1000_BASE_T != 0 {
		return "Ethernet: 1000BASE-T"
	}
	if id[SFF8636_GIGE_COMP_OFFSET]&SFF8636_GIGE_1000_BASE_CX != 0 {
		return "Ethernet: 1000BASE-CX"
	}
	if id[SFF8636_GIGE_COMP_OFFSET]&SFF8636_GIGE_1000_BASE_LX != 0 {
		return "Ethernet: 1000BASE-LX"
	}
	if id[SFF8636_GIGE_COMP_OFFSET]&SFF8636_GIGE_1000_BASE_SX != 0 {
		return "Ethernet: 1000BASE-SX"
	}

	/* Fibre Channel link length */
	if id[SFF8636_FC_LEN_OFFSET]&SFF8636_FC_LEN_VERY_LONG != 0 {
		return "FC: very long distance (V)"
	}
	if id[SFF8636_FC_LEN_OFFSET]&SFF8636_FC_LEN_SHORT != 0 {
		return "FC: short distance (S)"
	}
	if id[SFF8636_FC_LEN_OFFSET]&SFF8636_FC_LEN_INT != 0 {
		return "FC: intermediate distance (I)"
	}
	if id[SFF8636_FC_LEN_OFFSET]&SFF8636_FC_LEN_LONG != 0 {
		return "FC: long distance (L)"
	}
	if id[SFF8636_FC_LEN_OFFSET]&SFF8636_FC_LEN_MED != 0 {
		return "FC: medium distance (M)"
	}

	return ""
}

func Decode(id []byte) (*SFF8636, error) {
	s := &SFF8636{
		Identifier: sff8636ShowIdentifier(id),
	}

	if id[SFF8636_ID_OFFSET] == SFF8024_ID_QSFP ||
		id[SFF8636_ID_OFFSET] == SFF8024_ID_QSFP_PLUS ||
		id[SFF8636_ID_OFFSET] == SFF8024_ID_QSFP28 {

		s.ExtIdentifier = sff8636ShowExtIdentifier(id)
		s.ExtIdentifierDescr = sff8636ShowExtIdentifierDescr(id)
		//		s.Connector = sff8636ShowConnector(id)
		//		s.Transceiver = sff8636ShowTransceiver(id)
		//		s.Encoding = sff8636ShowEncoding(id)
	}

	return s, nil
}

/*
void sff8636_show_all(const __u8 *id, __u32 eeprom_len)
{
        sff8636_show_identifier(id);
        if ((id[SFF8636_ID_OFFSET] == SFF8024_ID_QSFP) ||
                (id[SFF8636_ID_OFFSET] == SFF8024_ID_QSFP_PLUS) ||
                (id[SFF8636_ID_OFFSET] == SFF8024_ID_QSFP28)) {
                sff8636_show_ext_identifier(id);
                sff8636_show_connector(id);
                sff8636_show_transceiver(id);
                sff8636_show_encoding(id);
                sff_show_value_with_unit(id, SFF8636_BR_NOMINAL_OFFSET,
                                "BR, Nominal", 100, "Mbps");
                sff8636_show_rate_identifier(id);
                sff_show_value_with_unit(id, SFF8636_SM_LEN_OFFSET,
                             "Length (SMF,km)", 1, "km");
                sff_show_value_with_unit(id, SFF8636_OM3_LEN_OFFSET,
                                "Length (OM3 50um)", 2, "m");
                sff_show_value_with_unit(id, SFF8636_OM2_LEN_OFFSET,
                                "Length (OM2 50um)", 1, "m");
                sff_show_value_with_unit(id, SFF8636_OM1_LEN_OFFSET,
                             "Length (OM1 62.5um)", 1, "m");
                sff_show_value_with_unit(id, SFF8636_CBL_LEN_OFFSET,
                             "Length (Copper or Active cable)", 1, "m");
                sff8636_show_wavelength_or_copper_compliance(id);
                sff_show_ascii(id, SFF8636_VENDOR_NAME_START_OFFSET,
                               SFF8636_VENDOR_NAME_END_OFFSET, "Vendor name");
                sff8636_show_oui(id);
                sff_show_ascii(id, SFF8636_VENDOR_PN_START_OFFSET,
                               SFF8636_VENDOR_PN_END_OFFSET, "Vendor PN");
                sff_show_ascii(id, SFF8636_VENDOR_REV_START_OFFSET,
                               SFF8636_VENDOR_REV_END_OFFSET, "Vendor rev");
                sff_show_ascii(id, SFF8636_VENDOR_SN_START_OFFSET,
                               SFF8636_VENDOR_SN_END_OFFSET, "Vendor SN");
                sff8636_show_revision_compliance(id);
                sff8636_show_dom(id, eeprom_len);
        }
}
*/

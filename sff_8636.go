package ethtool

import (
	"fmt"
)

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

func Decode(id []byte) (*SFF8636, error) {
	s := &SFF8636{
		Identifier: sff8636ShowIdentifier(id),
	}

	if id[SFF8636_ID_OFFSET] == SFF8024_ID_QSFP ||
		id[SFF8636_ID_OFFSET] == SFF8024_ID_QSFP_PLUS ||
		id[SFF8636_ID_OFFSET] == SFF8024_ID_QSFP28 {

		s.ExtIdentifier = sff8636ShowExtIdentifier(id)
		s.ExtIdentifierDescr = sff8636ShowExtIdentifierDescr(id)

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

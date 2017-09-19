package ethtool

func sff8024ShowIdentifier(id []byte, offset int) string {
	switch id[offset] {
	case SFF8024_ID_UNKNOWN:
		return "(no module present, unknown, or unspecified)"
	case SFF8024_ID_GBIC:
		return "(GBIC)"
	case SFF8024_ID_SOLDERED_MODULE:
		return "(module soldered to motherboard)"
	case SFF8024_ID_SFP:
		return "(SFP)"
	case SFF8024_ID_300_PIN_XBI:
		return "(300 pin XBI)"
	case SFF8024_ID_XENPAK:
		return "(XENPAK)"
	case SFF8024_ID_XFP:
		return "(XFP)"
	case SFF8024_ID_XFF:
		return "(XFF)"
	case SFF8024_ID_XFP_E:
		return "(XFP-E)"
	case SFF8024_ID_XPAK:
		return "(XPAK)"
	case SFF8024_ID_X2:
		return "(X2)"
	case SFF8024_ID_DWDM_SFP:
		return "(DWDM-SFP)"
	case SFF8024_ID_QSFP:
		return "(QSFP)"
	case SFF8024_ID_QSFP_PLUS:
		return "(QSFP+)"
	case SFF8024_ID_CXP:
		return "(CXP)"
	case SFF8024_ID_HD4X:
		return "(Shielded Mini Multilane HD 4X)"
	case SFF8024_ID_HD8X:
		return "(Shielded Mini Multilane HD 8X)"
	case SFF8024_ID_QSFP28:
		return "(QSFP28)"
	case SFF8024_ID_CXP2:
		return "(CXP2/CXP28)"
	case SFF8024_ID_CDFP:
		return "(CDFP Style 1/Style 2)"
	case SFF8024_ID_HD4X_FANOUT:
		return "(Shielded Mini Multilane HD 4X Fanout Cable)"
	case SFF8024_ID_HD8X_FANOUT:
		return "(Shielded Mini Multilane HD 8X Fanout Cable)"
	case SFF8024_ID_CDFP_S3:
		return "(CDFP Style 3)"
	case SFF8024_ID_MICRO_QSFP:
		return "(microQSFP)"
	}

	return "(reserved or unknown)"
}

func sff8024ShowConnector(id []byte, offset int) string {
	switch id[offset] {
	case SFF8024_CTOR_UNKNOWN:
		return "(unknown or unspecified)"
	case SFF8024_CTOR_SC:
		return "(SC)"
	case SFF8024_CTOR_FC_STYLE_1:
		return "(Fibre Channel Style 1 copper)"
	case SFF8024_CTOR_FC_STYLE_2:
		return "(Fibre Channel Style 2 copper)"
	case SFF8024_CTOR_BNC_TNC:
		return "(BNC/TNC)"
	case SFF8024_CTOR_FC_COAX:
		return "(Fibre Channel coaxial headers)"
	case SFF8024_CTOR_FIBER_JACK:
		return "(FibreJack)"
	case SFF8024_CTOR_LC:
		return "(LC)"
	case SFF8024_CTOR_MT_RJ:
		return "(MT-RJ)"
	case SFF8024_CTOR_MU:
		return "(MU)"
	case SFF8024_CTOR_SG:
		return "(SG)"
	case SFF8024_CTOR_OPT_PT:
		return "(Optical pigtail)"
	case SFF8024_CTOR_MPO:
		return "(MPO Parallel Optic)"
	case SFF8024_CTOR_MPO_2:
		return "(MPO Parallel Optic - 2x16)"
	case SFF8024_CTOR_HSDC_II:
		return "(HSSDC II)"
	case SFF8024_CTOR_COPPER_PT:
		return "(Copper pigtail)"
	case SFF8024_CTOR_RJ45:
		return "(RJ45)"
	case SFF8024_CTOR_NO_SEPARABLE:
		return "(No separable connector)"
	case SFF8024_CTOR_MXC_2x16:
		return "(MXC 2x16)"
	}
	return "(reserved or unknown)"
}

func sff8024ShowEncoding(id []byte, offset int, sffType int) string {
	switch id[offset] {
	case SFF8024_ENCODING_UNSPEC:
		return "(unspecified)"
	case SFF8024_ENCODING_8B10B:
		return "(8B/10B)"
	case SFF8024_ENCODING_4B5B:
		return "(4B/5B)"
	case SFF8024_ENCODING_NRZ:
		return "(NRZ)"
	case SFF8024_ENCODING_4h:
		if sffType == ETH_MODULE_SFF_8472 {
			return "(Manchester)"
		} else if sffType == ETH_MODULE_SFF_8636 {
			return "(SONET Scrambled)"
		}
	case SFF8024_ENCODING_5h:
		if sffType == ETH_MODULE_SFF_8472 {
			return "(SONET Scrambled)"
		} else if sffType == ETH_MODULE_SFF_8636 {
			return "(64B/66B)"
		}
	case SFF8024_ENCODING_6h:
		if sffType == ETH_MODULE_SFF_8472 {
			return "(64B/66B)"
		} else if sffType == ETH_MODULE_SFF_8636 {
			return "(Manchester)"
		}
	case SFF8024_ENCODING_256B:
		return "((256B/257B (transcoded FEC-enabled data))"
	case SFF8024_ENCODING_PAM4:
		return "(PAM4)"
	}
	return "(reserved or unknown)"
}

package ethtool

import ()

const (
	/* A0-based EEPROM offsets for DOM support checks */
	SFF_A0_DOM     = 92
	SFF_A0_OPTIONS = 93
	SFF_A0_COMP    = 94

	/* EEPROM bit values for various registers */
	SFF_A0_DOM_EXTCAL = (1 << 4)
	SFF_A0_DOM_INTCAL = (1 << 5)
	SFF_A0_DOM_IMPL   = (1 << 6)
	SFF_A0_DOM_PWRT   = (1 << 3)

	SFF_A0_OPTIONS_AW = (1 << 7)

	/* Offset for SFF-8472 */
	SFF_A2_BASE = 0x100

	/* A2-based offsets for DOM */
	SFF_A2_TEMP       = 96
	SFF_A2_TEMP_HALRM = 0
	SFF_A2_TEMP_LALRM = 2
	SFF_A2_TEMP_HWARN = 4
	SFF_A2_TEMP_LWARN = 6

	SFF_A2_VCC       = 98
	SFF_A2_VCC_HALRM = 8
	SFF_A2_VCC_LALRM = 10
	SFF_A2_VCC_HWARN = 12
	SFF_A2_VCC_LWARN = 14

	SFF_A2_BIAS       = 100
	SFF_A2_BIAS_HALRM = 16
	SFF_A2_BIAS_LALRM = 18
	SFF_A2_BIAS_HWARN = 20
	SFF_A2_BIAS_LWARN = 22

	SFF_A2_TX_PWR       = 102
	SFF_A2_TX_PWR_HALRM = 24
	SFF_A2_TX_PWR_LALRM = 26
	SFF_A2_TX_PWR_HWARN = 28
	SFF_A2_TX_PWR_LWARN = 30

	SFF_A2_RX_PWR       = 104
	SFF_A2_RX_PWR_HALRM = 32
	SFF_A2_RX_PWR_LALRM = 34
	SFF_A2_RX_PWR_HWARN = 36
	SFF_A2_RX_PWR_LWARN = 38

	SFF_A2_ALRM_FLG = 112
	SFF_A2_WARN_FLG = 116

	/* 32-bit little-endian calibration constants */
	SFF_A2_CAL_RXPWR4 = 56
	SFF_A2_CAL_RXPWR3 = 60
	SFF_A2_CAL_RXPWR2 = 64
	SFF_A2_CAL_RXPWR1 = 68
	SFF_A2_CAL_RXPWR0 = 72

	/* 16-bit little endian calibration constants */
	SFF_A2_CAL_TXI_SLP   = 76
	SFF_A2_CAL_TXI_OFF   = 78
	SFF_A2_CAL_TXPWR_SLP = 80
	SFF_A2_CAL_TXPWR_OFF = 82
	SFF_A2_CAL_T_SLP     = 84
	SFF_A2_CAL_T_OFF     = 86
	SFF_A2_CAL_V_SLP     = 88
	SFF_A2_CAL_V_OFF     = 90
)

type sff8472 struct {
}

func ParseSFF8472(id []byte) (sff8472, error) {
	return sff8472{}, nil
}

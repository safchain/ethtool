//go:build !linux

package ethtool

type LinkSettingSource string

// Constants defining the source of the LinkSettings data
const (
	SourceGLinkSettings LinkSettingSource = "GLINKSETTINGS"
	SourceGSet          LinkSettingSource = "GSET"
)

// EthtoolLinkSettingsFixed corresponds to struct ethtool_link_settings fixed part
type EthtoolLinkSettingsFixed struct {
	Cmd                 uint32
	Speed               uint32
	Duplex              uint8
	Port                uint8
	PhyAddress          uint8
	Autoneg             uint8
	MdixSupport         uint8 // Renamed from mdio_support
	EthTpMdix           uint8
	EthTpMdixCtrl       uint8
	LinkModeMasksNwords int8
	Transceiver         uint8
	MasterSlaveCfg      uint8
	MasterSlaveState    uint8
	Reserved1           [1]byte
	Reserved            [7]uint32
	// Flexible array member link_mode_masks[0] starts here implicitly
}

// LinkSettings is the user-friendly representation returned by GetLinkSettings
type LinkSettings struct {
	Speed                uint32
	Duplex               uint8
	Port                 uint8
	PhyAddress           uint8
	Autoneg              uint8
	MdixSupport          uint8
	EthTpMdix            uint8
	EthTpMdixCtrl        uint8
	Transceiver          uint8
	MasterSlaveCfg       uint8
	MasterSlaveState     uint8
	SupportedLinkModes   []string
	AdvertisingLinkModes []string
	LpAdvertisingModes   []string
	Source               LinkSettingSource // "GSET" or "GLINKSETTINGS"
}

// GetLinkSettings retrieves link settings, preferring ETHTOOL_GLINKSETTINGS and falling back to ETHTOOL_GSET.
// Uses a single ioctl call with the maximum expected buffer size.
func (e *Ethtool) GetLinkSettings(intf string) (*LinkSettings, error) {
	return nil, errOSUnsupported
}

// SetLinkSettings applies link settings, determining whether to use ETHTOOL_SLINKSETTINGS or ETHTOOL_SSET.
func (e *Ethtool) SetLinkSettings(intf string, settings *LinkSettings) error {
	return errOSUnsupported
}
